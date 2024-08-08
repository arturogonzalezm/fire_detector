package utils

import (
	"fmt"
	"sync"

	"gocv.io/x/gocv"
)

type singleton struct {
	webcam *gocv.VideoCapture
	window *gocv.Window
}

var (
	instance *singleton
	once     sync.Once
)

func InitializeSingleton() error {
	var err error
	once.Do(func() {
		webcam, webcamErr := gocv.OpenVideoCapture(0)
		if webcamErr != nil {
			err = fmt.Errorf("error opening video capture device: %v", webcamErr)
			return
		}

		window := gocv.NewWindow("Fire Detection")

		instance = &singleton{
			webcam: webcam,
			window: window,
		}
	})
	return err
}

func GetWebcamInstance() *gocv.VideoCapture {
	if instance == nil {
		return nil
	}
	return instance.webcam
}

func GetWindowInstance() *gocv.Window {
	if instance == nil {
		return nil
	}
	return instance.window
}
