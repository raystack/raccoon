package http

import (
	"raccoon/http/grpc"
	"raccoon/http/rest"
	"raccoon/http/websocket"
)

type Handler struct {
	wh *websocket.Handler
	rh *rest.Handler
	gh *grpc.Handler
}
