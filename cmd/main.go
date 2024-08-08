package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"fire_detector/pkg/detect"
	"fire_detector/pkg/utils"
	"gocv.io/x/gocv"
)

func main() {
	// Setup logging
	utils.SetupLogger()

	// Initialize the singleton instance
	err := utils.InitializeSingleton()
	if err != nil {
		log.Printf("Initialization failed: %v\n", err)
		fmt.Printf("Initialization failed: %v\n", err)
		return
	}

	webcam := utils.GetWebcamInstance()
	if webcam == nil {
		log.Println("Failed to open webcam")
		fmt.Println("Failed to open webcam")
		return
	}
	defer webcam.Close()

	window := utils.GetWindowInstance()
	if window == nil {
		log.Println("Failed to create window")
		fmt.Println("Failed to create window")
		return
	}
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	prevGray := gocv.NewMat()
	defer prevGray.Close()

	detector := detect.NewFireDetector(prevGray)

	for {
		if ok := webcam.Read(&img); !ok {
			log.Println("Device closed")
			fmt.Println("Device closed")
			return
		}
		if img.Empty() {
			continue
		}

		fireDetected, fireMask, debugInfo := detector.Detect(img)

		if fireDetected {
			detector.DrawFireBox(img, fireMask)
			gocv.PutText(&img, "FIRE DETECTED", image.Point{X: 10, Y: 30}, gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 0, 0}, 2) // Red text
			log.Println("FIRE DETECTED")
			fmt.Println("FIRE DETECTED")
		} else {
			gocv.PutText(&img, "NO FIRE DETECTED", image.Point{X: 10, Y: 30}, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 0}, 2) // Green text
			log.Println("NO FIRE DETECTED")
			fmt.Println("NO FIRE DETECTED")
		}

		// Display debug info
		gocv.PutText(&img, debugInfo, image.Point{X: 10, Y: 60}, gocv.FontHersheyPlain, 1.0, color.RGBA{255, 255, 255, 0}, 1)
		log.Println(debugInfo)
		fmt.Println(debugInfo)

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
