package elementTable

import (
	"encoding/hex"
	"time"
)

func (el *ElementTable) setName() {
	el.name = el.apc.Parameter(_KEY_TABLE_NAME)
	if el.name == "" {
		el.name = time.Now().String()
	}

	el.name = hex.EncodeToString([]byte(el.name))
}
