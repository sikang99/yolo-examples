package main

import (
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func main() {
	// open reference image
	shiftomo := gocv.IMRead("../data/images/julia.jpg", gocv.IMReadGrayScale)
	if shiftomo.Empty() {
		log.Fatalln("image not ready")
	}
	defer shiftomo.Close()

	// open USB camera
	video, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return
	}
	defer video.Close()

	// prepare Mat for captured image and grayscaled image
	img := gocv.NewMat()
	defer img.Close()
	imgGray := gocv.NewMat()
	defer imgGray.Close()

	// create a window
	window := gocv.NewWindow("Julia SIFT")
	defer window.Close()

	// create SIFT object
	sift := gocv.NewSIFT()
	defer sift.Close()

	// find SIFT key points on the reference image
	kp1, des1 := sift.DetectAndCompute(shiftomo, gocv.NewMat())

	for {
		if !video.Read(&img) {
			log.Fatalln("video not ready")
			return
		}
		if img.Empty() {
			continue
		}
		// convert to grayscale image
		gocv.CvtColor(img, &imgGray, gocv.ColorRGBToGray)

		// find key points on the captured image
		kp2, des2 := sift.DetectAndCompute(imgGray, gocv.NewMat())

		// find match points
		bf := gocv.NewBFMatcher()
		matches := bf.KnnMatch(des1, des2, 2)
		var matchedPoints []gocv.DMatch
		for _, m := range matches {
			if len(m) > 1 {
				if m[0].Distance < 0.6*m[1].Distance {
					matchedPoints = append(matchedPoints, m[0])
				}
			}
		}
		// match color
		c1 := color.RGBA{R: 0, G: 255, B: 0, A: 0}
		// point color
		c2 := color.RGBA{R: 255, G: 0, B: 0, A: 0}
		// show matching image
		if len(matchedPoints) > 0 {
			out := gocv.NewMat()
			gocv.DrawMatches(shiftomo, kp1, imgGray, kp2, matchedPoints, &out, c1, c2, make([]byte, 0), gocv.DrawDefault)
			window.IMShow(out)
		}
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
