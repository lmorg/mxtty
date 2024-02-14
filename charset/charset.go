package charset

/*
	Reference documentation used:
	- box drawing: https://en.wikipedia.org/wiki/Box-drawing_character#Unix,_CP/M,_BBS
	- character sets: https://en.wikipedia.org/wiki/National_Replacement_Character_Set
	- ascii table: https://upload.wikimedia.org/wikipedia/commons/thumb/1/1b/ASCII-Table-wide.svg/1280px-ASCII-Table-wide.svg.png
	- character codes: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Controls-beginning-with-ESC
*/

// DEC Special Character and Line Drawing Set, VT100.
var DecSpecialChar = map[rune]rune{
	'`': '◆',
	'a': '▒',
	'b': '␉',
	'c': '␌',
	'd': '␍',
	'e': '␊',
	'f': '°',
	'g': '±',
	'h': '␤',
	'i': '␋',
	'j': '┘',
	'k': '┐',
	'l': '┌',
	'm': '└',
	'n': '┼',
	'o': '⎺',
	'p': '⎻',
	'q': '─',
	'r': '⎼',
	's': '⎽',
	't': '├',
	'u': '┤',
	'v': '┴',
	'w': '┬',
	'x': '│',
	'y': '≤',
	'z': '≥',
	'{': 'π',
	'|': '≠',
	'}': '£',
	'~': '·',
}

var UnitedKingdom = map[rune]rune{
	'#': '£',
}

var Finnish = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	# 	@ 	Ä 	Ö 	Å 	Ü 	_ 	é 	ä 	ö 	å 	ü
	'[':  'Ä',
	'\\': 'Ö',
	']':  'Å',
	'^':  'Ü',
	'`':  'é',
	'{':  'ä',
	'|':  'ö',
	'}':  'å',
	'~':  'ü',
}

var Swedish = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	# 	É 	Ä 	Ö 	Å 	Ü 	_ 	é 	ä 	ö 	å 	ü
	'@':  'É',
	'[':  'Ä',
	'\\': 'Ö',
	']':  'Å',
	'^':  'Ü',
	'`':  'é',
	'{':  'ä',
	'|':  'ö',
	'}':  'å',
	'~':  'ü',
}

var German = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	# 	§ 	Ä 	Ö 	Ü 	^ 	_ 	` 	ä 	ö 	ü 	ß
	'@':  '§',
	'[':  'Ä',
	'\\': 'Ö',
	']':  'Ü',
	'{':  'ä',
	'|':  'ö',
	'}':  'ü',
	'~':  'ß',
}

var French = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	£ 	à 	° 	ç 	§ 	^ 	_ 	` 	é 	ù 	è 	¨
	'#':  '£',
	'@':  'à',
	'[':  '°',
	'\\': 'ç',
	']':  '§',
	'{':  'é',
	'|':  'ù',
	'}':  'è',
	'~':  '¨',
}

var FrenchCanadian = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	# 	à 	â 	ç 	ê 	î 	_ 	ô 	é 	ù 	è 	û
	'@':  'à',
	'[':  'â',
	'\\': 'ç',
	']':  'ê',
	'^':  'î',
	'`':  'ô',
	'{':  'é',
	'|':  'ù',
	'}':  'è',
	'~':  'û',
}

var Italian = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	£ 	§ 	° 	ç 	é 	^ 	_ 	ù 	à 	ò 	è 	ì
}

var Spanish = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	£ 	§ 	¡ 	Ñ 	¿ 	^ 	_ 	` 	˚ 	ñ 	ç 	~
}

var Dutch = map[rune]rune{
	//	# 	@ 	[ 	\ 	] 	^ 	_ 	` 	{ 	| 	} 	~
	//	£ 	¾ 	ĳ 	½ 	| 	^ 	_ 	` 	¨ 	ƒ 	¼ 	´
}
