package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	originalImageName := flag.String("image", "", "Original Image File")
	rRate := flag.Float64("rrate", 1.0, "Red multiply rate")
	gRate := flag.Float64("grate", 0.9, "Green multiply rate")
	bRate := flag.Float64("brate", 0.5, "Blue multiply rate")
	hShift := flag.Int("hshift", 0.0, "H Shift (degree)")
	vignette := flag.String("vignette", "none", "Vignette (none/low/high)")
	flag.Parse()

	// adjust for 360 degree
	*hShift = *hShift * 256 / 360

	img := gocv.IMRead(*originalImageName, gocv.IMReadColor)
	width := img.Cols()
	height := img.Rows()
	channels := img.Channels()
	fmt.Printf("image %s was loaded. widht = %d, height = %d\n", *originalImageName, width, height)

	startTime := time.Now()

	fmt.Println("Converting to HSV")
	imgHSV := gocv.NewMat()
	gocv.CvtColor(img, &imgHSV, gocv.ColorBGRToHSVFull)

	// Shift HUE
	ptrHSV, _ := imgHSV.DataPtrUint8()
	fmt.Println("Shifting Hue")
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ptrHSV[(y*width+x)*channels] = ptrHSV[(y*width+x)*channels] + uint8(*hShift)
		}
	}

	fmt.Println("Converting back to BGR")
	imgCrossProcess := gocv.NewMat()
	gocv.CvtColor(imgHSV, &imgCrossProcess, gocv.ColorHSVToBGRFull)

	// Apply RGB factor
	ptrBGR, _ := imgCrossProcess.DataPtrUint8()
	fmt.Println("Applying RGB factors")
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ptrBGR[(y*width+x)*channels+0] = uint8(float64(ptrBGR[(y*width+x)*channels+0]) * *bRate)
			ptrBGR[(y*width+x)*channels+1] = uint8(float64(ptrBGR[(y*width+x)*channels+1]) * *gRate)
			ptrBGR[(y*width+x)*channels+2] = uint8(float64(ptrBGR[(y*width+x)*channels+2]) * *rRate)
		}
	}

	// vignette
	if *vignette != "none" {
		vFactor := 0.7
		if *vignette == "high" {
			vFactor = 1.0
		}
		centerX := float64(width / 2)
		centerY := float64(height / 2)
		fmt.Println("Applying vignette map")
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				nx := (float64(x) - centerX) / centerX
				ny := (float64(y) - centerY) / centerY
				dist := math.Sqrt(nx*nx+ny*ny) / math.Sqrt(2)
				ptrBGR[(y*width+x)*channels+0] = uint8(float64(ptrBGR[(y*width+x)*channels+0]) * (1 - dist*vFactor))
				ptrBGR[(y*width+x)*channels+1] = uint8(float64(ptrBGR[(y*width+x)*channels+1]) * (1 - dist*vFactor))
				ptrBGR[(y*width+x)*channels+2] = uint8(float64(ptrBGR[(y*width+x)*channels+2]) * (1 - dist*vFactor))
			}
		}
	}

	endTime := time.Now()
	fmt.Printf("Process time = %s", endTime.Sub(startTime).String())

	windowOrginal := gocv.NewWindow("Orginal")
	windowCrossProcess := gocv.NewWindow("Cross Process")
	windowOrginal.IMShow(img)
	windowCrossProcess.IMShow(imgCrossProcess)

	windowOrginal.WaitKey(0)
}

//---------------------------------------------------------------------------------
