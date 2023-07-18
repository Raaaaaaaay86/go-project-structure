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
	"testing"
)

func TestForceDeleteCommentUseCase_Execute(t *testing.T) {
	type ForceDeleteTestCase struct {
		TestDescription string
		CommentId       primitive.ObjectID
		ExecutorId      uint
		RoleIds         []role.RoleId
		ExpectedError   error
	}

	testCases := []ForceDeleteTestCase{
		{
			TestDescription: "Able Force delete comment by ADMIN",
			CommentId:       primitive.NewObjectID(),
			ExecutorId:      1,
			RoleIds:         []role.RoleId{role.Admin},
			ExpectedError:   nil,
		},
		{
			TestDescription: "Able Force delete comment by SUPER_ADMIN",
			CommentId:       primitive.NewObjectID(),
			ExecutorId:      1,
			RoleIds:         []role.RoleId{role.SuperAdmin},
			ExpectedError:   nil,
		},
		{
			TestDescription: "Unable Force delete comment by USER",
			CommentId:       primitive.NewObjectID(),
			ExecutorId:      1,
			RoleIds:         []role.RoleId{role.User},
			ExpectedError:   exception.ErrUnauthorized,
		},
	}

	for i, tc := range testCases {
		t.Logf("Test case [%d]: %s", i, tc.TestDescription)
		videoCommentRepository := mocks.NewVideoCommentRepository(t)

		switch tc.ExpectedError {
		case nil:
			videoCommentRepository.On("ForceDeleteById", mock.Anything, tc.CommentId).Return(1, nil).Once()
		case exception.ErrUnauthorized:
		}

		cmd := comment.ForceDeleteCommentCommand{
			CommentId:  tc.CommentId,
			ExecutorId: tc.ExecutorId,
			RoleIds:    tc.RoleIds,
		}
		res, err := comment.NewForceDeleteCommentUseCase(tracing.NewEmptyTracerProvider(), videoCommentRepository).Execute(context.TODO(), cmd)
		if err != nil || tc.ExpectedError != nil {
			assert.ErrorIs(t, tc.ExpectedError, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, 1, res.DeleteCount)
	}
}
