package detect

import (
	"fmt"
	"image/color"

	"gocv.io/x/gocv"
)

// Adjustable parameters
const (
	HueMin              = 0
	HueMax              = 30
	SaturationMin       = 100
	ValueMin            = 100
	BrightnessThreshold = 150
	FireRatioThreshold  = 0.0005 // 0.05% of the image
	MinContourArea      = 50
	ConsistentFrames    = 2
	FlickerThreshold    = 20
)

type FireDetector struct {
	Detector
	FireCounter int
}

func NewFireDetector(prevGray gocv.Mat) *FireDetector {
	return &FireDetector{
		Detector: *NewDetector(prevGray),
	}
}

func (fd *FireDetector) Detect(img gocv.Mat) (bool, gocv.Mat, string) {
	hsv := gocv.NewMat()
	defer hsv.Close()
	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)

	// Color mask
	lowerFire := gocv.NewScalar(float64(HueMin), float64(SaturationMin), float64(ValueMin), 0)
	upperFire := gocv.NewScalar(float64(HueMax), 255, 255, 0)
	mask := gocv.NewMat()
	defer mask.Close()
	gocv.InRangeWithScalar(hsv, lowerFire, upperFire, &mask)

	// Brightness mask
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	brightMask := gocv.NewMat()
	defer brightMask.Close()
	gocv.Threshold(gray, &brightMask, float32(BrightnessThreshold), 255, gocv.ThresholdBinary)

	// Combine color and brightness
	fireMask := gocv.NewMat()
	gocv.BitwiseAnd(mask, brightMask, &fireMask)

	// Check for flickering
	flickerRatio := 0.0
	if !fd.PrevGray.Empty() {
		diff := gocv.NewMat()
		defer diff.Close()
		gocv.AbsDiff(gray, fd.PrevGray, &diff)
		flickerMask := gocv.NewMat()
		defer flickerMask.Close()
		gocv.Threshold(diff, &flickerMask, float32(FlickerThreshold), 255, gocv.ThresholdBinary)
		gocv.BitwiseAnd(fireMask, flickerMask, &fireMask)

		flickerPixels := gocv.CountNonZero(flickerMask)
		flickerRatio = float64(flickerPixels) / float64(img.Rows()*img.Cols())
	}
	gray.CopyTo(&fd.PrevGray)

	// Calculate fire ratio
	firePixels := gocv.CountNonZero(fireMask)
	totalPixels := img.Rows() * img.Cols()
	fireRatio := float64(firePixels) / float64(totalPixels)

	// Check for consistent detection
	if fireRatio > FireRatioThreshold {
		fd.FireCounter++
	} else {
		fd.FireCounter = 0
	}

	fireDetected := fd.FireCounter >= ConsistentFrames

	// Debug info
	debugInfo := fmt.Sprintf("Fire Ratio: %.4f, Flicker Ratio: %.4f, Counter: %d", fireRatio, flickerRatio, fd.FireCounter)

	return fireDetected, fireMask, debugInfo
}

func (fd *FireDetector) DrawFireBox(img gocv.Mat, mask gocv.Mat) {
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		area := gocv.ContourArea(contour)
		if area > MinContourArea {
			rect := gocv.BoundingRect(contour)
			gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 2) // Red box
		}
	}
}
