package codes

var (
	AnsiUp        = []byte{27, 91, 65}
	AnsiDown      = []byte{27, 91, 66}
	AnsiForwards  = []byte{27, 91, 67}
	AnsiBackwards = []byte{27, 91, 68}
	AnsiHome      = []byte{27, 91, 72}
	AnsiHomeSc    = []byte{27, 91, 49, 126}
	AnsiEnd       = []byte{27, 91, 70}
	AnsiEndSc     = []byte{27, 91, 52, 126}
	AnsiDelete    = []byte{27, 91, 51, 126}
	AnsiShiftTab  = []byte{27, 91, 90}
	AnsiPageUp    = []byte{27, 91, 53, 126}
	AnsiPageDown  = []byte{27, 91, 54, 126}
	AnsiOptUp     = []byte{27, 27, 91, 65}
	AnsiOptDown   = []byte{27, 27, 91, 66}
	AnsiOptLeft   = []byte{27, 27, 91, 68}
	AnsiOptRight  = []byte{27, 27, 91, 67}
	AnsiCtrlUp    = []byte{27, 91, 49, 59, 53, 65}
	AnsiCtrlDown  = []byte{27, 91, 49, 59, 53, 66}
	AnsiCtrlLeft  = []byte{27, 91, 49, 59, 53, 68}
	AnsiCtrlRight = []byte{27, 91, 49, 59, 53, 67}

	AnsiF1VT100 = []byte{27, 79, 80}
	AnsiF2VT100 = []byte{27, 79, 81}
	AnsiF3VT100 = []byte{27, 79, 82}
	AnsiF4VT100 = []byte{27, 79, 83}
	AnsiF1      = []byte{27, 91, 49, 49, 126}
	AnsiF2      = []byte{27, 91, 49, 50, 126}
	AnsiF3      = []byte{27, 91, 49, 51, 126}
	AnsiF4      = []byte{27, 91, 49, 52, 126}
	AnsiF5      = []byte{27, 91, 49, 53, 126}
	AnsiF6      = []byte{27, 91, 49, 55, 126}
	AnsiF7      = []byte{27, 91, 49, 56, 126}
	AnsiF8      = []byte{27, 91, 49, 57, 126}
	AnsiF9      = []byte{27, 91, 50, 48, 126}
	AnsiF10     = []byte{27, 91, 50, 49, 126}
	AnsiF11     = []byte{27, 91, 50, 51, 126}
	AnsiF12     = []byte{27, 91, 50, 52, 126}

	AnsiShiftF1  = []byte{27, 91, 49, 59, 50, 80}
	AnsiShiftF2  = []byte{27, 91, 49, 59, 50, 81}
	AnsiShiftF3  = []byte{27, 91, 49, 59, 50, 82}
	AnsiShiftF4  = []byte{27, 91, 49, 59, 50, 83}
	AnsiShiftF5  = []byte{27, 91, 49, 53, 59, 50, 126}
	AnsiShiftF6  = []byte{27, 91, 49, 55, 59, 50, 126}
	AnsiShiftF7  = []byte{27, 91, 49, 56, 59, 50, 126}
	AnsiShiftF8  = []byte{27, 91, 49, 57, 59, 50, 126}
	AnsiShiftF9  = []byte{27, 91, 50, 48, 59, 50, 126}
	AnsiShiftF10 = []byte{27, 91, 50, 49, 59, 50, 126}
	AnsiShiftF11 = []byte{27, 91, 50, 51, 59, 50, 126}
	AnsiShiftF12 = []byte{27, 91, 50, 52, 59, 50, 126}
)
