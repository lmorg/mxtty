package rendererwebkit

import (
	"testing"

	"github.com/lmorg/ttyphoon/types"
)

type paneLayoutTestTile struct {
	types.Tile
	id   string
	x, y int32
	term types.Term
}

func (tile *paneLayoutTestTile) Id() string          { return tile.id }
func (tile *paneLayoutTestTile) Left() int32         { return tile.x }
func (tile *paneLayoutTestTile) Top() int32          { return tile.y }
func (tile *paneLayoutTestTile) GetTerm() types.Term { return tile.term }

type paneLayoutTestTerm struct {
	types.Term
	size types.XY
}

func (term *paneLayoutTestTerm) GetSize() *types.XY { return &term.size }

func TestPaneLayoutTracksVisiblePaneGeometry(t *testing.T) {
	first := &paneLayoutTestTile{id: "%1", term: &paneLayoutTestTerm{size: types.XY{X: 80, Y: 24}}}
	secondTerm := &paneLayoutTestTerm{size: types.XY{X: 80, Y: 24}}
	second := &paneLayoutTestTile{id: "%2", x: 81, y: 3, term: secondTerm}
	renderer := &webkitRender{termWin: &types.AppWindowTerms{Tiles: []types.Tile{second, nil, first}}}
	panes := renderer.GetPaneLayout()
	if len(panes) != 2 || panes[0].ID != "%1" || panes[1].ID != "%2" {
		t.Fatalf("unexpected pane order: %+v", panes)
	}
	if panes[1].X != 81 || panes[1].Y != 3 || panes[1].Width != 80 || panes[1].Height != 24 || panes[1].CanClose {
		t.Fatalf("unexpected pane geometry: %+v", panes[1])
	}
	second.x = 0
	second.y = 0
	secondTerm.size = types.XY{X: 120, Y: 40}
	renderer.termWin.Tiles = []types.Tile{second}
	panes = renderer.GetPaneLayout()
	if len(panes) != 1 || panes[0].ID != "%2" || panes[0].X != 0 || panes[0].Y != 0 || panes[0].Width != 120 || panes[0].Height != 40 {
		t.Fatalf("stale pane layout: %+v", panes)
	}
	renderer.termWin.Tiles = nil
	if len(renderer.GetPaneLayout()) != 0 {
		t.Fatal("removed panes remain in layout")
	}
}
