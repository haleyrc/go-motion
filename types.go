package main

import (
  "github.com/lazywei/go-opencv/opencv"
  "github.com/gorilla/websocket"
  "text/template"
  "net"
)

type connection struct {
  ws *websocket.Conn
  send chan []byte
}

type hub struct {
  connections map[*connection]bool
  broadcast chan []byte
  register chan *connection
  unregister chan *connection
}

var h = hub{
  broadcast: make(chan []byte),
  register: make(chan *connection),
  unregister: make(chan *connection),
  connections: make(map[*connection]bool),
}

var l = hub{
  broadcast: make(chan []byte),
  register: make(chan *connection),
  unregister: make(chan *connection),
  connections: make(map[*connection]bool),
}

var (
  homeTempl *template.Template
  cap *opencv.Capture
  avgFrame []float64
  blurBytes []byte
  colorBytes []byte
  host = "127.0.0.1"
  port = "1890"
  remote = host + ":" + port
  conn net.Conn
)

const (
  BLUE = 0
  GREEN = 1
  RED = 2
  FRAME_WIDTH = 640
  FRAME_HEIGHT = 480
)
