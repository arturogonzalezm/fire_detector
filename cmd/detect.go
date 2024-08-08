package cmd

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"gocv.io/x/gocv"
)

// Detect starts the detection process
func Detect() error {
	// Load the pre-trained YOLOv4 model and configuration
	model := "models/yolov4.weights"
	config := "models/yolov4.cfg"
	classes := "models/coco.names"

	net := gocv.ReadNet(model, config)
	if net.Empty() {
		return fmt.Errorf("error reading network model")
	}

	// Open the video capture (webcam or video file)
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		return fmt.Errorf("error opening video capture device: %v", err)
	}
	defer func() {
		if err := webcam.Close(); err != nil {
			log.Fatalf("error closing video capture device: %v", err)
		}
	}()

	// Open a window to display the results
	window := gocv.NewWindow("Fire Detection")
	defer func() {
		if err := window.Close(); err != nil {
			log.Fatalf("error closing window: %v", err)
		}
	}()

	// Load the class names
	classNames, err := loadClassNames(classes)
	if err != nil {
		return err
	}

	// Process each frame
	for {
		frame := gocv.NewMat()
		if ok := webcam.Read(&frame); !ok {
			log.Println("Error reading frame from video capture device")
			return nil
		}
		if frame.Empty() {
			continue
		}

		// Perform detection
		blob := gocv.BlobFromImage(frame, 1.0/255.0, image.Pt(416, 416), gocv.NewScalar(0, 0, 0, 0), true, false)
		net.SetInput(blob, "")

		detections := net.ForwardLayers([]string{"yolo_139", "yolo_150", "yolo_161"}) // Update these layer names

		// Draw detections
		for _, detection := range detections {
			for i := 0; i < detection.Rows(); i++ {
				row := detection.RowRange(i, i+1)
				scores := row.ColRange(5, row.Cols())
				_, confidence, _, maxClassID := gocv.MinMaxLoc(scores)

				if confidence > 0.5 {
					centerX := int(row.GetFloatAt(0, 0) * float32(frame.Cols()))
					centerY := int(row.GetFloatAt(0, 1) * float32(frame.Rows()))
					width := int(row.GetFloatAt(0, 2) * float32(frame.Cols()))
					height := int(row.GetFloatAt(0, 3) * float32(frame.Rows()))
					left := centerX - width/2
					top := centerY - height/2

					// Ensure maxClassID is extracted correctly as an integer index
					classIdx := int(maxClassID.X) // Assuming maxClassID.X gives the correct index

					// Draw a rectangle around the detection
					gocv.Rectangle(&frame, image.Rect(left, top, left+width, top+height), color.RGBA{R: 0, G: 255, B: 0, A: 0}, 2)
					label := fmt.Sprintf("%s: %.2f", classNames[classIdx], confidence)
					gocv.PutText(&frame, label, image.Pt(left, top-10), gocv.FontHersheyPlain, 1.0, color.RGBA{R: 0, G: 255, B: 0, A: 0}, 2)
				}
			}
		}

		// Show the image in the window, and wait 1 millisecond
		window.IMShow(frame)
		window.WaitKey(1)
	}
}

// loadClassNames loads the class names from a file
func loadClassNames(file string) ([]string, error) {
	var classNames []string
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading class names file: %v", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		classNames = append(classNames, line)
	}
	return classNames, nil
}
