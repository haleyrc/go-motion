package main

import (
  "github.com/gorilla/websocket"
  "net/http"
  "fmt"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
  ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
  if _, ok := err.(websocket.HandshakeError); ok {
    http.Error(w, "Not a websocket handshake", 400)
    return
  } else if err != nil {
    return
  }

  c := &connection{send: make(chan []byte, 256), ws: ws}
  fmt.Println( "Registering image client..." )
  h.register <- c
  defer func() { h.unregister <- c }()
  go captureFrames()
  c.writer()
}

func logHandler(w http.ResponseWriter, r *http.Request) {
  ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
  if _, ok := err.(websocket.HandshakeError); ok {
    http.Error(w, "Not a websocket handshake", 400)
    return
  } else if err != nil {
    return
  }

  c := &connection{send: make(chan []byte, 256), ws: ws}
  fmt.Println( "Registering log client..." )
  h.register <- c
  defer func() { h.unregister <- c }()
  c.writer()
}

func homeHandler( c http.ResponseWriter, req *http.Request) {
  homeTempl.Execute(c, req.Host)
}
