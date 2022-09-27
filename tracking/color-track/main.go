package main

import (
	"image"
	"image/color"
	"log"

	"github.com/matipan/computer-vision/queue"
	"gocv.io/x/gocv"
)

var (
	rcolor = color.RGBA{G: 255, A: 255}
	lcolor = color.RGBA{R: 255, A: 255}

	lhsv = gocv.Scalar{Val1: 49, Val2: 89, Val3: 0}
	hhsv = gocv.Scalar{Val1: 109, Val2: 255, Val3: 255}

	size = image.Point{X: 600, Y: 600}
	blur = image.Point{X: 11, Y: 11}

	wt     = gocv.NewWindow("thersholded")
	wi     = gocv.NewWindow("images")
	img    = gocv.NewMat()
	mask   = gocv.NewMat()
	frame  = gocv.NewMat()
	hsv    = gocv.NewMat()
	kernel = gocv.NewMat()
)

func main() {
	defer close()

	wt.ResizeWindow(600, 600)
	wt.MoveWindow(0, 0)
	wi.MoveWindow(600, 0)
	wi.ResizeWindow(600, 600)

	video, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalln(err)
	}
	defer video.Close()

	queue := queue.New(40)

	for {
		if !video.Read(&img) {
			log.Fatalln("video read error")
		}

		gocv.Flip(img, &img, 1)
		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101)
		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
		gocv.Erode(mask, &mask, kernel)
		gocv.Dilate(mask, &mask, kernel)
		cnt := bestContour(mask, 2000)
		if cnt.IsNil() {
			queue.Clear()
			if imShow() {
				break
			}
			continue
		}

		rect := gocv.BoundingRect(cnt)
		gocv.Rectangle(&img, rect, rcolor, 2)
		x, y := middle(rect)
		queue.Push(image.Point{X: x, Y: y})
		queue.RangePrevious(func(c image.Point, p image.Point) {
			gocv.Line(&img, p, c, lcolor, 2)
		})

		if imShow() {
			break
		}
	}
}

func imShow() bool {
	wi.IMShow(img)
	wt.IMShow(mask)
	return wi.WaitKey(1) == 27 || wt.WaitKey(1) == 27
}

// bestContour obtains the biggest contour in the frame(provided is bigger) than the minArea.
func bestContour(frame gocv.Mat, minArea float64) gocv.PointVector {
	cnts := gocv.FindContours(frame, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var (
		bestCnt  gocv.PointVector
		bestArea = minArea
	)
	for i := 0; i < cnts.Size(); i++ {
		area := gocv.ContourArea(cnts.At(i))
		if area > bestArea {
			bestArea = area
			bestCnt = cnts.At(i)
		}
	}
	return bestCnt
}

// middle calculates the middle x and y of a rectangle.
func middle(rect image.Rectangle) (x int, y int) {
	return (rect.Max.X-rect.Min.X)/2 + rect.Min.X, (rect.Max.Y-rect.Min.Y)/2 + rect.Min.Y
}

func close() {
	defer wi.Close()
	defer wt.Close()
	defer img.Close()
	defer mask.Close()
	defer frame.Close()
	defer hsv.Close()
	defer kernel.Close()
}
