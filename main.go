package main

import (
  "github.com/lazywei/go-opencv/opencv"
  // "text/template"
	"net/http"
  "time"
  "log"
  "fmt"
  "C"
)

func captureFrames() {
  for {
    if cap.GrabFrame() {
      img := cap.RetrieveFrame(1)
      if img != nil {
        processImage(img)
      } else {
        fmt.Println("Nil image frame!")
      }
      time.Sleep(40 * time.Millisecond)
    }
  }
}

func serveLog() {
  port := 8081
  fmt.Println("Serving log on port", port)
  if err := http.ListenAndServe(":8081", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}

func main() {
  cap = opencv.NewCameraCapture(0)
  if cap == nil {
    panic("Can not open camera!")
  }
  defer cap.Release()

  avgFrame = make([]float64, FRAME_WIDTH*FRAME_HEIGHT*3)
  captureFrames()

  /*
  homeTempl = template.Must(template.ParseFiles("html/home.html"))
  go h.run()
  go l.run()
  http.HandleFunc("/", homeHandler)
  http.HandleFunc("/ws", wsHandler)
  http.HandleFunc("/log", logHandler)

  go serveLog()
  port := 8080
  fmt.Println("Serving web on port", port)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
  */
}
