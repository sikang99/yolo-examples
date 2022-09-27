// What it does:
//
// This example uses the Window class to open an image file, and then display
// the image in a Window class.
//
// How to run:
//
// go run ./cmd/showimage/main.go /home/ron/Pictures/mcp23017.jpg
//
//go:build example
// +build example

package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\tshowimage [imgfile]")
		return
	}

	filename := os.Args[1]
	window := gocv.NewWindow("Features for " + filename)

	img := gocv.IMRead(filename, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Error reading image from:", filename)
		return
	}
	defer img.Close()

	grayImage := gocv.NewMat()
	defer grayImage.Close()

	gocv.CvtColor(img, &grayImage, gocv.ColorBGRToGray)

	destImage := gocv.NewMat()
	defer destImage.Close()

	gocv.Threshold(grayImage, &destImage, 100, 255, gocv.ThresholdBinaryInv)
	resultImage := gocv.NewMatWithSize(img.Cols(), img.Rows(), gocv.MatTypeCV8U)

	gocv.Resize(destImage, &resultImage, image.Pt(resultImage.Rows(), resultImage.Cols()), 0, 0, gocv.InterpolationCubic)
	gocv.Dilate(resultImage, &resultImage, gocv.NewMat())
	gocv.GaussianBlur(resultImage, &resultImage, image.Pt(5, 5), 0, 0, gocv.BorderWrap)

	results := gocv.FindContours(resultImage, gocv.RetrievalTree, gocv.ChainApproxSimple)
	imageForShowing := gocv.NewMatWithSize(resultImage.Rows(), resultImage.Cols(), gocv.MatChannels4)

	for i := 0; i < results.Size(); i++ {
		fmt.Println(i)
		gocv.DrawContours(&imageForShowing, results, i, color.RGBA{R: 0, G: 0, B: 255, A: 255}, 1)
		gocv.Rectangle(&imageForShowing,
			gocv.BoundingRect(results.At(i)),
			color.RGBA{R: 0, G: 255, B: 0, A: 100}, 1)
	}

	for {
		window.IMShow(imageForShowing)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
