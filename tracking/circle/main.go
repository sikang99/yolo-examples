package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"sort"
	"time"

	"github.com/sigtot/byggern-rest/serial"
	"github.com/sigtot/kalman"
	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
)

const y0 = 50
const xLeft = 90
const xRight = 540
const clusterThresh = 10
const boxPadding = 30

const serialName = "/dev/ttyACM0"
const serialBaud = 9600
const serialStopBits = 2

const xScale = 0.38
const xOffset = 2 * xLeft

const kickThresh = 6

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	defer webcam.Close()

	window := gocv.NewWindow("Web camera")

	A := mat.NewDense(4, 4, []float64{1, 0, 0, 0.1677, 0, 1, 0.1677, 0, 0, 0, 1, 0, 0, 0, 0, 1})
	B := mat.NewDense(4, 1, []float64{0, 0.0001406, 6 * 0.1677, 0})
	C := mat.NewDense(2, 4, []float64{1, 0, 0, 0, 0, 1, 0, 0})
	D := mat.NewDense(2, 1, []float64{0, 0})
	G := mat.NewDiagDense(4, []float64{0.2, 0.2, 0.1, 0.1})
	H := mat.NewDense(2, 2, []float64{0.1, 0.1, 0.2, 0.2})
	R := mat.NewDiagDense(2, []float64{10, 10})
	Q := mat.NewDiagDense(4, []float64{0.2, 0.2, 1, 1})

	aPriErrCovInit := mat.NewDense(4, 4, []float64{1, 0, 2, 0, 0, 1, 0, 2, 2, 0, 1, 0, 0, 2, 0, 1})
	aPriStateEstInit := mat.NewVecDense(4, []float64{300, 200, 0, 0})
	input := mat.NewVecDense(1, []float64{-4})
	outputInit := mat.NewVecDense(2, []float64{300, 200})

	f := kalman.NewFilter(A, B, C, D, H, G, R, Q, aPriErrCovInit, aPriStateEstInit, input, outputInit)

	conn, err := serial.CreateConnection(
		serialName,
		serialBaud,
		serialStopBits)
	if err != nil {
		log.Println(err)
		return
	}

	ticker := time.NewTicker(16 * time.Millisecond)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			drawPredictImage(window, webcam, &f, conn)
			//sendPredictPos(webcam, &f, conn)
		case <-quit:
			ticker.Stop()
			return
		}
	}
	err = conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func findGameBox(webcam *gocv.VideoCapture) {
	window := gocv.NewWindow("Game box Canny")
	img := gocv.NewMat()
	defer img.Close()

	window2 := gocv.NewWindow("Gray scale image")

	if ok := webcam.Read(&img); !ok {
		log.Fatal("Webcam closed")
	}

	if img.Empty() {
		log.Println("Warning: Read empty image when trying to find game box")
		return
	}

	grayImg := gocv.NewMat()
	gocv.CvtColor(img, &grayImg, gocv.ColorRGBAToGray)

	gocv.MedianBlur(grayImg, &grayImg, 3)
	canny := gocv.NewMat()
	defer canny.Close()

	gocv.Canny(grayImg, &canny, 3, 3)

	erodeKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(1, 3))
	gocv.Erode(canny, &canny, erodeKernel)

	dilateKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	gocv.Dilate(canny, &canny, dilateKernel)

	lines := gocv.NewMat()
	gocv.HoughLinesP(canny, &lines, 1, 3.14/180, 80)

	var xValues []int
	for i := 0; i < lines.Rows(); i++ {
		pt1 := image.Pt(int(lines.GetVeciAt(i, 0)[0]), int(lines.GetVeciAt(i, 0)[1]))
		pt2 := image.Pt(int(lines.GetVeciAt(i, 0)[2]), int(lines.GetVeciAt(i, 0)[3]))
		if math.Sqrt(math.Pow(float64(pt2.X-pt1.X), 2)+math.Pow(float64(pt2.Y-pt1.Y), 2)) > 30 {
			x0 := findIntersection(pt1, pt2, img.Rows()-y0)
			xValues = append(xValues, x0)
			gocv.Line(&img, pt1, image.Pt(x0, y0), color.RGBA{0, 255, 0, 50}, 2)
		}
	}

	xClusters := []int{0}
	var clusterValues []int
	if len(xValues) > 0 {
		sort.Ints(xValues)
		for i := 0; i < len(xValues); i++ {
			if xValues[i] > xClusters[len(xClusters)-1]+clusterThresh {
				xClusters = append(xClusters, xValues[i])
				clusterValues = []int{}
			} else {
				clusterValues = append(clusterValues, xValues[i])
				xClusters[len(xClusters)-1] = average(clusterValues)
			}
		}
	}

	blue := color.RGBA{0, 0, 255, 0}
	for _, x := range xClusters {
		gocv.Circle(&img, image.Pt(x, y0), 2, blue, 2)
	}

	window.IMShow(canny)
	window2.IMShow(img)
	window.WaitKey(1)
}

func average(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum / len(values)
}

func findIntersection(pt1 image.Point, pt2 image.Point, y0 int) int {
	return pt1.X + (pt1.X-pt2.X)/(pt1.Y-pt2.Y)*(y0-pt1.Y)
}

func sendPredictPos(webcam *gocv.VideoCapture, f *kalman.Filter, conn serial.Connection) {
	err := conn.Write(fmt.Sprintf("{servo=%d}", 25))
	if err != nil {
		fmt.Print(err)
	}
	img := gocv.NewMat()
	defer img.Close()
	if ok := webcam.Read(&img); !ok {
		log.Fatal("Webcam closed")
	}

	if img.Empty() {
		log.Println("Warning: Read empty image")
		return
	}

	cimg := gocv.NewMat()
	defer cimg.Close()

	mask := gocv.NewMat()
	defer mask.Close()

	hsvImg := gocv.NewMat()
	defer hsvImg.Close()

	rot := float64(10)
	hlsBlueBelow := gocv.NewScalar(rot, 80, 80, 0)
	hlsBlueAbove := gocv.NewScalar(rot+15, 255, 255, 0)

	gocv.CvtColor(img, &hsvImg, gocv.ColorRGBAToBGR)
	gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorBGRToHLS)
	gocv.InRangeWithScalar(hsvImg, hlsBlueBelow, hlsBlueAbove, &mask)

	cnts := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxNone)
	largestContour := 0
	largestContourCount := 0
	for i := 0; i < cnts.Size(); i++ {
		length := cnts.At(i).Size()
		if length > largestContourCount {
			largestContour = i
			largestContourCount = length
		}
	}

	if !cnts.IsNil() {
		rect := gocv.MinAreaRect(cnts.At(largestContour))

		x := rect.Center.X
		y := rect.Center.Y

		output := mat.NewVecDense(2, []float64{float64(x), float64(y)})
		f.AddOutput(output)

		//kicked := false
		// Draw 100 predicted positions
		lastPredY := f.APostStateEst(f.CurrentK()).At(1, 0)
		for i := 0; i < 100; i++ {
			aPostStateEst := f.APostStateEst(f.CurrentK() + i)
			predX := aPostStateEst.At(0, 0)
			predY := aPostStateEst.At(1, 0)

			// Ball crosses bottom line at this i
			if lastPredY > y0 && predY < y0 {
				xPos := confinedToRange(int(predX), xLeft, xRight)

				cartReference := int(float64(img.Rows()-xPos+xOffset) * xScale)
				err := conn.Write(fmt.Sprintf("{motor=%d}", cartReference))
				if err != nil {
					log.Println(err)
				}
				/*
					                if i < kickThresh && kicked == false {
										err = conn.Write("{kick}")
										kicked = true
									}
				*/
			}
			lastPredY = predY
		}
	}
}

func drawPredictImage(window *gocv.Window, webcam *gocv.VideoCapture, f *kalman.Filter, conn serial.Connection) {
	img := gocv.NewMat()
	defer img.Close()
	if ok := webcam.Read(&img); !ok {
		log.Fatal("Webcam closed")
	}

	if img.Empty() {
		log.Println("Warning: Read empty image")
		return
	}

	cimg := gocv.NewMat()
	defer cimg.Close()

	mask := gocv.NewMat()
	defer mask.Close()

	hsvImg := gocv.NewMat()
	defer hsvImg.Close()

	rot := float64(10)
	hlsBlueBelow := gocv.NewScalar(rot, 80, 80, 0)
	hlsBlueAbove := gocv.NewScalar(rot+15, 255, 255, 0)

	gocv.CvtColor(img, &hsvImg, gocv.ColorRGBAToBGR)
	gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorBGRToHLS)
	gocv.InRangeWithScalar(hsvImg, hlsBlueBelow, hlsBlueAbove, &mask)

	cnts := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxNone)
	largestContour := 0
	largestContourCount := 0
	for i := 0; i < cnts.Size(); i++ {
		length := cnts.At(i).Size()
		if length > largestContourCount {
			largestContour = i
			largestContourCount = length
		}
	}

	red := color.RGBA{255, 0, 0, 0}
	green := color.RGBA{0, 255, 0, 0}
	blue := color.RGBA{0, 0, 255, 0}
	yellow := color.RGBA{255, 255, 0, 0}

	if !cnts.IsNil() {
		rect := gocv.MinAreaRect(cnts.At(largestContour))

		x := rect.Center.X
		y := rect.Center.Y

		gocv.Circle(&img, image.Pt(x, y), 7, red, 13)

		output := mat.NewVecDense(2, []float64{float64(x), float64(y)})
		f.AddOutput(output)

		// Draw 100 predicted positions
		lastPredY := f.APostStateEst(f.CurrentK()).At(1, 0)
		for i := 0; i < 100; i++ {
			aPostStateEst := f.APostStateEst(f.CurrentK() + i)
			predX := aPostStateEst.At(0, 0)
			predY := aPostStateEst.At(1, 0)
			gocv.Circle(&img, image.Pt(int(predX), int(predY)), 2, green, 2)

			// Ball crosses bottom line at this k
			if lastPredY > y0 && predY < y0 {
				xPos := confinedToRange(int(predX), xLeft, xRight)

				gocv.Circle(&img, image.Pt(xPos, (int(predY)+int(lastPredY))/2), 5, yellow, 5)

				/*
									cartReference := int(float64(img.Rows() - xPos + 250) * xScale)
					                err := conn.Write(fmt.Sprintf("{motor=%d}", cartReference))
					                if err != nil {
					                    log.Println(err)
					                }
				*/

			}
			lastPredY = predY
		}
	}

	// Draw bottom line
	gocv.Line(&img, image.Pt(xLeft-boxPadding, y0), image.Pt(xRight+boxPadding, y0), blue, 2)
	gocv.Line(&img, image.Pt(xLeft, y0), image.Pt(xRight, y0), green, 2)

	window.IMShow(img)
	window.WaitKey(1)
}

func confinedToRange(value, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
func printLocation(webcam *gocv.VideoCapture) {
	img := gocv.NewMat()
	if ok := webcam.Read(&img); !ok {
		log.Fatal("Webcam closed")
	}

	if img.Empty() {
		log.Println("Warning: Read empty image")
		return
	}

	gocv.CvtColor(img, &img, gocv.ColorRGBToGray)

	circles := gocv.NewMat()
	defer circles.Close()

	gocv.HoughCirclesWithParams(
		img,
		&circles,
		gocv.HoughStandard,
		1,                     // dp
		float64(img.Rows()/8), // minDist
		75,                    // param1
		10,                    // param2
		25,                    // minRadius
		28,                    // maxRadius
	)

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		// if circles are found
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])
			fmt.Printf("pos=(%d, %d) r=%d\n", x, y, r)
		}
	}
}

func showImg(webcam *gocv.VideoCapture, window *gocv.Window) {
	img := gocv.NewMat()
	if ok := webcam.Read(&img); !ok {
		log.Fatal("Webcam closed")
	}

	if img.Empty() {
		log.Println("Warning: Read empty image")
		return
	}

	cimg := gocv.NewMat()
	defer cimg.Close()

	gocv.CvtColor(img, &img, gocv.ColorRGBToGray)
	gocv.CvtColor(img, &cimg, gocv.ColorGrayToBGR)

	circles := gocv.NewMat()
	defer circles.Close()

	gocv.HoughCirclesWithParams(
		img,
		&circles,
		gocv.HoughGradient,
		1,                     // dp
		float64(img.Rows()/8), // minDist
		75,                    // param1
		25,                    // param2
		25,                    // minRadius
		28,                    // maxRadius
	)

	blue := color.RGBA{0, 0, 255, 0}
	red := color.RGBA{255, 0, 0, 0}

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		// if circles are found
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])

			gocv.Circle(&cimg, image.Pt(x, y), r, blue, 2)
			gocv.Circle(&cimg, image.Pt(x, y), 2, red, 3)
		}
	}

	window.IMShow(cimg)
	window.WaitKey(1)
}
