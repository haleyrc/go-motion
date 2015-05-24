package main

import (
  "github.com/lazywei/go-opencv/opencv"
  "encoding/base64"
  // "image/png"
  // "utils"
  // "image"
  // "bytes"
  "math"
  // "time"
  "net"
  "fmt"
  "C"
)

func processImage(img *opencv.IplImage) {
  // defer utils.Un(utils.Trace("Processing image"))

  width := img.Width()
  height := img.Height()

  blur := opencv.CreateImage(width, height, opencv.IPL_DEPTH_8U, 3)
  defer blur.Release()

  // Blur the image
  opencv.Smooth(img, blur, opencv.CV_BLUR, 5, 5, 0, 0)
  blurBytes = C.GoBytes(blur.ImageData(), C.int(blur.ImageSize()))

  numPixels   := len(blurBytes) / 3
  numC        := img.Channels()
  colorBytes   = make([]byte, numPixels*numC)
  changeCount := 0

  for i := 0; i < numPixels; i++ {
    blueChan  := float64(blurBytes[numC*i+BLUE])
    greenChan := float64(blurBytes[numC*i+GREEN])
    redChan   := float64(blurBytes[numC*i+RED])
    grayScale := byte(0.2126 * redChan + 0.7152 * greenChan + 0.0722 * blueChan)

    // Accumulate weighted
    // C.cvRunningAvg(img, avg_frame, 0.20)
    avgFrame[numC*i+BLUE]  = (1-0.2)*avgFrame[numC*i+BLUE]  + 0.2*blueChan
    avgFrame[numC*i+GREEN] = (1-0.2)*avgFrame[numC*i+GREEN] + 0.2*greenChan
    avgFrame[numC*i+RED]   = (1-0.2)*avgFrame[numC*i+RED]   + 0.2*redChan

    // Diff the blur and scaled average
    // C.cvAbsDiff(blur, scaled_avg, diff)
    blueChan  = math.Abs(avgFrame[numC*i+BLUE]  - blueChan)
    greenChan = math.Abs(avgFrame[numC*i+GREEN] - greenChan)
    redChan   = math.Abs(avgFrame[numC*i+RED]   - redChan)

    // Grayscale the diff
    // opencv.CvtColor(diff, gray, opencv.CV_BGR2GRAY)
    grayValue :=   0.2126 * redChan + 0.7152 * greenChan + 0.0722 * blueChan

    // Threshold the blurred
    // C.cvThreshold(blur_gray, bw, 40, 255, opencv.THRESH_BINARY)
    var threshValue = 0
    if (grayValue > 40) {
      threshValue = 255
      changeCount++
    }

    // Convert back to three channel RGB
    // opencv.CvtColor(gray, color, opencv.CV_GRAY2GBR)
    if (threshValue == 0) {
      colorBytes[numC*i+BLUE]  = grayScale
      colorBytes[numC*i+GREEN] = grayScale
      colorBytes[numC*i+RED]   = grayScale
    } else {
      /*
      colorBytes[numC*i+BLUE]  = blurBytes[numC*i+BLUE]
      colorBytes[numC*i+GREEN] = blurBytes[numC*i+GREEN]
      colorBytes[numC*i+RED]   = blurBytes[numC*i+RED]
      */
      srcAlpha := 0.5
      colorBytes[numC*i+BLUE]  = byte(float64(grayScale) * (1-srcAlpha) + 1.0 * srcAlpha)
      colorBytes[numC*i+GREEN] = byte(float64(grayScale) * (1-srcAlpha) + 1.0 * srcAlpha)
      colorBytes[numC*i+RED]   = grayScale
    }
  }

  /*
  if (changeCount > 20) {
    msg := "Motion detected at " + time.Now().String()
  }
  */

  rgba := make([]uint8, numPixels * 4)
  for i := 0; i < numPixels; i++ {
    rgba[4*i+2] = uint8(colorBytes[3*i])
    rgba[4*i+1] = uint8(colorBytes[3*i+1])
    rgba[4*i]   = uint8(colorBytes[3*i+2])
    rgba[4*i+3] = uint8(255)
  }

  hexString := base64.StdEncoding.EncodeToString(rgba)

  var err error
  conn, err = net.Dial("tcp", remote)
  defer conn.Close()
  if err != nil {
    panic("Could not open tcp socket.")
  }

  if conn != nil {
    in, err := conn.Write([]byte(hexString))
    if err != nil {
      fmt.Printf("Could not write to the socket: %v, in:%v\n", err, in)
    }
  } else {
    fmt.Println("Connection is nil!")
  }


  /*
  hexString := base64.StdEncoding.EncodeToString(b.Bytes())
  h.broadcast <- []byte(hexString)
  */
}
