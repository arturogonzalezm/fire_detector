package detect

import "gocv.io/x/gocv"

type Detector struct {
	PrevGray gocv.Mat
}

func NewDetector(prevGray gocv.Mat) *Detector {
	return &Detector{PrevGray: prevGray}
}

func (d *Detector) Detect(img gocv.Mat) (bool, gocv.Mat, string) {
	// This method should be overridden by child structs
	return false, gocv.NewMat(), ""
}

func (d *Detector) DrawFireBox(img gocv.Mat, mask gocv.Mat) {
	// This method should be overridden by child structs
}
