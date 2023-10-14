package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	exifremove "github.com/scottleedavis/go-exif-remove"
	"golang.org/x/image/draw"
)

func main() {
	err := Run("icon_1.png", 1000, 1000)
	if err != nil {
		panic(err)
	}

	err = Run("icon_2.png", 1000, 1000)
	if err != nil {
		panic(err)
	}

	err = Run("icon_3.jpeg", 1000, 1000)
	if err != nil {
		panic(err)
	}

}

func Run(fileName string, width int, height int) error {
	path := filepath.Join("images", fileName)
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	img, t, err := image.Decode(file)
	if err != nil {
		return err
	}

	fmt.Printf("type: %s\n", t)

	rctSrc := img.Bounds()

	imgDst := image.NewRGBA(
		image.Rect(
			0,
			0,
			width,
			height,
		),
	)
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), img, rctSrc, draw.Over, nil)

	buff := bytes.NewBuffer([]byte{})
	writer := io.Writer(buff)

	err = png.Encode(writer, imgDst)
	if err != nil {
		return err
	}

	fileNameNoExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	exportPath := filepath.Join("dist", fileNameNoExt+".png")

	// Remove EXIF
	noExifBytes, err := exifremove.Remove(buff.Bytes())
	if err != nil {
		return err
	}

	err = os.WriteFile(exportPath, noExifBytes, 0666)
	if err != nil {
		return err
	}

	return nil
}
