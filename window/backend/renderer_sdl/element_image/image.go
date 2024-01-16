package elementImage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"

	"golang.org/x/image/bmp"
)

func (el *ElementImage) decode() error {
	s := el.apc.Parameter(_KEY_BASE64)
	if len(s) == 0 {
		return fmt.Errorf("no image supplied in \"base64\" parameter")
	}

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("unable to decode base64 string: %s", err.Error())
	}

	buf := bytes.NewBuffer(b)

	img, imgType, err := image.Decode(buf)
	if err != nil {
		return fmt.Errorf("unable to load image: %s", err.Error())
	}

	buf.Reset()
	err = bmp.Encode(buf, img)
	if err != nil {
		return fmt.Errorf("unable to convert %s to bitmap: %s", imgType, err.Error())
	}

	el.bmp = buf.Bytes()
	return nil
}
