package main

import (
	"log"

	"gocv.io/x/gocv"
)

func main() {
	video, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Println(err)
		return
	}
	defer video.Close()

	src_window := gocv.NewWindow("Video : FromBytes")
	defer src_window.Close()
	dst_window := gocv.NewWindow("Video : Source")
	defer dst_window.Close()

	src_mat := gocv.NewMat()

	for video.Read(&src_mat) {
		buf := src_mat.ToBytes()

		dst_mat, err := gocv.NewMatFromBytes(src_mat.Rows(), src_mat.Cols(), src_mat.Type(), buf)
		if err != nil {
			log.Println(err)
			return
		}

		src_window.IMShow(src_mat)
		dst_window.IMShow(dst_mat)
		if src_window.WaitKey(1) >= 0 || dst_window.WaitKey(1) >= 0 {
			break
		}
	}
}
