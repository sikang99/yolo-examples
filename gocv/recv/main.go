package main

import (
	"log"
	"runtime"

	"gocv.io/x/gocv"
)

func init() {
	runtime.LockOSThread()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	src := `udpsrc port=9000 ! application/x-rtp,encoding-name=JPEG,payload=26 ! rtpjpegdepay ! jpegdec ! videoconvert ! appsink`
	video, err := gocv.OpenVideoCapture(src)
	if err != nil {
		log.Fatalln(err)
	}

	window := gocv.NewWindow("Recv")
	window.Close()

	img := gocv.NewMat()
	img.Close()

	for {
		video.Read(&img)
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
