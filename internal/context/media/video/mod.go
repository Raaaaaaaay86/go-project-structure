package video

import "github.com/google/wire"

const pkg = "internal.context.media.video"

var ContextVideoSet = wire.NewSet(
	NewCreateVideoUseCase,
	wire.Bind(new(ICreateVideoUseCase), new(*CreateVideoUseCase)),
	NewLikeVideoUseCase,
	wire.Bind(new(ILikeVideoUseCase), new(*LikeVideoUseCase)),
	NewUnLikeVideoUseCase,
	wire.Bind(new(IUnLikeVideoUseCase), new(*UnLikeVideoUseCase)),
	NewUpdateVideoInfoUseCase,
	wire.Bind(new(IUpdateVideoInfoUseCase), new(*UpdateVideoInfoUseCase)),
	NewUploadVideoUseCase,
	wire.Bind(new(IUploadVideoUseCase), new(*UploadVideoUseCase)),
)
