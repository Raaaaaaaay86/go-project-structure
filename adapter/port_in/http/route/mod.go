package route

import "github.com/google/wire"

const pkg = "adapter.port_in.http.route"

var RouteSet = wire.NewSet(
	NewAuthenticationController,
	wire.Bind(new(IAuthenticationController), new(*AuthenticationController)),
	NewVideoController,
	wire.Bind(new(IVideoController), new(*VideoController)),
	NewCommentController,
	wire.Bind(new(ICommentController), new(*CommentController)),
)
