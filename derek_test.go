package derek

import (
	"crypto/sha1"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"strconv"
	"testing"
)

func TestProcess(t *testing.T) {
	savePng, _ := strconv.ParseBool(os.Getenv("SAVE_PNG"))

	src, err := loadImage("derekify/sample.in.jpg")
	if err != nil {
		t.Fatal(err)
	}

	want := "0bbf280772c08884cb1ffad785da93d2e40d2357"

	t.Run("test.png", func(t *testing.T) {
		dst := Process(src, 8, 90)
		sum := hash(toGray(dst))
		if r := hex.EncodeToString(sum); r != want {
			t.Error("test.png hash:", r, " want:", want)
		}
		if savePng {
			if err := savePNG("test.png", dst); err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("gray.png", func(t *testing.T) {
		dst := ProcessGray(toGray(src), 8, 90)
		sum := hash(dst)
		if r := hex.EncodeToString(sum); r != want {
			t.Error("gray.png hash:", r, " want:", want)
		}
		if savePng {
			if err := savePNG("gray.png", dst); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func BenchmarkProcess(b *testing.B) {
	src, err := loadImage("derekify/sample.in.jpg")
	if err != nil {
		b.Fatal(err)
	}

	b.Run("src", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Process(src, 8, 90)
		}
	})

	b.Run("toGray", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ProcessGray(toGray(src), 8, 90)
		}
	})

	gray := toGray(src)

	b.Run("gray", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Process(gray, 8, 90)
		}
	})

	b.Run("Gray", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ProcessGray(gray, 8, 90)
		}
	})
}

func loadImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	i, _, err := image.Decode(f)
	return i, err
}

func toGray(src image.Image) *image.Gray {
	return NewGray(src)
}

func savePNG(file string, img image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func hash(gray *image.Gray) []byte {
	h := sha1.New()
	h.Write(gray.Pix)
	return h.Sum(nil)
}
