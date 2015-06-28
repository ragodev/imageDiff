package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"time"
)

var debug = false
var save = true
var sensitivityPercentage = 10

func main() {

	if debug {
		startTime := time.Now()
		defer runTime(startTime)
	}

	img := loadJPEG("testr.jpg")
	imgOld := loadJPEG("testr2.jpg")
	//img := loadJPEG("2.jpg")
	//imgOld := loadJPEG("1.jpg")

	_, _, rmvGreen := rmvGreenAndCommon(img, imgOld)

	if save {
		//fmt.Println(time.Now().UnixNano())
		saveImage(rmvGreen, "output") //time.Now().String())
	}
}

func runTime(startTime time.Time) {
	endTime := time.Now()
	fmt.Println("time: ", endTime.Sub(startTime))
}

func saveImage(img image.Image, fileName string) {
	finalFile, _ := os.Create(fileName + ".jpeg")
	defer finalFile.Close()
	jpeg.Encode(finalFile, img, &jpeg.Options{jpeg.DefaultQuality})
}

//rmvGreenAndCommon comares two images and
func rmvGreenAndCommon(img image.Image, imgOld image.Image) (int, int, image.Image) {

	bounds := img.Bounds()

	rmvGreen := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rmvGreen, bounds, img, bounds.Min, draw.Src)

	greens := 0
	match := 0
	white := color.RGBA{255, 255, 255, 255}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			r, g, b, _ := img.At(x, y).RGBA()
			if debug {
				fmt.Println("x:", x, " y:", y, " r:", r, " g:", g, " b:", b, " lu:", luminance(r, g, b))
			}

			if isGreen(r, g, b) {
				greens++
				rmvGreen.Set(x, y, white)

			} else {

				rOld, gOld, bOld, _ := imgOld.At(x, y).RGBA()
				if isSimilar(r, g, b, rOld, gOld, bOld) {
					match++
					rmvGreen.Set(x, y, white)

				}
			}
		}
	}

	totalPixels := bounds.Max.Y * bounds.Max.X
	mismatchPix := totalPixels - greens - match
	fmt.Println("total pixels", totalPixels)
	fmt.Println("green pixels", greens)
	fmt.Println("matched pixels", match)
	fmt.Println("new pixels", mismatchPix)

	if mismatchPix > (totalPixels / sensitivityPercentage) {
		fmt.Println("**ALERT** Over", sensitivityPercentage, "percent change **ALERT**")
	} else {
		fmt.Println("Change under", sensitivityPercentage, "percent")
	}

	return 0, 0, rmvGreen
}

//isSimilar compares two RGB value triplets, allowing for a slight change in color or luminance.
func isSimilar(r uint32, g uint32, b uint32, rOld uint32, gOld uint32, bOld uint32) bool {
	var r1 uint32
	if rOld > 10000 {
		r1 = rOld - 10000
	}

	// var r2 uint32 = 65000
	// if rOld < 55000 {
	// 	r2 = rOld + 10000
	// }

	var b1 uint32
	if bOld > 10000 {
		b1 = bOld - 10000
	}

	var g1 uint32
	if gOld > 10000 {
		g1 = gOld - 10000
	}
	//fmt.Println(rOld, " g:", gOld, " b:", bOld)
	//fmt.Println("r:", r, " g:", g, " b:", b, " lu:", luminance(r, g, b))
	//fmt.Println("r1", r1, " r2:", r2)

	if (r < rOld+10000 && r > r1) && (b < bOld+10000 && b > b1) && (g < gOld+10000 && g > g1) {
		return true
	}
	return false
}

//isGreen returns true if the RBG input equates to a green-ish color
func isGreen(r uint32, g uint32, b uint32) bool {
	if g > r && g > b && g > 10000 {
		return true
	}
	return false
}

//luminance returns the luminance out ~65021 for RGB values r, g, b
//credit goes to https://github.com/esdrasbeleza/blzimg for luminance
func luminance(r uint32, g uint32, b uint32) uint32 {
	return uint32(0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b))
}

//loadJPEG decodes the .jpeg file "fileName" and returns an image.Image
func loadJPEG(fileName string) image.Image {
	reader, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}
