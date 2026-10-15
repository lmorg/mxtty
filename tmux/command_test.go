package tmux

import (
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/lmorg/ttyphoon/types"
)

func TestCommandLine(t *testing.T) {
	expected := `list-windows  -F 'Name:#{window_name}|||Id:#{window_id}|||Width:#{window_width}|||Height:#{window_height}|||Active:#{?window_active,true,false}'`

	var _test_CMD_LIST_WINDOWS = "list-windows"

	type _test_CMD_LIST_WINDOWS_T struct {
		Name       string `tmux:"window_name"`
		Id         string `tmux:"window_id"`
		Width      int    `tmux:"window_width"`
		Height     int    `tmux:"window_height"`
		Active     bool   `tmux:"?window_active,true,false"`
		panes      map[string]*PaneT
		activePane *PaneT
	}

	actual := string(mkCmdLine(
		_test_CMD_LIST_WINDOWS,
		reflect.TypeOf(_test_CMD_LIST_WINDOWS_T{}),
	))

	if expected != actual {
		t.Errorf("Incorrect string returned:")
		t.Logf("  Expected: %s", expected)
		t.Logf("  Actual:   %s", actual)
		t.Logf("  exp byte: %v", []byte(expected))
		t.Logf("  act byte: %v", []byte(actual))
	}
}

func TestQuoteStringEscapesTmuxArgument(t *testing.T) {
	if got, want := quoteString(`SSH_AUTH_SOCK=/tmp/agent\"value`), `"SSH_AUTH_SOCK=/tmp/agent\\\"value"`; got != want {
		t.Fatalf("quoteString() = %q, want %q", got, want)
	}
}

type layoutChangeTestRenderer struct {
	types.Renderer
	refreshed chan struct{}
	release   <-chan struct{}
}

func (renderer *layoutChangeTestRenderer) RefreshWindowList() {
	renderer.refreshed <- struct{}{}
	if renderer.release != nil {
		<-renderer.release
	}
}

func TestPaneResizeDoesNotEchoKnownTmuxGeometry(t *testing.T) {
	tmux := &Tmux{}
	tmux.readerDead.Store(true)
	pane := &PaneT{id: "%7", width: 80, height: 24, tmux: tmux}
	if err := pane.Resize(&types.XY{X: 80, Y: 24}); err != nil {
		t.Fatalf("unchanged geometry attempted a tmux command: %v", err)
	}
	if err := pane.Resize(&types.XY{X: 100, Y: 24}); err == nil {
		t.Fatal("a requested size change did not attempt a tmux command")
	}
}

type paneGeometryTestTerm struct {
	types.Term
	pane      *PaneT
	resizeErr error
	resizes   int
}

func (term *paneGeometryTestTerm) MakeVisible(bool) {}
func (term *paneGeometryTestTerm) SetFocus(bool)    {}
func (term *paneGeometryTestTerm) Resize(size *types.XY) {
	term.resizes++
	term.resizeErr = term.pane.Resize(size)
}

func TestPaneInfoResizeUpdatesLocalTerminalWithoutEchoingTmux(t *testing.T) {
	tmux := &Tmux{panes: newPaneMap(), wins: newWindowMap()}
	tmux.readerDead.Store(true)
	win := &WindowT{id: "@1", panes: newPaneMap()}
	tmux.wins.Set(win.id, win)
	pane := &PaneT{id: "%7", width: 40, height: 12, tmux: tmux}
	term := &paneGeometryTestTerm{pane: pane}
	pane.term = term
	tmux.panes.Set(pane.id, pane)
	info := &paneInfo{Id: pane.id, WindowId: win.id, Width: 80, Height: 24}
	info.updatePane(tmux)
	info.updatePane(tmux)
	if term.resizes != 2 || term.resizeErr != nil || pane.width != 80 || pane.height != 24 {
		t.Fatalf("geometry refresh echoed a resize: calls=%d error=%v size=%dx%d", term.resizes, term.resizeErr, pane.width, pane.height)
	}
}

func TestLayoutChangeRefreshesPaneGeometry(t *testing.T) {
	renderer := &layoutChangeTestRenderer{refreshed: make(chan struct{}, 1)}
	tmuxCommandMap[_RESP_LAYOUT_CHANGE](&Tmux{renderer: renderer}, []byte("%layout-change @1 layout visible-layout *"))
	select {
	case <-renderer.refreshed:
	case <-time.After(time.Second):
		t.Fatal("tmux layout change did not refresh pane geometry")
	}
}

func TestLayoutChangeBurstCoalescesIntoOneTrailingRefresh(t *testing.T) {
	release := make(chan struct{})
	defer close(release)
	renderer := &layoutChangeTestRenderer{refreshed: make(chan struct{}, 128), release: release}
	tmux := &Tmux{renderer: renderer}
	_respLayoutChange(tmux, nil)
	select {
	case <-renderer.refreshed:
	case <-time.After(time.Second):
		t.Fatal("initial layout refresh did not start")
	}
	for range 100 {
		_respLayoutChange(tmux, nil)
	}
	tmux.layoutRefreshMu.Lock()
	running, pending := tmux.layoutRefreshRunning, tmux.layoutRefreshPending
	tmux.layoutRefreshMu.Unlock()
	if !running || !pending || len(renderer.refreshed) != 0 {
		t.Fatal("layout events were not coalesced behind the active refresh")
	}
	release <- struct{}{}
	select {
	case <-renderer.refreshed:
	case <-time.After(time.Second):
		t.Fatal("final layout refresh was lost")
	}
	tmux.layoutRefreshMu.Lock()
	pending = tmux.layoutRefreshPending
	tmux.layoutRefreshMu.Unlock()
	if pending || len(renderer.refreshed) != 0 {
		t.Fatal("layout burst scheduled more than one trailing refresh")
	}
}

func TestCompletedPaneStateAndCloseGuard(t *testing.T) {
	tmux := &Tmux{panes: newPaneMap()}
	pane := &PaneT{id: "%7", tmux: tmux}
	tmux.panes.Set(pane.id, pane)
	if pane.CommandFinished() {
		t.Fatal("running pane is marked completed")
	}
	if err := tmux.CloseCompletedPane(pane.id); err == nil || err.Error() != "pane command is still running" {
		t.Fatalf("running pane close error = %v", err)
	}
	if err := tmux.CloseCompletedPane("%missing"); err == nil {
		t.Fatal("missing pane was accepted")
	}
	pane.dead.Store(true)
	if !pane.CommandFinished() {
		t.Fatal("retained dead pane is not marked completed")
	}
	pane.dead.Store(false)
	for _, result := range []int32{1, 2} {
		pane.actionResult.Store(result)
		if !pane.CommandFinished() {
			t.Fatalf("interactive result %d is not marked completed", result)
		}
	}
}

func TestRespawnedPaneClearsCompletionState(t *testing.T) {
	tmux := &Tmux{panes: newPaneMap(), wins: newWindowMap()}
	win := &WindowT{id: "@1", panes: newPaneMap()}
	tmux.wins.Set(win.id, win)
	pane := &PaneT{id: "%7", tmux: tmux}
	pane.actionResult.Store(2)
	pane.autoCloseAt.Store(time.Now().Add(30 * time.Second).UnixMilli())
	tmux.panes.Set(pane.id, pane)
	info := &paneInfo{Id: pane.id, WindowId: win.id, Dead: true}
	info.updatePane(tmux)
	info.Dead = false
	info.updatePane(tmux)
	if pane.CommandFinished() || pane.AutoCloseAt() != 0 {
		t.Fatal("respawned pane still has completed-command state")
	}
}

func TestClonePaneCommandKeepsExitedPaneAndStartsWrapper(t *testing.T) {
	command := clonePaneCommand("%7", `-e "TOKEN=a b" `, "/tmp/ttyphoon action.sh")
	want := `split-window -v -k -P -F '#{pane_id}' -t %7 -c "#{pane_current_path}" -e "TOKEN=a b" bash "/tmp/ttyphoon action.sh"`
	if command != want {
		t.Fatalf("clonePaneCommand() = %q, want %q", command, want)
	}
}

func TestInteractiveActionCallbackRequiresMatchingUUIDAndRunsOnce(t *testing.T) {
	tmux := &Tmux{panes: newPaneMap()}
	pane := &PaneT{id: "%7", tmux: tmux}
	tmux.panes.Set("%7", pane)
	var calls int
	var gotOutput []byte
	var gotExitCode int
	var gotErr error
	id, err := newInteractiveActionID()
	if err != nil {
		t.Fatal(err)
	}
	tmux.registerInteractiveAction(id, "", func(output []byte, exitCode int, callbackErr error) {
		calls++
		gotOutput = output
		gotExitCode = exitCode
		gotErr = callbackErr
	})
	if tmux.completeInteractiveAction("different-id", "%7", []byte("wrong"), 0) {
		t.Fatal("completion with a different UUID was accepted")
	}
	if !tmux.completeInteractiveAction(id, "%7", []byte("git output"), 1) {
		t.Fatal("early completion was not queued")
	}
	if calls != 0 {
		t.Fatal("callback ran before pane initialization completed")
	}
	if !tmux.bindInteractiveActionPane(id, "%7") {
		t.Fatal("failed to bind callback to initialized pane")
	}
	if tmux.completeInteractiveAction(id, "%8", []byte("wrong pane"), 0) {
		t.Fatal("completion from a different pane was accepted")
	}
	if tmux.completeInteractiveAction(id, "%7", []byte("duplicate"), 0) {
		t.Fatal("duplicate completion was accepted")
	}
	if calls != 1 || string(gotOutput) != "git output" || gotExitCode != 1 || gotErr != nil {
		t.Fatalf("callback result = calls:%d output:%q exit:%d error:%v", calls, gotOutput, gotExitCode, gotErr)
	}
	if pane.InteractiveActionChrome() != types.COLOR_ERROR {
		t.Fatal("failed interactive action did not retain error chrome")
	}
	if pane.AutoCloseAt() != 0 {
		t.Fatal("failed interactive action has an auto-close deadline")
	}
}

func TestInteractiveActionRegistryDoesNotReplaceExistingUUID(t *testing.T) {
	tmux := &Tmux{}
	firstCalls, secondCalls := 0, 0
	if !tmux.registerInteractiveAction("same-id", "", func([]byte, int, error) { firstCalls++ }) {
		t.Fatal("first UUID registration failed")
	}
	if tmux.registerInteractiveAction("same-id", "", func([]byte, int, error) { secondCalls++ }) {
		t.Fatal("duplicate UUID registration unexpectedly succeeded")
	}
	if !tmux.bindInteractiveActionPane("same-id", "%11") {
		t.Fatal("failed to bind original callback to pane")
	}
	if !tmux.completeInteractiveAction("same-id", "%11", nil, 1) {
		t.Fatal("original UUID callback did not complete")
	}
	if firstCalls != 1 || secondCalls != 0 {
		t.Fatalf("callback calls = first:%d second:%d, want first only", firstCalls, secondCalls)
	}
}

func TestSuccessfulInteractiveActionUsesSuccessChrome(t *testing.T) {
	tmux := &Tmux{panes: newPaneMap()}
	pane := &PaneT{id: "%10", tmux: tmux}
	tmux.panes.Set("%10", pane)
	t.Cleanup(func() { tmux.panes.Delete(pane.id) })
	id, err := newInteractiveActionID()
	if err != nil {
		t.Fatal(err)
	}
	tmux.registerInteractiveAction(id, "", func([]byte, int, error) {})
	if !tmux.bindInteractiveActionPane(id, pane.id) {
		t.Fatal("failed to bind callback to pane")
	}
	if !tmux.completeInteractiveAction(id, pane.id, []byte("success"), 0) {
		t.Fatal("matching successful completion was not accepted")
	}
	if pane.InteractiveActionChrome() != types.COLOR_OK {
		t.Fatal("successful interactive action did not retain success chrome")
	}
	deadline := pane.AutoCloseAt()
	remaining := time.Until(time.UnixMilli(deadline))
	if remaining < 29*time.Second || remaining > 30*time.Second {
		t.Fatalf("auto-close deadline is %s away, want 30 seconds", remaining)
	}
	if err := tmux.CancelPaneAutoClose(pane.id); err != nil {
		t.Fatal(err)
	}
	if pane.AutoCloseAt() != 0 || pane.autoCloseAt.CompareAndSwap(deadline, -1) {
		t.Fatal("cancelled auto-close can still fire")
	}
	if !pane.CommandFinished() {
		t.Fatal("cancelling auto-close disabled manual close")
	}
}

func TestCancelPaneAutoCloseRejectsAlreadyClosingPane(t *testing.T) {
	tmux := &Tmux{panes: newPaneMap()}
	pane := &PaneT{id: "%10", tmux: tmux}
	tmux.panes.Set(pane.id, pane)
	pane.autoCloseAt.Store(-1)
	if err := tmux.CancelPaneAutoClose(pane.id); err == nil {
		t.Fatal("cancellation was accepted after auto-close started")
	}
	if err := tmux.CancelPaneAutoClose("%missing"); err == nil {
		t.Fatal("cancellation was accepted for a missing pane")
	}
}

func TestClosingInteractiveActionPaneFailsPendingCallback(t *testing.T) {
	tmux := &Tmux{}
	var callbackErr error
	id, err := newInteractiveActionID()
	if err != nil {
		t.Fatal(err)
	}
	tmux.registerInteractiveAction(id, "", func(_ []byte, _ int, err error) {
		callbackErr = err
	})
	if !tmux.bindInteractiveActionPane(id, "%9") {
		t.Fatal("failed to bind callback to pane")
	}

	tmux.cancelInteractiveActionsForPane("%9")
	if callbackErr == nil || callbackErr.Error() != "interactive action pane closed before completion" {
		t.Fatalf("callback error = %v, want a pane-closed error", callbackErr)
	}
	if tmux.completeInteractiveAction(id, "%9", nil, 0) {
		t.Fatal("completion was accepted after pane-close cancellation")
	}
}

func TestWriteBashScriptEmitsCallbackUUIDAndExitCode(t *testing.T) {
	path, err := writeBashScript("git commit -m 'message'", "550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })

	script, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"set +e",
		"_ttyphoon_exit_code=$?",
		`ttyphoon;interactive-action-complete;{"id":"%s","exitCode":%s}`,
		"550e8400-e29b-41d4-a716-446655440000",
	} {
		if !strings.Contains(string(script), want) {
			t.Errorf("Bash script does not contain %q", want)
		}
	}
}

func TestBashScriptEmitsCompletionAfterCommandFailure(t *testing.T) {
	const id = "550e8400-e29b-41d4-a716-446655440000"
	path, err := writeBashScript("printf 'hook failed\\n'; false", id)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })

	output, err := exec.Command("bash", path).CombinedOutput()
	if err != nil {
		t.Fatalf("wrapper script exited unexpectedly: %v, output %q", err, output)
	}
	if !strings.Contains(string(output), "hook failed") || !strings.Contains(string(output), `ttyphoon;interactive-action-complete;{"id":"`+id+`","exitCode":1}`) {
		t.Fatalf("wrapper output %q does not include command output and failed completion APC", output)
	}
}
