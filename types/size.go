package types

type XY struct {
	X int32
	Y int32
}

type Rect struct {
	Start *XY
	End   *XY
}
