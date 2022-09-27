package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func main() {

	// define default hog descriptor
	hog := gocv.NewHOGDescriptor()
	defer hog.Close()
	hog.SetSVMDetector(gocv.HOGDefaultPeopleDetector())

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// read image
	img := gocv.IMRead("images/person_010.bmp", 0)
	if img.Empty() {
		log.Fatalln("image not ready")
	}

	//resize image
	fact := float64(400) / float64(img.Cols())
	newY := float64(img.Rows()) * fact
	gocv.Resize(img, &img, image.Point{X: 400, Y: int(newY)}, 0, 0, 1)

	// detect people in image
	rects := hog.DetectMultiScaleWithParams(img, 0, image.Point{X: 8, Y: 8}, image.Point{X: 16, Y: 16}, 1.05, 2, false)

	// print found points
	printStruct(rects)

	// draw a rectangle around each face on the original image,
	// along with text identifing as "Human"
	for _, r := range rects {
		gocv.Rectangle(&img, r, blue, 3)

		size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
	}

	if ok := gocv.IMWrite("loool.jpg", img); !ok {
		fmt.Println("Error")
	}

}

func printStruct(i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
