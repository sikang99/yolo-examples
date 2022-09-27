//---------------------------------------------------------------------------------
// Function: Barcode Detection using zar
//---------------------------------------------------------------------------------
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/bieber/barcode"
	"gocv.io/x/gocv"
)

//---------------------------------------------------------------------------------
type Program struct {
	Name    string
	framech chan gocv.Mat
	fdetect bool
	ok      bool
}

//---------------------------------------------------------------------------------
func main() {
	var err error
	pg := Program{
		Name:    "Barcode Detection",
		framech: make(chan gocv.Mat),
		ok:      true,
	}

	fmt.Println(pg.Name)
	flag.BoolVar(&pg.fdetect, "detect", pg.fdetect, "start barcode detection")
	flag.Parse()

	video, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Println(err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(pg.Name)
	defer window.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	go pg.detectBarcodeByFrame()

	for pg.ok {
		if !video.Read(&frame) {
			log.Println("video.Read error")
			return
		}

		if pg.fdetect {
			select {
			case pg.framech <- frame:
				// log.Println("detect and update boxes")
			default:
				// log.Println("not ready for detection. so skip")
			}
		}

		window.IMShow(frame)
		switch window.WaitKey(1) {
		case 27: // Esc key
			pg.ok = false
			return
		case 'd': // detect toggle
			pg.fdetect = !pg.fdetect
		}
	}
}

//---------------------------------------------------------------------------------
func (d *Program) detectBarcodeByFrame() (err error) {
	log.Println("i.detectBarcodeByFrame:")

	textColor := color.RGBA{255, 0, 0, 0} // red
	dotColor := color.RGBA{0, 255, 0, 0}  // green

	for d.ok {
		frame, ok := <-d.framech
		// --- for safe channel handling
		if frame.Empty() || !ok {
			continue
		}

		scanner := barcode.NewScanner().SetEnabledAll(true)

		img, _ := frame.ToImage()
		src := barcode.NewImage(img)
		symbols, _ := scanner.ScanImage(src)
		log.Println(symbols)

		for _, s := range symbols {
			data := s.Data
			points := s.Boundary // Data points that zbar returns

			x0 := points[0].X
			y0 := points[0].Y

			size := gocv.GetTextSize(data, gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(x0-size.X, y0-size.Y)
			gocv.PutText(&frame, data, pt, gocv.FontHersheyPlain, 1.2, textColor, 2)

			for _, p := range points {
				x0 := p.X
				y0 := p.Y
				pt := image.Pt(x0, y0)
				gocv.PutText(&frame, ".", pt, gocv.FontHersheyPlain, 1.2, dotColor, 2)
			}
		}

	}
	return
}

//---------------------------------------------------------------------------------
