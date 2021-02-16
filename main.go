package main

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
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
	X        int
	Y        int
	Text     string
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	imageName := strconv.Itoa(rand.Int()) + ".png"
	log.Println("Image name: ", imageName)

	background := &Thumbnail{
		Name:   imageName,
		Width:  1280,
		Height: 720,
		Color:  color.RGBA{R: 100, G: 200, B: 200, A: 0xff},
	}

	caption := &Caption{
		FontSize: 20,
		X:        20,
		Y:        20,
		Text:     "Use filepath.Join to create the path from the directory dir and the file name.\n\n",
	}

	img, err := background.generateImageBg()
	if err != nil {
		log.Fatal(err)
	}

	err = caption.GenerateCaption(img, background)
	if err != nil {
		log.Fatal(err)
	}

}

func (i *Thumbnail) generateImageBg() (*image.RGBA, error) {
	// Create a colored image of the given width and height.
	img := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))

	// Set color for each pixel
	for x := 0; x < i.Width; x++ {
		for y := 0; y < i.Height; y++ {
			img.Set(x, y, i.Color)
		}
	}

	log.Println("Image background success generated")

	return img, nil
}

func (c *Caption) GenerateCaption(img *image.RGBA, i *Thumbnail) error {
	point := fixed.Point26_6{X: fixed.Int26_6(c.X * 64), Y: fixed.Int26_6(c.Y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(c.Text)

	err := i.SaveToDisk(i.Name, img)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

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
