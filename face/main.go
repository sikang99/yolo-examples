package main

// Face detection, based on Viola-Jones method

import (
	"fmt"
	"image/color"

	"gocv.io/x/gocv"
)

func main() {

	webCam, err := gocv.VideoCaptureDevice(0)
	if nil != err {
		fmt.Println("ErrIno:VideoCaptureDevice! ", err)
		return
	}
	defer webCam.Close()

	window := gocv.NewWindow("viola-jones")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	width := int(webCam.Get(gocv.VideoCaptureFrameWidth))
	height := int(webCam.Get(gocv.VideoCaptureFrameHeight))

	fmt.Printf("video resolution: [%v], [%v]\n", width, height)

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if ok := classifier.Load("../../face/haar/haarcascade_frontalface_default.xml"); !ok {
		fmt.Printf("error load xml!")
		return
	}

	blue := color.RGBA{0, 0, 255, 0}

	for {
		if ok := webCam.Read(&img); !ok {
			fmt.Println("ErrIno:Read! ")
			return
		}

		if img.Empty() {
			continue
		}

		rects := classifier.DetectMultiScale(img)

		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)
		}

		window.IMShow(img)

		if 27 == window.WaitKey(1) {
			break
		}
	}
}
