package escpos

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

func convertToRaster(img image.Image) []byte {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	widthBytes := (w + 7) / 8
	header := []byte{0x1D, 0x76, 0x30, 0x00, byte(widthBytes % 256), byte(widthBytes / 256), byte(h % 256), byte(h / 256)}
	data := make([]byte, widthBytes*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8((r + g + b) / 3 >> 8)
			if gray < 128 {
				byteIndex := y*widthBytes + x/8
				bit := 7 - (x % 8)
				data[byteIndex] |= 1 << uint(bit)
			}
		}
	}
	return append(header, data...)
}

func convertToBitImage(img image.Image) []byte {
	g := toGray(img)
	bounds := g.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	out := bytes.Buffer{}
	for y := 0; y < h; y += 8 {
		nL := byte(w % 256)
		nH := byte(w / 256)
		out.Write([]byte{0x1B, 0x2A, 0x00, nL, nH})
		for x := 0; x < w; x++ {
			var b byte
			for k := 0; k < 8; k++ {
				yy := y + k
				if yy >= h {
					break
				}
				r, g2, b2, _ := g.At(x, yy).RGBA()
				gray := uint8((r + g2 + b2) / 3 >> 8)
				if gray < 128 {
					b |= (1 << uint(7-k))
				}
			}
			out.WriteByte(b)
		}
		out.WriteByte('\n')
	}
	return out.Bytes()
}

func resizeToWidth(src image.Image, targetWidth int) image.Image {
    if targetWidth <= 0 {
        return src
    }
    b := src.Bounds()
    w := b.Dx()
    h := b.Dy()
    if targetWidth >= w {
        return src
    }
    targetHeight := int(float64(h) * float64(targetWidth) / float64(w))
    dst := image.NewGray(image.Rect(0, 0, targetWidth, targetHeight))
    for y := 0; y < targetHeight; y++ {
        sy := int(float64(y) * float64(h) / float64(targetHeight))
        for x := 0; x < targetWidth; x++ {
            sx := int(float64(x) * float64(w) / float64(targetWidth))
            dst.Set(x, y, src.At(sx, sy))
        }
    }
    return dst
}

func loadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(body)
	img, format, err := image.Decode(r)
	if err != nil {
		if format == "jpeg" {
			r = bytes.NewReader(body)
			return jpeg.Decode(r)
		}
		if format == "png" {
			r = bytes.NewReader(body)
			return png.Decode(r)
		}
		return nil, err
	}
	return img, nil
}

func toGray(src image.Image) *image.Gray {
	bounds := src.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, src, bounds.Min, draw.Src)
	return gray
}

func loadLocalImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, format, err := image.Decode(f)
	if err != nil {
		if format == "jpeg" {
			if _, err = f.Seek(0, 0); err != nil {
				return nil, err
			}
			return jpeg.Decode(f)
		}
		if format == "png" {
			if _, err = f.Seek(0, 0); err != nil {
				return nil, err
			}
			return png.Decode(f)
		}
		return nil, err
	}
	return img, nil
}
