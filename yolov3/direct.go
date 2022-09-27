package main

import (
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func main() {
	// open capture device
	video, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Println(err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow("YOLOv3")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	// open DNN object tracking model
	net := gocv.ReadNet("../assets/models/yolov3.weights", "../assets/models/yolov3.cfg")
	if net.Empty() {
		log.Println("Error reading network model")
		return
	}
	defer net.Close()

	// for Macbook
	net.SetPreferableBackend(gocv.NetBackendType(gocv.NetBackendDefault))
	net.SetPreferableTarget(gocv.NetTargetType(gocv.NetTargetCPU))

	var ratio float64 = 0.00392
	var mean gocv.Scalar = gocv.NewScalar(0, 0, 0, 0)
	var swapRGB bool = true

	log.Println("Start reading device")
	firsttime := true

	for {
		if !video.Read(&img) {
			log.Printf("Device closed")
			return
		}
		if img.Empty() {
			continue
		}

		// convert image Mat to 416x416 blob
		blob := gocv.BlobFromImage(img, ratio, image.Pt(416, 416), mean, swapRGB, false)

		// feed the blob into the detector
		net.SetInput(blob, "")

		// run a forward pass thru the network
		prob := net.Forward("")
		if firsttime {
			log.Printf("prob.Total() = %v, prob.Size() = %v\n", prob.Total(), prob.Size())
			firsttime = false
		}

		performDetection(&img, prob)

		prob.Close()
		blob.Close()

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func performDetection(frame *gocv.Mat, results gocv.Mat) {
	totalResults := results.Total() / 85
	for i := 0; i < totalResults; i++ {
		confidence := results.GetFloatAt(i, 4)
		if confidence > 0.03 {
			center_x := int(results.GetFloatAt(i, 0) * float32(frame.Cols()))
			center_y := int(results.GetFloatAt(i, 1) * float32(frame.Rows()))
			width := int(results.GetFloatAt(i, 2) * float32(frame.Cols()))
			height := int(results.GetFloatAt(i, 3) * float32(frame.Rows()))
			gocv.Rectangle(frame, image.Rect(center_x-width/2, center_y-height/2, center_x+width/2, center_y+height/2), color.RGBA{0, 255, 0, 0}, 2)
		}
	}
}
