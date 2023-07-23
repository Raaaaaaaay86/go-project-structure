package http

import "github.com/google/wire"

//const pkg = "adapter.port_in.http"

var HttpServerSet = wire.NewSet(
	NewHttpServer,
)
