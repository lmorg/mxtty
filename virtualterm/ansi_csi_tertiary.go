package virtualterm

import "log"

func lookupTertiaryCsi(term *Term, code []rune) {
	param := string(code[:len(code)-1])
	r := code[len(code)-1]
	switch r {
	case 'B':
		switch param {
		case "1":
			log.Printf("DEBUG: BEGIN 1")

		default:
			log.Printf("Tertiary CSI parameter not implemented in %s: %v [param: %s]", string(r), string(code), param)
		}

	case 'E':
		switch param {
		case "1":
			log.Printf("DEBUG: END 1")
		default:
			log.Printf("Tertiary CSI parameter not implemented in %s: %v [param: %s]", string(r), string(code), param)
		}

	default:
		log.Printf("Tertiary CSI code not implemented: %s (%s)", string(r), string(code))
	}
}
