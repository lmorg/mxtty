package elementTable

import (
	"encoding/hex"
	"time"
)

func (el *ElementTable) setName() {
	if el.parameters.Name == "" {
		el.parameters.Name = time.Now().String()
	}

	el.name = hex.EncodeToString([]byte(el.name))
}
