//---------------------------------------------------------------------------------
// Function: YOLOv3 Detection
//---------------------------------------------------------------------------------
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/wimspaargaren/yolov3"
	"gocv.io/x/gocv"
)

//---------------------------------------------------------------------------------
type Program struct {
	Name              string
	yolov3WeightsPath string
	yolov3ConfigPath  string
	cocoNamesPath     string
	yolonet           yolov3.Net
	objects           []yolov3.ObjectDetection
	framech           chan gocv.Mat
	fcuda             bool
	mtype             string
	fdetect           bool
	ok                bool
}

//---------------------------------------------------------------------------------
func main() {
	var err error
	pg := Program{
		Name:    "YOLOv3 Detection",
		framech: make(chan gocv.Mat),
		fcuda:   false,
		ok:      true,
	}

	fmt.Println(pg.Name)
	flag.StringVar(&pg.mtype, "model", pg.mtype, "use model type [normal|tiny|spp]")
	flag.BoolVar(&pg.fdetect, "detect", pg.fdetect, "start YOLO detection")
	flag.BoolVar(&pg.fcuda, "cuda", pg.fcuda, "use CUDA backend")
	flag.Parse()

	switch pg.mtype {
	case "tiny":
		pg.yolov3WeightsPath = "../assets/models/yolov3-tiny.weights"
		pg.yolov3ConfigPath = "../assets/models/yolov3-tiny.cfg"
	case "spp": // with spatial pyramid pooling
		pg.yolov3WeightsPath = "../assets/models/yolov3-spp.weights"
		pg.yolov3ConfigPath = "../assets/models/yolov3-spp.cfg"
	default:
		pg.yolov3WeightsPath = "../assets/models/yolov3.weights"
		pg.yolov3ConfigPath = "../assets/models/yolov3.cfg"
	}
	pg.cocoNamesPath = "../assets/models/coco.names"

	netconf := yolov3.DefaultConfig()
	if pg.fcuda {
		netconf.NetBackendType = gocv.NetBackendCUDA
		netconf.NetTargetType = gocv.NetTargetCUDA
		pg.Name += " using CUDA"
	}

	pg.yolonet, err = yolov3.NewNetWithConfig(pg.yolov3WeightsPath, pg.yolov3ConfigPath, pg.cocoNamesPath, netconf)
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
				// log.Println("detect and update boxes")
			default:
				// log.Println("not ready for detection. so skip")
			}
			yolov3.DrawDetections(&frame, pg.objects)
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
		if !ok || frame.Empty() {
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
