package elementImage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/bmp"
)

func (el *ElementImage) decode() error {
	var (
		b   []byte
		err error
	)
	s := el.apc.Parameter(_KEY_FILENAME)
	if s != "" {
		b, err = _loadImage(s)
	} else {
		s = el.apc.Parameter(_KEY_BASE64)
		if s == "" {
			return fmt.Errorf("no image supplied in \"base64\" nor \"filename\" parameters")
		}
		b, err = _decodeImage(s)
	}
	if err != nil {
		return err
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

func _loadImage(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open image: %s", err.Error())
	}

	return io.ReadAll(f)
}

func _decodeImage(b64 string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64 string: %s", err.Error())
	}
	return b, nil
}
