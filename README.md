# Fire Detection System

This project implements a fire detection system using the Go programming language and the OpenCV (gocv) library. The system processes video frames to detect fire based on color and brightness masks and checks for flickering.

## Project Structure

- `cmd/` - Contains the main executable
    - `main.go` - The entry point of the application
- `pkg/` - Contains the library code
    - `detect/` - Contains the detection logic
        - `detector.go` - Base detection logic
        - `fire_detector.go` - Fire detection logic
    - `utils/` - Contains utility code
        - `logger.go` - Logging setup
        - `singleton.go` - Singleton pattern implementation

## Design

```mermaid
sequenceDiagram
    participant User
    participant Main as main.go
    participant Utils as utils/singleton.go
    participant Detect as detect/fire_detector.go
    participant OpenCV as gocv

    User->>Main: Start Application
    Main->>Utils: Setup Logger
    Main->>Utils: Initialize Singleton
    Utils-->>Main: Webcam and Window Instances
    Main->>OpenCV: Open Webcam
    OpenCV-->>Main: Webcam Instance
    Main->>OpenCV: Create Window
    OpenCV-->>Main: Window Instance
    Main->>OpenCV: Read Frame
    OpenCV-->>Main: Frame
    Main->>Detect: Detect Fire
    Detect->>OpenCV: Convert to HSV
    Detect->>OpenCV: Create Color Mask
    Detect->>OpenCV: Convert to Gray
    Detect->>OpenCV: Create Brightness Mask
    Detect->>OpenCV: Combine Masks
    Detect->>OpenCV: Check Flickering
    Detect-->>Main: Fire Detected?
    Main->>OpenCV: Draw Fire Box
    Main->>OpenCV: Display Frame
    Main->>User: Show Result
```

## Getting Started

### Prerequisites

- Go programming language installed
- OpenCV library installed with Go bindings (gocv)

### Installation

1. Clone the repository:
   ```sh
   git clone hhttps://github.com/arturogonzalezm/fire_detector.git
   cd fire_detector
    ```
