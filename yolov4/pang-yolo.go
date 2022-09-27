//==========================================================================================
// Name: Moth WebSocket Pang Video Caster using YOLO detection
// Author: Stony Kang, sikang99@gmail.com
// Reference: https://godoc.org/gocv.io/x/gocv
// Copyright: TeamGRIT, 2021
//==========================================================================================
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
	// "gocv.io/x/gocv/cuda"
)

//---------------------------------------------------------------------------------
type Program struct {
	name       string
	url        string
	device     int
	media      string
	label      string
	scale      float64
	version    bool
	method     string
	dnn        DNN
	boxColor   color.RGBA
	labelColor color.RGBA
}

// Deep Neural Network (DNN) parameters
type DNN struct {
	model   string
	config  string
	coco    string
	classes []string
	names   []string
	net     gocv.Net
	backend gocv.NetBackendType
	target  gocv.NetTargetType
}

var (
	redColor   = color.RGBA{255, 0, 0, 0}
	greenColor = color.RGBA{0, 255, 0, 0}
	blueColor  = color.RGBA{0, 0, 255, 0}
)

var pg = &Program{
	name:   "MothCam Pang Cast, (c)TeamGRIT: YOLOv4",
	url:    "wss://localhost:8277/pang/ws/pub?channel=c40hp6epjh65aeq6ne50",
	label:  "objects detected",
	method: "yolov4",
	dnn: DNN{
		model:  "../assets/models/yolov4-tiny.weights",
		config: "../assets/models/yolov4-tiny.cfg",
		coco:   "../assets/models/coco.names",
		// backend: gocv.NetBackendDefault,
		// target:  gocv.NetTargetCPU,
		backend: gocv.NetBackendVKCOM,
		target:  gocv.NetTargetVulkan,
	},
	device:     0,
	scale:      1.0,
	boxColor:   greenColor,
	labelColor: redColor,
}

//---------------------------------------------------------------------------------
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//---------------------------------------------------------------------------------
func main() {
	fmt.Println(pg.name)

	flag.StringVar(&pg.url, "url", pg.url, "target url to send video")
	flag.StringVar(&pg.media, "media", pg.media, "media file or url(rtsp) to use")
	flag.StringVar(&pg.label, "label", pg.label, "label string on screen")
	flag.IntVar(&pg.device, "device", pg.device, "local device id to use")
	flag.Float64Var(&pg.scale, "scale", pg.scale, "scale of video screen")
	flag.BoolVar(&pg.version, "version", pg.version, "version number of engine")
	flag.Parse()

	if pg.version {
		fmt.Println("GoCV:", gocv.Version())
		fmt.Println("OpenCV:", gocv.OpenCVVersion())
		// devices := cuda.GetCudaEnabledDeviceCount()
		// for i := 0; i < devices; i++ {
		// 	fmt.Print("  ")
		// 	cuda.PrintShortCudaDeviceInfo(i)
		// }
		return
	}

	fmt.Println("[Config] method:", pg.method, "media:", pg.media, "device:", pg.device)

	video, err := pg.openVideoCapture(pg.media, pg.device)
	if err != nil {
		log.Println(err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(pg.name)
	defer window.Close()

	ws, _, err := websocket.DefaultDialer.Dial(pg.url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	mime := fmt.Sprintf("video/jpeg; width=%.0f; height=%.0f",
		video.Get(gocv.VideoCaptureFrameWidth), video.Get(gocv.VideoCaptureFrameHeight))
	log.Println(mime)
	err = ws.WriteMessage(websocket.TextMessage, []byte(mime))
	if err != nil {
		log.Println(err)
		return
	}

	img := gocv.NewMat()
	defer img.Close()

	// setup for yolo detection
	pg.beforeDection()
	defer pg.afterDection()

	for video.IsOpened() {
		ok := video.Read(&img)
		if !ok {
			log.Println("capture device closed")
			break
		}
		if img.Empty() {
			continue
		}

		if pg.scale != 1.0 {
			gocv.Resize(img, &img, image.Point{}, pg.scale, pg.scale, 0)
		}

		nb, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Println(err)
			break
		}

		data := nb.GetBytes()

		err = ws.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Println(err)
			break
		}

		img = pg.doDetection(&img)

		window.IMShow(img)
		if window.WaitKey(1) == 27 { // Esc key
			break
		}
	}
}

//------------------------------------------------------------------------------------------
// open a Video Source such as camera device, media file, ...
func (d *Program) openVideoCapture(media string, device int) (video *gocv.VideoCapture, err error) {
	if media != "" {
		video, err = gocv.VideoCaptureFile(media)
		if err != nil {
			log.Println(err, media)
			return
		}
	} else {
		video, err = gocv.VideoCaptureDevice(device)
		if err != nil {
			log.Println(err, device)
			return
		}
	}
	return
}

//------------------------------------------------------------------------------------------
func (d *Program) beforeDection() (err error) {
	d.dnn.classes = d.readCOCO()

	d.dnn.net = gocv.ReadNet(d.dnn.model, d.dnn.config)
	if d.dnn.net.Empty() {
		log.Println("read error for network model:", d.dnn.model, d.dnn.config)
		return
	}
	d.dnn.net.SetPreferableBackend(d.dnn.backend)
	d.dnn.net.SetPreferableTarget(d.dnn.target)

	d.dnn.names = d.getOutputsNames(&d.dnn.net)
	return
}

//------------------------------------------------------------------------------------------
func (d *Program) afterDection() (err error) {
	d.dnn.net.Close()
	return
}

//------------------------------------------------------------------------------------------
func (d *Program) doDetection(frame *gocv.Mat) (detectImg gocv.Mat) {
	detectImg, detectClass := d.detect(&d.dnn.net, frame.Clone(), 0.40, 0.5, d.dnn.names, d.dnn.classes)
	if len(detectClass) > 0 {
		log.Printf("Detect Class : %v\n", detectClass)
	}
	return
}

//------------------------------------------------------------------------------------------
// readCOCO : Read coco.names
func (d *Program) readCOCO() (classes []string) {
	read, _ := os.Open(d.dnn.coco)
	defer read.Close()
	for {
		var t string
		_, err := fmt.Fscan(read, &t)
		if err != nil {
			break
		}
		classes = append(classes, t)
	}
	return
}

//------------------------------------------------------------------------------------------
// getOutputsNames : YOLO Layer
func (d *Program) getOutputsNames(net *gocv.Net) (outputs []string) {
	for _, i := range net.GetUnconnectedOutLayers() {
		layer := net.GetLayer(i)
		layerName := layer.GetName()
		if layerName != "_input" {
			outputs = append(outputs, layerName)
		}
	}
	return
}

//------------------------------------------------------------------------------------------
// detect : Run YOLOv4 Process
func (d *Program) detect(net *gocv.Net, src gocv.Mat, scoreThreshold float32, nmsThreshold float32, OutputNames []string, classes []string) (gocv.Mat, []string) {
	img := src.Clone()
	img.ConvertTo(&img, gocv.MatTypeCV32F)
	blob := gocv.BlobFromImage(img, 1/255.0, image.Pt(416, 416), gocv.NewScalar(0, 0, 0, 0), true, false)
	net.SetInput(blob, "")
	probs := net.ForwardLayers(OutputNames)
	boxes, confidences, classIds := d.postProcess(img, &probs)

	indices := make([]int, 100)
	if len(boxes) == 0 { // No Classes
		return src, []string{}
	}
	gocv.NMSBoxes(boxes, confidences, scoreThreshold, nmsThreshold, indices)

	return d.drawRect(src, boxes, classes, classIds, indices)
}

//------------------------------------------------------------------------------------------
// PostProcess : All Detect Box
func (d *Program) postProcess(frame gocv.Mat, outs *[]gocv.Mat) ([]image.Rectangle, []float32, []int) {
	var classIds []int
	var confidences []float32
	var boxes []image.Rectangle

	for _, out := range *outs {
		data, _ := out.DataPtrFloat32()
		for i := 0; i < out.Rows(); i, data = i+1, data[out.Cols():] {
			scoresCol := out.RowRange(i, i+1)

			scores := scoresCol.ColRange(5, out.Cols())
			_, confidence, _, classIDPoint := gocv.MinMaxLoc(scores)

			if confidence > 0.5 {
				centerX := int(data[0] * float32(frame.Cols()))
				centerY := int(data[1] * float32(frame.Rows()))
				width := int(data[2] * float32(frame.Cols()))
				height := int(data[3] * float32(frame.Rows()))

				left := centerX - width/2
				top := centerY - height/2
				classIds = append(classIds, classIDPoint.X)
				confidences = append(confidences, float32(confidence))
				boxes = append(boxes, image.Rect(left, top, width, height))
			}
		}
	}

	return boxes, confidences, classIds
}

//------------------------------------------------------------------------------------------
// drawRect : Detect Class to Draw Rect
func (d *Program) drawRect(img gocv.Mat, boxes []image.Rectangle, classes []string, classIds []int, indices []int) (gocv.Mat, []string) {
	var detectClass []string
	for _, idx := range indices {
		if idx == 0 {
			continue
		}
		gocv.Rectangle(&img, image.Rect(boxes[idx].Max.X, boxes[idx].Max.Y, boxes[idx].Max.X+boxes[idx].Min.X, boxes[idx].Max.Y+boxes[idx].Min.Y), color.RGBA{255, 0, 0, 0}, 2)
		gocv.PutText(&img, classes[classIds[idx]], image.Point{boxes[idx].Max.X, boxes[idx].Max.Y + 30}, gocv.FontHersheyPlain, 1.0, color.RGBA{0, 0, 255, 0}, 3)
		detectClass = append(detectClass, classes[classIds[idx]])
	}
	return img, detectClass
}

//==========================================================================================
