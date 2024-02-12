package virtualterm

import (
	"testing"
)

func TestCsiScrollUp(t *testing.T) {
	test := testTerm{
		Tests: []testCondition{
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "2222222222\n3333333333\n4444444444",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollUp(1)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "3333333333\n4444444444",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollUp(2)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "4444444444",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollUp(3)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollUp(4)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollUp(5)
				},
			},
		},
	}

	test.RunTests(t)
}

func TestCsiScrollDown(t *testing.T) {
	test := testTerm{
		Tests: []testCondition{
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n1111111111\n2222222222\n3333333333",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollDown(1)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n1111111111\n2222222222",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollDown(2)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n..........\n1111111111",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollDown(3)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n..........\n..........",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollDown(4)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n..........\n..........",
				Operation: func(t *testing.T, term *Term) {
					term.csiScrollDown(5)
				},
			},
		},
	}

	test.RunTests(t)
}

func Test_scrollDown(t *testing.T) {
	test := testTerm{
		Tests: []testCondition{
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n1111111111\n2222222222\n3333333333",
				Operation: func(t *testing.T, term *Term) {
					top, bottom := term.getScrollingRegion()
					term._scrollDown(top, bottom, 1)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n1111111111\n2222222222",
				Operation: func(t *testing.T, term *Term) {
					top, bottom := term.getScrollingRegion()
					term._scrollDown(top, bottom, 2)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n..........\n1111111111",
				Operation: func(t *testing.T, term *Term) {
					top, bottom := term.getScrollingRegion()
					term._scrollDown(top, bottom, 3)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n..........\n..........",
				Operation: func(t *testing.T, term *Term) {
					top, bottom := term.getScrollingRegion()
					term._scrollDown(top, bottom, 4)
				},
			},
			{
				Screen:   "1111111111\n2222222222\n3333333333\n4444444444",
				Expected: "..........\n..........\n..........\n..........",
				Operation: func(t *testing.T, term *Term) {
					top, bottom := term.getScrollingRegion()
					term._scrollDown(top, bottom, 5)
				},
			},
		},
	}

	test.RunTests(t)
}

func TestCsiInsertLines(t *testing.T) {
	test := testTerm{
		Tests: []testCondition{
			{
				Screen:   "1234567890\nabcde",
				Expected: "1234567890\n..........\nabcde",
			},
			{
				Screen:   "1234567890\n          \nabcde",
				Expected: "1234567890\n          \n..........\nabcde",
			},
		},
		Operation: func(t *testing.T, term *Term) {
			term.csiInsertLines(1)
		},
	}

	test.RunTests(t)
}

func TestCsiInsertCharacters(t *testing.T) {
	test := testTerm{
		Tests: []testCondition{
			{
				Screen:   "1234567890",
				Expected: "123..45678",
			},
		},
		Operation: func(t *testing.T, term *Term) {
			term.curPos.X = 3
			term.csiInsertCharacters(2)
		},
	}

	test.RunTests(t)
}
