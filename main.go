package main

import (
	"flag"
	"github.com/golang/freetype"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Thumbnail need canvas, so here it's
type Thumbnail struct {
	Name   string
	Width  int
	Height int
	Color  color.RGBA
}

type Caption struct {
	FontSize float64
	FontPath string
	FontType string
	Color    image.Image
	DPI      float64
	Spacing  float64
	X        int
	Y        int
	Text     string
}

func main() {
	text := flag.String("text", "Default text", "This text will be convert to image")
	fontsize := flag.Float64("fontsize", 25, "Define font size")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	imageName := strconv.Itoa(rand.Int()) + ".png"
	log.Println("Image name: ", imageName)

	background := &Thumbnail{
		Name:   imageName,
		Width:  1280,
		Height: 720,
		Color:  color.RGBA{R: 242, G: 239, B: 228, A: 0xff},
	}

	caption := &Caption{
		FontSize: *fontsize,
		X:        20,
		Y:        background.Height / 4 * 3,
		Text:     *text,
		FontPath: "./font/nutino-sans/",
		FontType: "NunitoSans-SemiBold.ttf",
		Color:    image.Black,
		DPI:      72,
		Spacing:  1.5,
	}
	img, err := background.generateImageBg()
	if err != nil {
		log.Fatal(err)
	}

	img, err = caption.GenerateCaption(img, background)
	if err != nil {
		log.Fatal(err)
	}

	err = background.SaveToDisk(background.Name, img)
	if err != nil {
		log.Fatal(err)

	}
	log.Println("Image background success generated")

}

func (i *Thumbnail) generateImageBg() (*image.RGBA, error) {
	// Create a colored image of the given width and height.
	img := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
	// Draw image with given color
	draw.Draw(img, img.Bounds(), &image.Uniform{i.Color}, image.Point{X: i.Width, Y: i.Height}, draw.Src)
	return img, nil
}

func (c *Caption) GenerateCaption(img *image.RGBA, i *Thumbnail) (*image.RGBA, error) {
	//create freetype context
	f := freetype.NewContext()

	//read given font
	fontBytes, err := ioutil.ReadFile(c.FontPath + c.FontType)
	if err != nil {
		return nil, err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// set free type properties
	f.SetDPI(c.DPI)
	f.SetFont(font)
	f.SetFontSize(c.FontSize)
	f.SetClip(img.Bounds())
	f.SetDst(img)
	f.SetSrc(c.Color)

	// Set caption position
	pt := freetype.Pt(c.X, c.Y+int(f.PointToFixed(c.FontSize)>>6))

	// Draw captionn text
	_, err = f.DrawString(c.Text, pt)

	if err != nil {
		log.Println(err)
		return img, nil
	}

	pt.Y += f.PointToFixed(c.FontSize * c.Spacing)

	return img, nil
}

func (i *Thumbnail) SaveToDisk(name string, img *image.RGBA) error {
	const path = "images"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	file, err := os.Create(filepath.Join(path, filepath.Base(name)))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	// Make image
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
