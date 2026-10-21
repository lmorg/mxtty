package prompts

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/lmorg/ttyphoon/ai/agent/aitypes"
)

func TestImageAttachmentsBecomeMultimodalUserParts(t *testing.T) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	images := []aitypes.ImageAttachment{{MIMEType: "image/png", Base64: base64.StdEncoding.EncodeToString(buffer.Bytes())}}
	message := buildUserMessageWithImages("Describe this image", images)
	if message.Role != schema.User || len(message.UserInputMultiContent) != 2 {
		t.Fatalf("not a multimodal user message: %+v", message)
	}
	parts := message.UserInputMultiContent
	if parts[0].Text != "Describe this image" || parts[1].Type != schema.ChatMessagePartTypeImageURL || parts[1].Image == nil {
		t.Fatalf("missing text or image parts: %+v", parts)
	}
	if parts[1].Image.MIMEType != "image/png" || parts[1].Image.Base64Data == nil || *parts[1].Image.Base64Data != images[0].Base64 {
		t.Fatal("model did not receive the image bytes")
	}
	if plain := buildUserMessageWithImages("Text only", nil); plain.Content != "Text only" || len(plain.UserInputMultiContent) != 0 {
		t.Fatal("text-only prompts changed")
	}
}
