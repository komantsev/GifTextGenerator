package giftextgenerator

import (
	"bufio"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func Generate(dst io.Writer, text string, textFont string, textSize float64, textColor color.RGBA, bgColor color.RGBA) {

	dpi := float64(300)

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(textFont)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(textSize)

	height := int(c.PointToFixed(textSize)>>6)
	imageHeight := int(float64(height) * 1.2)

	imageWidth := len(text) * height * 2;

	// Initialize the context.
	fg, bg := image.NewUniform(textColor), image.NewUniform(bgColor)

	rgba := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)

	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(imageHeight, int(float64(height) * 0.9))

	_, err = c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
		return
	}

	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    textSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}

	textWidth := d.MeasureString(text).Round()
	imageWidth = textWidth + int(float64(imageHeight) * 2.3)

	//SaveToFile(rgba, "out.png");

	step := imageHeight / 5
	size := imageHeight

	outGif := gif.GIF{}

	for i := 0; i < imageWidth - size; i += step {
		croppedImage := CropImage(rgba, size, i)

		//SaveToFile(croppedImage, "croped" + strconv.Itoa(i) + ".png");

		bounds := croppedImage.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.WebSafe)
		draw.Draw(palettedImage, palettedImage.Rect, croppedImage, bounds.Min, draw.Src)

		//SaveToFile(palettedImage, "croped" + strconv.Itoa(i) + ".gif");

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 8)
	}

	gif.EncodeAll(dst, &outGif)
}

func CropImage(rgba *image.RGBA, size int, offset int) image.Image {
	croppedImg := imaging.Crop(rgba, image.Rect(offset, 0, size + offset, size))
	return croppedImg
}

func SaveToFile(rgba image.Image, filename string) {
	outFile, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote " + filename + " OK.")
}