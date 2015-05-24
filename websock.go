package main

import (
  "github.com/gorilla/websocket"
)

func (h *hub) run() {
  for {
    select {
    case c := <-h.register:
      h.connections[c] = true
    case c := <-h.unregister:
      delete(h.connections, c)
      close(c.send)
    case m := <-h.broadcast:
      for c := range h.connections {
        select {
        case c.send <- m:
        default:
          delete(h.connections, c)
          close(c.send)
          go c.ws.Close()
        }
      }
    }
  }
}

func (c *connection) writer() {
  for message := range c.send {
    err := c.ws.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      break
    }
  }
  c.ws.Close()
}
