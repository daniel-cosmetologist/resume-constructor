package latex

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"
)

const maxPhotoSizeBytes = 2 * 1024 * 1024

var latexReplacer = strings.NewReplacer(
	`\\`, `\textbackslash{}`,
	`{`, `\{`,
	`}`, `\}`,
	`#`, `\#`,
	`$`, `\$`,
	`%`, `\%`,
	`&`, `\&`,
	`_`, `\_`,
	`~`, `\textasciitilde{}`,
	`^`, `\textasciicircum{}`,
)

// escapeLatex экранирует спецсимволы LaTeX в строке.
func escapeLatex(s string) string {
	return latexReplacer.Replace(s)
}

// processPhoto декодирует base64, приводит фото к соотношению 3:4 и
// пытается уложить размер в maxPhotoSizeBytes, возвращая байты и расширение файла.
func processPhoto(dataBase64, mimeType string) ([]byte, string, error) {
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(dataBase64))
	if err != nil {
		return nil, "", fmt.Errorf("decode base64: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		// Если не удалось декодировать как изображение, но размер допустимый — вернуть как есть.
		if len(raw) <= maxPhotoSizeBytes {
			ext := detectExt(mimeType)
			return raw, ext, nil
		}
		return nil, "", fmt.Errorf("photo is too large and cannot be decoded as image")
	}

	cropped := centerCropToRatio(img, 3, 4)

	out, err := encodeJPEGWithLimit(cropped, maxPhotoSizeBytes)
	if err != nil {
		return nil, "", err
	}

	return out, "jpg", nil
}

func detectExt(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return "bin"
	}
}

// centerCropToRatio обрезает центр изображения до заданного соотношения сторон.
func centerCropToRatio(img image.Image, num, den int) image.Image {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	if w == 0 || h == 0 {
		return img
	}

	targetW := w
	targetH := w * den / num

	if targetH > h {
		targetH = h
		targetW = h * num / den
	}

	if targetW <= 0 || targetH <= 0 {
		return img
	}

	x0 := b.Min.X + (w-targetW)/2
	y0 := b.Min.Y + (h-targetH)/2

	srcPoint := image.Point{X: x0, Y: y0}
	dstRect := image.Rect(0, 0, targetW, targetH)

	dst := image.NewRGBA(dstRect)
	draw.Draw(dst, dstRect, img, srcPoint, draw.Src)

	return dst
}

// encodeJPEGWithLimit кодирует изображение в JPEG, стараясь уложиться в limit байт.
func encodeJPEGWithLimit(img image.Image, limit int) ([]byte, error) {
	qualities := []int{90, 80, 70, 60, 50}

	var last []byte

	for i, q := range qualities {
		buf := &bytes.Buffer{}
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: q}); err != nil {
			return nil, fmt.Errorf("jpeg encode failed: %w", err)
		}
		last = buf.Bytes()
		if len(last) <= limit || i == len(qualities)-1 {
			break
		}
	}

	if len(last) == 0 {
		return nil, fmt.Errorf("encoded photo is empty")
	}

	return last, nil
}

// optional helpers for future PNG-specific handling (not used directly, но готовы при необходимости).
func decodeImage(raw []byte, mimeType string) (image.Image, error) {
	reader := bytes.NewReader(raw)
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/jpg":
		return jpeg.Decode(reader)
	case "image/png":
		return png.Decode(reader)
	default:
		return image.Decode(reader)
	}
}
