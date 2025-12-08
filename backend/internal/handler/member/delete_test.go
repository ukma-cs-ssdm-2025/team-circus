package member_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/member"
	"go.uber.org/zap"
)

type mockDeleteMemberService struct {
	mock.Mock
}

func (m *mockDeleteMemberService) DeleteMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID) error {
	args := m.Called(ctx, userUUID, groupUUID, memberUUID)
	return args.Error(0)
}

func TestNewDeleteMemberHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockDeleteMemberService, gin.HandlerFunc) {
		mockService := &mockDeleteMemberService{}
		handler := member.NewDeleteMemberHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulDelete", func(t *testing.T) {
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		memberUUID := uuid.New()
		authUserUUID := uuid.New()

		mockService.On("DeleteMemberByUser", mock.Anything, authUserUUID, groupUUID, memberUUID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{
			{Key: "uuid", Value: groupUUID.String()},
			{Key: "user_uuid", Value: memberUUID.String()},
		}
		c.Set("user_uid", authUserUUID)

		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "Member deleted successfully", response["message"])
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		_, handler := setup(t)

		req := httptest.NewRequest(http.MethodDelete, "/groups/invalid/members/"+uuid.NewString(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{
			{Key: "uuid", Value: "invalid"},
			{Key: "user_uuid", Value: uuid.NewString()},
		}

		handler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "invalid group uuid format", response["error"])
	})

	t.Run("MissingUserContext", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()
		memberUUID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{
			{Key: "uuid", Value: groupUUID.String()},
			{Key: "user_uuid", Value: memberUUID.String()},
		}

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "user context missing", response["error"])
	})

	t.Run("InvalidUserContextType", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()
		memberUUID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{
			{Key: "uuid", Value: groupUUID.String()},
			{Key: "user_uuid", Value: memberUUID.String()},
		}
		c.Set("user_uid", "not-a-uuid")

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "invalid user context", response["error"])
	})

	t.Run("InvalidMemberUUID", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/groups/"+groupUUID.String()+"/members/invalid", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{
			{Key: "uuid", Value: groupUUID.String()},
			{Key: "user_uuid", Value: "invalid"},
		}
		c.Set("user_uid", uuid.New())

		handler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "invalid user uuid format", response["error"])
	})

	serviceErrorTests := []struct {
		name           string
		serviceErr     error
		expectedCode   int
		expectedErrMsg string
	}{
		{"Forbidden", domain.ErrForbidden, http.StatusForbidden, "access forbidden"},
		{"GroupNotFound", domain.ErrGroupNotFound, http.StatusNotFound, "group not found"},
		{"UserNotFound", domain.ErrUserNotFound, http.StatusNotFound, "member not found"},
		{"OnlyAuthor", domain.ErrOnlyAuthor, http.StatusBadRequest, "cannot remove the last author"},
		{"Internal", errors.New("db down"), http.StatusInternalServerError, "failed to remove member"},
	}

	for _, tc := range serviceErrorTests {
		tc := tc
		t.Run("ServiceError_"+tc.name, func(t *testing.T) {
			mockService, handler := setup(t)

			groupUUID := uuid.New()
			memberUUID := uuid.New()
			authUserUUID := uuid.New()

			mockService.On("DeleteMemberByUser", mock.Anything, authUserUUID, groupUUID, memberUUID).
				Return(tc.serviceErr).Once()

			req := httptest.NewRequest(http.MethodDelete, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{
				{Key: "uuid", Value: groupUUID.String()},
				{Key: "user_uuid", Value: memberUUID.String()},
			}
			c.Set("user_uid", authUserUUID)

			handler(c)

			assert.Equal(t, tc.expectedCode, w.Code)
			var response map[string]interface{}
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
			assert.Equal(t, tc.expectedErrMsg, response["error"])
		})
	}
}
