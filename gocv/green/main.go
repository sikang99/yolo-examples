package main

import (
	"fmt"
	"log"
	"runtime"

	"gocv.io/x/gocv"
)

//---------------------------------------------------------------------------------
type Program struct {
	Name  string
	Label string
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
		Name: "Split color channel",
		ok:   true,
	}

	window := gocv.NewWindow(pg.Name)
	defer window.Close()

	window.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowAutosize)

	read := gocv.IMRead("../../data/images/julia.jpg", gocv.IMReadColor)
	fmt.Println(read.Type())

	var rgbChan1 []gocv.Mat
	bgrChan := gocv.Split(read)

	//blue:= bgrChan[0]
	green := bgrChan[1]
	//red := bgrChan[2]

	back_ch := gocv.Zeros(read.Rows(), read.Cols(), gocv.MatTypeCV8UC1)

	rgbChan1 = append(rgbChan1, back_ch)
	rgbChan1 = append(rgbChan1, green)
	rgbChan1 = append(rgbChan1, back_ch)

	image := gocv.NewMatWithSize(read.Rows(), read.Cols(), gocv.MatTypeCV8UC3)

	gocv.Merge(rgbChan1, &image)

	window.IMShow(image)
	window.WaitKey(0)
}

//---------------------------------------------------------------------------------
