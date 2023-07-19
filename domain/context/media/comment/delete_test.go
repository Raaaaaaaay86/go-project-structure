package comment_test

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestDeleteCommentUseCase_Execute(t *testing.T) {
	type DeleteCommentTestCase struct {
		CommentId       primitive.ObjectID
		CommentAuthorId uint
		ExecutorId      uint
		RoleIds         []role.RoleId
		ExpectedError   error
		TestDescription string
	}

	testCases := []DeleteCommentTestCase{
		{
			TestDescription: "Delete comment by author",
			CommentId:       primitive.NewObjectID(),
			CommentAuthorId: 1,
			ExecutorId:      1,
			RoleIds:         []role.RoleId{role.User},
			ExpectedError:   nil,
		},
		{
			TestDescription: "Comment cannot deleted by non-author USER role",
			CommentId:       primitive.NewObjectID(),
			CommentAuthorId: 2,
			ExecutorId:      1,
			RoleIds:         []role.RoleId{role.User},
			ExpectedError:   exception.ErrUnauthorized,
		},
		{
			TestDescription: "Delete a non-exists comment",
			CommentId:       [12]byte{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			CommentAuthorId: 0,
			ExecutorId:      1,
			RoleIds:         []role.RoleId{role.Admin},
			ExpectedError:   mongo.ErrNoDocuments,
		},
	}

	for i, tc := range testCases {
		t.Logf("Test case [%d] - %s", i, tc.TestDescription)

		videoCommentRepository := mocks.NewVideoCommentRepository(t)
		switch tc.ExpectedError {
		case nil:
			videoCommentRepository.
				On("DeleteById", mock.Anything, tc.CommentId, tc.ExecutorId).
				Return(0, nil).
				Once()
		case mongo.ErrNoDocuments:
			videoCommentRepository.
				On("DeleteById", mock.Anything, tc.CommentId, tc.ExecutorId).
				Return(0, mongo.ErrNoDocuments).
				Once()
		case exception.ErrUnauthorized:
			videoCommentRepository.
				On("DeleteById", mock.Anything, tc.CommentId, tc.ExecutorId).
				Return(0, exception.ErrUnauthorized).
				Once()
		}

		useCase := comment.NewDeleteCommentUseCase(tracing.NewEmptyTracerProvider(), videoCommentRepository)
		cmd := comment.DeleteCommentCommand{
			CommentId:  tc.CommentId,
			ExecutorId: tc.ExecutorId,
			RoleIds:    tc.RoleIds,
		}
		_, err := useCase.Execute(context.TODO(), cmd)
		if err != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}

		assert.NoError(t, err)
	}
}
