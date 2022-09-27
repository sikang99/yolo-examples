package main

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	runtime.LockOSThread()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	args := strings.Fields("gst-launch-1.0 videotestsrc ! video/x-raw,width=1280,height=720 ! autovideosink")
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(out))
}
