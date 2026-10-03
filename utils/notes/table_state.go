package notes

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/lmorg/ttyphoon/utils/cache"
	"github.com/lmorg/ttyphoon/utils/graphs"
)

type TableState struct {
	SortColumn int          `json:"sortColumn"`
	SortDesc   bool         `json:"sortDesc"`
	Filter     string       `json:"filter"`
	Graph      *graphs.Spec `json:"graph,omitempty"`
}

func tableStateCacheKey(workspace, filename string, index int, headings []string) string {
	normalized := make([]string, len(headings))
	for i := range headings {
		normalized[i] = strings.TrimSpace(headings[i])
	}
	hash := sha1.Sum([]byte(strings.Join(normalized, "\x1f")))
	return fmt.Sprintf("%s:%s:%d:%s", workspace, filename, index, hex.EncodeToString(hash[:]))
}

func GetTableState(workspace, filename string, index int, headings []string) TableState {
	var state TableState
	cache.Read(cache.NS_NOTESW_TABLE_STATE, tableStateCacheKey(workspace, filename, index, headings), &state)
	return state
}

func GetTableStateWithLegacyIndex(workspace, filename string, index, legacyIndex int, headings []string) TableState {
	var state TableState
	if cache.Read(cache.NS_NOTESW_TABLE_STATE, tableStateCacheKey(workspace, filename, index, headings), &state) {
		return state
	}
	if legacyIndex != index {
		if cache.Read(cache.NS_NOTESW_TABLE_STATE, tableStateCacheKey(workspace, filename, legacyIndex, headings), &state) {
			return state
		}
	}
	for distance := 1; distance <= 8; distance++ {
		for _, candidateIndex := range []int{legacyIndex + distance, legacyIndex - distance} {
			if candidateIndex < 0 || candidateIndex == index || candidateIndex == legacyIndex {
				continue
			}
			var candidate TableState
			if !cache.Read(cache.NS_NOTESW_TABLE_STATE, tableStateCacheKey(workspace, filename, candidateIndex, headings), &candidate) || candidate.Graph == nil {
				continue
			}
			state.Graph = candidate.Graph
			SetTableState(workspace, filename, index, headings, state)
			return state
		}
	}
	return state
}

func SetTableState(workspace, filename string, index int, headings []string, state TableState) {
	cache.Write(cache.NS_NOTESW_TABLE_STATE, tableStateCacheKey(workspace, filename, index, headings), state, cache.Days(30))
}
