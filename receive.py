import socket
import pygame
import base64

host = "localhost"
port = 1890

screen = pygame.display.set_mode( (640, 480) )

while True:
  s = socket.socket( socket.AF_INET, socket.SOCK_STREAM )
  s.setsockopt( socket.SOL_SOCKET, socket.SO_REUSEADDR, 1 )
  s.bind( (host, port) )
  s.listen(1)
  conn, addr = s.accept()
  message = []
  while True:
    d = conn.recv( 1228800 )
    if not d: break
    else: message.append( d )
  conn.close()
  data = ''.join( message )
  data = base64.b64decode( data )
  image = pygame.image.fromstring( data, (640,480), "RGBA" )
  screen.blit( image, (0,0) )
  pygame.display.update()
