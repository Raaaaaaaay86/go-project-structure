package comment_test

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestDeleteCommentUseCase_Execute(t *testing.T) {
	type DeleteCommentTestCase struct {
		CommentId       primitive.ObjectID
		CommentAuthorId uint
		ExecutorId      uint
		RoleId          role.RoleId
		ExpectedError   error
		Description     string
	}

	testCases := []DeleteCommentTestCase{
		{
			Description:     "Delete comment by author",
			CommentId:       primitive.NewObjectID(),
			CommentAuthorId: 1,
			ExecutorId:      1,
			RoleId:          role.User,
			ExpectedError:   nil,
		},
		{
			Description:     "Comment cannot deleted by non-author USER role",
			CommentId:       primitive.NewObjectID(),
			CommentAuthorId: 2,
			ExecutorId:      1,
			RoleId:          role.User,
			ExpectedError:   exception.ErrUnauthorized,
		},
		{
			Description:     "Comment can deleted by ADMIN role",
			CommentId:       primitive.NewObjectID(),
			CommentAuthorId: 2,
			ExecutorId:      1,
			RoleId:          role.Admin,
			ExpectedError:   nil,
		},
		{
			Description:     "Comment can deleted by SUPER_ADMIN role",
			CommentId:       primitive.NewObjectID(),
			CommentAuthorId: 2,
			ExecutorId:      1,
			RoleId:          role.SuperAdmin,
			ExpectedError:   nil,
		},
		{
			Description:     "Delete a non-exists comment",
			CommentId:       [12]byte{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			CommentAuthorId: 0,
			ExecutorId:      1,
			RoleId:          role.Admin,
			ExpectedError:   mongo.ErrNoDocuments,
		},
	}

	for i, tc := range testCases {
		t.Logf("Test case [%d]", i)

		videoCommentRepository := mocks.NewVideoCommentRepository(t)
		switch tc.ExpectedError {
		case nil:
			videoCommentRepository.
				On("FindById", tc.CommentId).
				Return(&entity.VideoComment{Id: tc.CommentId, AuthorId: tc.CommentAuthorId}, nil).
				Once()
			videoCommentRepository.
				On("DeleteById", tc.CommentId).
				Return(nil).
				Once()
		case mongo.ErrNoDocuments:
			videoCommentRepository.
				On("FindById", tc.CommentId).
				Return(nil, mongo.ErrNoDocuments).
				Once()
		case exception.ErrUnauthorized:
			videoCommentRepository.
				On("FindById", tc.CommentId).
				Return(&entity.VideoComment{Id: tc.CommentId, AuthorId: tc.CommentAuthorId}, nil).
				Once()
		}

		useCase := comment.NewDeleteCommentUseCase(videoCommentRepository)
		cmd := comment.DeleteCommentCommand{
			CommentId:  tc.CommentId,
			ExecutorId: tc.ExecutorId,
			RoleId:     tc.RoleId,
		}
		_, err := useCase.Execute(cmd)
		if err != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}

		assert.NoError(t, err)
	}
}
