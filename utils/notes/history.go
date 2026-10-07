package notes

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/lmorg/murex/utils/lists"
	"github.com/lmorg/ttyphoon/utils/cache"
)

type HistoryListT struct {
	History []HistoryEntryT
	Index   int
}

type HistoryEntryT struct {
	Filename string `json:"filename"`
	Line     int    `json:"line"`
}

func (entry *HistoryEntryT) UnmarshalJSON(data []byte) error {
	var filename string
	if err := json.Unmarshal(data, &filename); err == nil {
		entry.Filename = filename
		entry.Line = 1
		return nil
	}

	type historyEntryAlias HistoryEntryT
	var decoded historyEntryAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*entry = HistoryEntryT(decoded)
	if entry.Line < 1 {
		entry.Line = 1
	}
	return nil
}

func getHistoryList(projectRoot string) *HistoryListT {
	h := &HistoryListT{}

	cache.Read(cache.NS_NOTESW_HISTORY, projectRoot, h)

	return h
}

func setHistoryList(projectRoot string, h *HistoryListT) {
	cache.Write(cache.NS_NOTESW_HISTORY, projectRoot, h, cache.Days(30))
}

func HistoryListCurrent(projectRoot string) string {
	return HistoryListCurrentEntry(projectRoot).Filename
}

func HistoryListCurrentEntry(projectRoot string) HistoryEntryT {
	h := getHistoryList(projectRoot)
	if len(h.History) == 0 {
		return HistoryEntryT{}
	}
	if h.Index < 0 || h.Index >= len(h.History) {
		return HistoryEntryT{}
	}

	return h.History[h.Index]
}

func HistoryListAdd(projectRoot, filename string) error {
	return HistoryListAddAt(projectRoot, filename, 1)
}

func HistoryListAddAt(projectRoot, filename string, line int) error {
	if filename == "" {
		return errors.New("no filename to add to recent list")
	}

	filename = filepath.Clean(filename)
	if line < 1 {
		line = 1
	}

	h := getHistoryList(projectRoot)
	if h.Index >= 0 && h.Index < len(h.History) {
		current := h.History[h.Index]
		if filepath.Clean(current.Filename) == filename && abs(current.Line-line) < 20 {
			return nil
		}
	}

	truncateTo := h.Index + 1
	if truncateTo < 0 {
		truncateTo = 0
	}
	if truncateTo > len(h.History) {
		truncateTo = len(h.History)
	}

	h.History = h.History[:truncateTo]

	h.History = append(h.History, HistoryEntryT{Filename: filename, Line: line})
	if len(h.History) > 50 {
		h.History = h.History[len(h.History)-50:]
	}

	h.Index = len(h.History) - 1
	setHistoryList(projectRoot, h)
	return nil
}

func HistoryListRename(projectRoot, oldFile, newFile string) error {
	if oldFile == "" || newFile == "" || oldFile == newFile {
		return fmt.Errorf("cannot update recent list: from '%s' to '%s'", oldFile, newFile)
	}

	oldFile = filepath.Clean(oldFile)
	newFile = filepath.Clean(newFile)

	h := getHistoryList(projectRoot)

	for i := range h.History {
		if filepath.Clean(h.History[i].Filename) == oldFile {
			h.History[i].Filename = newFile
		}
	}

	setHistoryList(projectRoot, h)
	return nil
}

func HistoryListDelete(projectRoot, filename string) error {
	if filename == "" {
		return errors.New("no filename to delete from recent list")
	}

	filename = filepath.Clean(filename)

	h := getHistoryList(projectRoot)

	var err error
	for {
		i := -1
		for idx := range h.History {
			if filepath.Clean(h.History[idx].Filename) == filename {
				i = idx
				break
			}
		}
		if i == -1 {
			if len(h.History) == 0 {
				h.Index = 0
			} else if h.Index >= len(h.History) {
				h.Index = len(h.History) - 1
			} else if h.Index < 0 {
				h.Index = 0
			}
			setHistoryList(projectRoot, h)
			return nil
		}

		if i < h.Index {
			h.Index--
		}

		h.History, err = lists.RemoveOrdered(h.History, i)
		if err != nil {
			return err
		}
	}
}

// Previous returns a filename or error
func HistoryListPrevious(projectRoot string) (string, error) {
	entry, err := HistoryListPreviousEntry(projectRoot)
	return entry.Filename, err
}

func HistoryListPreviousEntry(projectRoot string) (HistoryEntryT, error) {
	h := getHistoryList(projectRoot)

	if h.Index == 0 {
		return HistoryEntryT{}, errors.New("You already have the first history item open")
	}

	h.Index--
	setHistoryList(projectRoot, h)
	return h.History[h.Index], nil
}

// HistoryListNext returns a filename or error
func HistoryListNext(projectRoot string) (string, error) {
	entry, err := HistoryListNextEntry(projectRoot)
	return entry.Filename, err
}

func HistoryListNextEntry(projectRoot string) (HistoryEntryT, error) {
	h := getHistoryList(projectRoot)

	if h.Index+1 == len(h.History) {
		return HistoryEntryT{}, errors.New("You already have the last history item open")
	}

	h.Index++
	setHistoryList(projectRoot, h)
	return h.History[h.Index], nil
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
