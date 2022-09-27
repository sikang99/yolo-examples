package main

import (
	"image"
	"image/color"
	"log"
	"runtime"

	"gocv.io/x/gocv"
)

//---------------------------------------------------------------------------------
var (
	green = color.RGBA{0, 255, 0, 0}
	red   = color.RGBA{255, 0, 0, 0}
)

type Program struct {
	Name  string
	Label string
	fps   float64
	ok    bool
}

//---------------------------------------------------------------------------------
func init() {
	runtime.LockOSThread()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//---------------------------------------------------------------------------------
func main() {
	pg := Program{
		Name: "Circle Detection",
		fps:  25.0,
		ok:   true,
	}

	video, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalln(err)
	}

	// --- for testing
	video.Set(gocv.VideoCaptureFPS, pg.fps)
	fps := video.Get(gocv.VideoCaptureFPS)
	log.Println("FPS:", fps)

	window := gocv.NewWindow(pg.Name)
	defer window.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	for pg.ok {
		if !video.Read(&frame) {
			log.Fatalln("video not ready")
		}

		err = pg.detectRectangle(&frame)
		if err != nil {
			log.Fatalln(err)
		}

		window.IMShow(frame)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}

//---------------------------------------------------------------------------------
func (d *Program) detectRectangle(pframe *gocv.Mat) (err error) {
	// log.Println("i.detectCircles:")

	img := gocv.NewMat()
	defer img.Close()

	gocv.CvtColor(*pframe, &img, gocv.ColorBGRToGray)
	gocv.MedianBlur(img, &img, 5)

	circles := gocv.NewMat()
	defer circles.Close()

	gocv.HoughCirclesWithParams(
		img,
		&circles,
		gocv.HoughGradient,
		1,                     // dp
		float64(img.Rows()/8), // minDist
		// 26,
		200, // param1
		48,  // param2
		0,   // minRadius
		0,   // maxRadius
	)

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		// if circles are found
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])

			gocv.Circle(pframe, image.Pt(x, y), r, green, 2)
			gocv.Circle(pframe, image.Pt(x, y), 2, red, 3)
		}
	}
	return
}

//---------------------------------------------------------------------------------
