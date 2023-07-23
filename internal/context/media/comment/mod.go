package comment

import "github.com/google/wire"

const pkg = "internal.context.media.comment"

var ContextCommentSet = wire.NewSet(
	NewCreateCommentUseCase,
	wire.Bind(new(ICreateCommentUseCase), new(*CreateCommentUseCase)),
	NewDeleteCommentUseCase,
	wire.Bind(new(IDeleteCommentUseCase), new(*DeleteCommentUseCase)),
	NewFindByVideoUseCase,
	wire.Bind(new(IFindByVideoUseCase), new(*FindByVideoUseCase)),
	NewForceDeleteCommentUseCase,
	wire.Bind(new(IForceDeleteCommentUseCase), new(*ForceDeleteCommentUseCase)),
)
