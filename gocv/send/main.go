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
	src := `videotestsrc ! video/x-raw,width=1280,height=720 ! appsink`
	webcam, err := gocv.OpenVideoCapture(src)
	if err != nil {
		log.Fatalln(err)
	}

	window := gocv.NewWindow("Send")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
