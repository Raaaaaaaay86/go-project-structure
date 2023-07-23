package repository

import (
	"github.com/google/wire"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
)

const pkg = "adapter.port_out.repository"

var RepositorySet = wire.NewSet(
	NewUserRepository,
	wire.Bind(new(repository.UserRepository), new(*UserRepository)),
	NewVideoCommentRepository,
	wire.Bind(new(repository.VideoCommentRepository), new(*VideoCommentRepository)),
	NewVideoPostRepository,
	wire.Bind(new(repository.VideoPostRepository), new(*VideoPostRepository)),
)
