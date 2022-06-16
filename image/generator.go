package image

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"net/http"

	"github.com/AidenHadisi/MyDailyBibleBot/pkg/httpclient"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
)

//go:embed KeepCalm-Medium.ttf
var font []byte

type ImageGenerator struct {
	client httpclient.HttpClient
}

func NewImageProcessor(client httpclient.HttpClient) *ImageGenerator {
	return &ImageGenerator{
		client: client,
	}
}

func (i *ImageGenerator) Process(url, text string, fontSize int) ([]byte, error) {
	bgImage, err := i.getImage(url)
	if err != nil {
		return nil, err
	}

	//first resize the image to our max width allowed
	bgImage = resize.Resize(1200, 0, bgImage, resize.Lanczos3)

	//first draw the image
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	//then draw a semi transparent rectangle over the image to make it darker
	dc.SetColor(color.RGBA{0, 0, 0, 100})
	dc.DrawRectangle(0, 0, float64(imgWidth), float64(imgHeight))
	dc.Fill()

	//get anchor points
	x := float64(imgWidth / 2)
	y := float64((imgHeight / 2))
	maxWidth := float64(imgWidth) - 60.0
	// maxHeight := float64(imgHeight) - 100

	// //load in our font
	// dc.LoadFontFace("../../assets/KeepCalm-Medium.ttf", float64(fontSize))

	// //measure the text and make an estimate for font size
	// fontWidth, fontHeight := dc.MeasureString(text)

	// textLines := fontWidth / maxWidth
	// textHeight := textLines * fontHeight

	// estimated_font_size := float64(fontSize) * maxHeight / textHeight

	//Now draw the text inside rectangle

	f, err := truetype.Parse(font)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: float64(fontSize),
	})

	dc.SetColor(color.White)
	dc.SetFontFace(face)
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	var buf bytes.Buffer
	err = dc.EncodePNG(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i *ImageGenerator) getImage(url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
