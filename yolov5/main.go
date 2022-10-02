//---------------------------------------------------------------------------------
// Function: YOLOv5 Detection
//---------------------------------------------------------------------------------
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/wimspaargaren/yolov5"
	"gocv.io/x/gocv"
)

//---------------------------------------------------------------------------------
type Program struct {
	Name        string
	yolov5Model string
	cocoPath    string
	yolonet     yolov5.Net
	objects     []yolov5.ObjectDetection
	framech     chan gocv.Mat
	fcuda       bool
	fdetect     bool
	ok          bool
}

//---------------------------------------------------------------------------------
func main() {
	var err error
	pg := Program{
		Name:    "YOLOv5 Detection",
		framech: make(chan gocv.Mat),
		fcuda:   false,
		ok:      true,
	}

	fmt.Println(pg.Name)
	flag.BoolVar(&pg.fdetect, "detect", pg.fcuda, "start YOLO detection")
	flag.BoolVar(&pg.fcuda, "cuda", pg.fcuda, "use CUDA backend")
	flag.Parse()

	assetDir := "/Users/stoney/assets/models/"
	// pg.yolov5Model = assetDir + "yolov5s.onnx"
	pg.yolov5Model = assetDir + "yolov7-tiny.onnx"
	pg.cocoPath = assetDir + "coco.names"

	netconf := yolov5.DefaultConfig()
	if pg.fcuda {
		netconf.NetBackendType = gocv.NetBackendCUDA
		netconf.NetTargetType = gocv.NetTargetCUDA
		pg.Name += " using CUDA"
	}

	pg.yolonet, err = yolov5.NewNetWithConfig(pg.yolov5Model, pg.cocoPath, netconf)
	if err != nil {
		log.Println(err)
		return
	}
	defer pg.yolonet.Close()

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

	go pg.detectObjectsByFrame()

	for pg.ok {
		if !video.Read(&frame) {
			log.Println("video.Read error")
			return
		}

		if pg.fdetect {
			select {
			case pg.framech <- frame:
				if len(pg.objects) > 0 {
					yolov5.DrawDetections(&frame, pg.objects)
				}
			default:
				// log.Println("not ready for detection. so skip it")
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
func (d *Program) detectObjectsByFrame() (err error) {
	log.Println("i.detectObjectsByFrame:")

	for d.ok {
		frame, ok := <-d.framech
		// --- for safe channel handling
		if frame.Empty() || !ok {
			continue
		}

		d.objects, err = d.yolonet.GetDetections(frame)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

//---------------------------------------------------------------------------------
