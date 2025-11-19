package member_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/member"
	"go.uber.org/zap"
)

type mockGetAllMembersService struct {
	mock.Mock
}

func (m *mockGetAllMembersService) GetAllMembersForUser(ctx context.Context, userUUID, groupUUID uuid.UUID) ([]*domain.Member, error) {
	args := m.Called(ctx, userUUID, groupUUID)
	members, _ := args.Get(0).([]*domain.Member)
	return members, args.Error(1)
}

func TestNewGetAllMembersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockGetAllMembersService, gin.HandlerFunc) {
		mockService := &mockGetAllMembersService{}
		handler := member.NewGetAllMembersHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulList", func(t *testing.T) {
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		authUserUUID := uuid.New()

		member1 := &domain.Member{
			GroupUUID: groupUUID,
			UserUUID:  uuid.New(),
			Role:      domain.RoleAuthor,
			CreatedAt: time.Now().UTC(),
		}
		member2 := &domain.Member{
			GroupUUID: groupUUID,
			UserUUID:  uuid.New(),
			Role:      domain.RoleViewer,
			CreatedAt: time.Now().UTC(),
		}

		mockService.On("GetAllMembersForUser", mock.Anything, authUserUUID, groupUUID).
			Return([]*domain.Member{member1, member2}, nil)

		req := httptest.NewRequest(http.MethodGet, "/groups/"+groupUUID.String()+"/members", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", authUserUUID)

		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		members, ok := response["members"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, members, 2)
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		_, handler := setup(t)

		req := httptest.NewRequest(http.MethodGet, "/groups/invalid/members", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: "invalid"}}

		handler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "invalid uuid format", response["error"])
	})

	t.Run("MissingUserContext", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/groups/"+groupUUID.String()+"/members", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "user context missing", response["error"])
	})

	t.Run("InvalidUserContextType", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/groups/"+groupUUID.String()+"/members", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", "not-a-uuid")

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "invalid user context", response["error"])
	})

	serviceErrorTests := []struct {
		name           string
		serviceErr     error
		expectedCode   int
		expectedErrMsg string
	}{
		{"Forbidden", domain.ErrForbidden, http.StatusForbidden, "access forbidden"},
		{"GroupNotFound", domain.ErrGroupNotFound, http.StatusNotFound, "group not found"},
		{"Internal", errors.New("unexpected"), http.StatusInternalServerError, "failed to list members"},
	}

	for _, tc := range serviceErrorTests {
		tc := tc
		t.Run("ServiceError_"+tc.name, func(t *testing.T) {
			mockService, handler := setup(t)
			groupUUID := uuid.New()
			authUserUUID := uuid.New()

			mockService.On("GetAllMembersForUser", mock.Anything, authUserUUID, groupUUID).
				Return(([]*domain.Member)(nil), tc.serviceErr).Once()

			req := httptest.NewRequest(http.MethodGet, "/groups/"+groupUUID.String()+"/members", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
			c.Set("user_uid", authUserUUID)

			handler(c)

			assert.Equal(t, tc.expectedCode, w.Code)
			var response map[string]interface{}
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
			assert.Equal(t, tc.expectedErrMsg, response["error"])
		})
	}
}
