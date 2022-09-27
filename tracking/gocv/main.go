//=======================================================================================
// Name: GoCV Object Tracking Examples
// Author: Stony Kang, sikang@teamgrit.kr
// Copyright: TeamGRIT, 2022
//=======================================================================================

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

//---------------------------------------------------------------------------------
var (
	red   = color.RGBA{R: 255}
	green = color.RGBA{G: 255}
	blue  = color.RGBA{B: 255}
	black = color.RGBA{0, 0, 0, 0}
	white = color.RGBA{255, 255, 255, 0}
)

type Program struct {
	Name        string
	Device      string
	TrackerType string
	Verbose     bool
	// --- internal variables
	tracker gocv.Tracker
}

//---------------------------------------------------------------------------------
func init() {
	// runtime.LockOSThread()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//---------------------------------------------------------------------------------
func main() {
	pg := Program{
		Device:      "0",
		TrackerType: "csrt",
		Verbose:     true,
	}

	fmt.Println(pg.Name)
	flag.StringVar(&pg.Device, "device", pg.Device, "device to use")
	flag.StringVar(&pg.TrackerType, "tracker", pg.TrackerType, "tracker type [kcf|csrt]")
	flag.BoolVar(&pg.Verbose, "verbose", pg.Verbose, "verbose mode to display logs")
	flag.Parse()

	// open webcam
	webcam, err := gocv.OpenVideoCapture(pg.Device)
	if err != nil {
		log.Println("open video capture device:", pg.Device)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Tracking Method: " + pg.TrackerType)
	defer window.Close()

	// create a tracker instance
	// (one of MIL, KCF, TLD, MedianFlow, Boosting, MOSSE or CSRT, GOTURN)
	switch pg.TrackerType {
	case "kcf":
		pg.tracker = contrib.NewTrackerKCF()
	case "csrt":
		pg.tracker = contrib.NewTrackerCSRT()
	default:
		log.Println("unknown tracker:", pg.TrackerType)
		return
	}
	defer pg.tracker.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// read an initial image
	if ok := webcam.Read(&img); !ok {
		log.Println("cannot read device:", pg.Device)
		return
	}

	gocv.PutText(&img, "Select a ROI with mouse dragging and press SPACE or ENTER",
		image.Point{20, 40}, gocv.FontHersheyPlain, 1.0, black, 2)

	// let the user mark a ROI to track in the window
	rect := window.SelectROI(img)
	if rect.Max.X == 0 {
		log.Println("user cancelled roi selection")
		return
	}

	// initialize the tracker with the image & the selected roi
	init := pg.tracker.Init(img, rect)
	if !init {
		log.Println("could not initialize the tracker")
		return
	}

	for {
		if ok := webcam.Read(&img); !ok {
			log.Println("closed device:", pg.Device)
			return
		}
		if img.Empty() {
			continue
		}

		// update the roi
		rect, _ := pg.tracker.Update(img)

		// draw it
		gocv.Rectangle(&img, rect, blue, 3)

		// show the image in the window, and wait 10 millisecond
		window.IMShow(img)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}

//=======================================================================================
