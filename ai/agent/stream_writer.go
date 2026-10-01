package agent

import (
	"sync"
	"time"

	"github.com/lmorg/ttyphoon/ai/agent/sessiondb"
	"github.com/lmorg/ttyphoon/types"
)

type aiStreamBlockWriter struct {
	workspace string
	sessionID int64
	runID     uint64
	emit      func(sessiondb.AIStreamBlock)

	mu          sync.Mutex
	nextOrdinal int64
	// Rolling blocks are keyed by the owning scope (a sub-agent's block ID, or ""
	// for the root agent) so concurrent sub-agents never share a block.
	rolling map[rollingBlockKey]*aiStreamBlockHandle
}

type rollingBlockKey struct {
	kind  sessiondb.StreamBlockKind
	scope string
}

type aiStreamBlockHandle struct {
	writer  *aiStreamBlockWriter
	block   sessiondb.AIStreamBlock
	mu      sync.Mutex
	closed  bool
	pending string
	timer   *time.Timer
}

const streamBlockFlushInterval = types.AI_EMIT_INTERVAL * time.Millisecond

func newAIStreamBlockWriter(workspace string, emit func(sessiondb.AIStreamBlock)) *aiStreamBlockWriter {
	identity := sessiondb.GetActiveStreamIdentity(workspace)
	if identity.SessionID <= 0 || identity.RunID == 0 {
		return nil
	}
	return &aiStreamBlockWriter{
		workspace: workspace,
		sessionID: identity.SessionID,
		runID:     identity.RunID,
		emit:      emit,
	}
}

func (w *aiStreamBlockWriter) Open(kind sessiondb.StreamBlockKind, parentID string) *aiStreamBlockHandle {
	return w.OpenWithLabel(kind, parentID, "")
}

func (w *aiStreamBlockWriter) OpenWithLabel(kind sessiondb.StreamBlockKind, parentID, label string) *aiStreamBlockHandle {
	if w == nil {
		return nil
	}
	// Only the scope opening the block stops thinking; a sibling sub-agent's
	// tool call must not truncate another agent's reasoning mid-sentence.
	if kind != sessiondb.StreamBlockThinking {
		w.CloseCurrentThinking(parentID)
	}
	w.mu.Lock()
	ordinal := w.nextOrdinal
	w.nextOrdinal++
	w.mu.Unlock()

	block := sessiondb.AIStreamBlock{
		SessionID: w.sessionID,
		RunID:     w.runID,
		Workspace: w.workspace,
		ParentID:  parentID,
		Kind:      kind,
		Label:     label,
		Ordinal:   ordinal,
		Status:    "open",
	}
	blockID, err := sessiondb.CreateStreamBlockWithDatabaseID(w.workspace, block, streamBlockNow())
	if err != nil {
		return nil
	}
	block.BlockID = blockID
	w.emitBlock(block)
	return &aiStreamBlockHandle{writer: w, block: block}
}

func (w *aiStreamBlockWriter) AppendText(text, scope string) {
	w.appendRolling(sessiondb.StreamBlockText, text, scope)
}

func (w *aiStreamBlockWriter) AppendThinking(text, scope string) {
	w.appendRolling(sessiondb.StreamBlockThinking, text, scope)
}

func (w *aiStreamBlockWriter) AppendNotice(text, scope string) {
	w.appendRolling(sessiondb.StreamBlockNotice, text, scope)
}

// appendRolling appends to the scope's current block of this kind, reopening it
// if a concurrent writer closed it between lookup and append.
func (w *aiStreamBlockWriter) appendRolling(kind sessiondb.StreamBlockKind, text, scope string) {
	if w == nil || text == "" {
		return
	}
	for attempt := 0; attempt < 3; attempt++ {
		handle := w.rollingHandle(kind, scope)
		if handle == nil {
			return
		}
		if handle.appendIfOpen(text) {
			return
		}
	}
}

func (w *aiStreamBlockWriter) rollingHandle(kind sessiondb.StreamBlockKind, scope string) *aiStreamBlockHandle {
	key := rollingBlockKey{kind: kind, scope: scope}

	w.mu.Lock()
	existing := w.rolling[key]
	w.mu.Unlock()
	if existing != nil && !existing.isClosed() {
		return existing
	}

	opened := w.Open(kind, scope)
	if opened == nil {
		return nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if current := w.rolling[key]; current != nil && !current.isClosed() {
		// Another goroutine won the race; discard the block we just opened.
		opened.Close()
		return current
	}
	if w.rolling == nil {
		w.rolling = make(map[rollingBlockKey]*aiStreamBlockHandle)
	}
	w.rolling[key] = opened
	return opened
}

func (w *aiStreamBlockWriter) AppendNoticeForParent(text, parentID string) {
	if parentID == "" {
		w.AppendNotice(text, "")
		return
	}
	w.AppendStandalone(sessiondb.StreamBlockNotice, text, parentID)
}

func (w *aiStreamBlockWriter) AppendStandalone(kind sessiondb.StreamBlockKind, text, parentID string) string {
	return w.AppendStandaloneWithLabel(kind, text, parentID, "")
}

func (w *aiStreamBlockWriter) AppendStandaloneWithLabel(kind sessiondb.StreamBlockKind, text, parentID, label string) string {
	if w == nil || text == "" {
		return ""
	}
	block := w.OpenWithLabel(kind, parentID, label)
	if block == nil {
		return ""
	}
	block.Append(text)
	block.Close()
	return block.block.BlockID
}

func (w *aiStreamBlockWriter) CloseCurrentThinking(scope string) {
	if w == nil {
		return
	}
	w.mu.Lock()
	handle := w.rolling[rollingBlockKey{kind: sessiondb.StreamBlockThinking, scope: scope}]
	w.mu.Unlock()
	handle.Close()
}

func (w *aiStreamBlockWriter) CloseAll() {
	if w == nil {
		return
	}
	w.mu.Lock()
	blocks := make([]*aiStreamBlockHandle, 0, len(w.rolling))
	for _, block := range w.rolling {
		blocks = append(blocks, block)
	}
	w.mu.Unlock()
	for _, block := range blocks {
		block.Close()
	}
}

func (h *aiStreamBlockHandle) Append(delta string) {
	h.appendIfOpen(delta)
}

// appendIfOpen reports false when the block was closed concurrently, so the
// caller can reopen rather than silently dropping the delta.
func (h *aiStreamBlockHandle) appendIfOpen(delta string) bool {
	if h == nil || h.writer == nil || delta == "" {
		return true
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return false
	}
	h.pending += delta
	if h.timer == nil {
		h.timer = time.AfterFunc(streamBlockFlushInterval, h.flush)
	}
	return true
}

func (h *aiStreamBlockHandle) Close() {
	if h == nil || h.writer == nil {
		return
	}
	h.mu.Lock()
	if h.closed {
		h.mu.Unlock()
		return
	}
	h.closed = true
	if h.timer != nil {
		h.timer.Stop()
		h.timer = nil
	}
	pending := h.pending
	h.pending = ""
	h.mu.Unlock()
	if err := sessiondb.AppendStreamBlock(h.writer.workspace, h.block.SessionID, h.block.RunID, h.block.BlockID, pending, "closed", streamBlockNow()); err == nil {
		event := h.block
		event.Content = ""
		event.Delta = pending
		event.Status = "closed"
		h.writer.emitBlock(event)
	}
}

func (h *aiStreamBlockHandle) flush() {
	if h == nil || h.writer == nil {
		return
	}
	h.mu.Lock()
	if h.closed || h.pending == "" {
		h.timer = nil
		h.mu.Unlock()
		return
	}
	pending := h.pending
	h.pending = ""
	h.timer = nil
	err := sessiondb.AppendStreamBlock(h.writer.workspace, h.block.SessionID, h.block.RunID, h.block.BlockID, pending, "", streamBlockNow())
	if err != nil {
		h.pending = pending + h.pending
	}
	h.mu.Unlock()
	if err == nil {
		event := h.block
		event.Content = ""
		event.Delta = pending
		h.writer.emitBlock(event)
	}
	if err != nil {
		h.mu.Lock()
		if !h.closed && h.timer == nil {
			h.timer = time.AfterFunc(streamBlockFlushInterval, h.flush)
		}
		h.mu.Unlock()
	}
}

func (h *aiStreamBlockHandle) isClosed() bool {
	if h == nil {
		return true
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.closed
}

func streamBlockNow() string {
	return time.Now().Format("2006-01-02 15:04:05.000000000Z07:00")
}

func (w *aiStreamBlockWriter) emitBlock(block sessiondb.AIStreamBlock) {
	if w != nil && w.emit != nil {
		w.emit(block)
	}
}
