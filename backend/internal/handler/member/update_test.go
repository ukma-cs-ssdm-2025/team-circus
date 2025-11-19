package member_test

import (
	"bytes"
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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/member/requests"
	"go.uber.org/zap"
)

type mockUpdateMemberService struct {
	mock.Mock
}

func (m *mockUpdateMemberService) UpdateMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID, role string) (*domain.Member, error) {
	args := m.Called(ctx, userUUID, groupUUID, memberUUID, role)
	member, _ := args.Get(0).(*domain.Member)
	return member, args.Error(1)
}

func TestNewUpdateMemberHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockUpdateMemberService, gin.HandlerFunc) {
		mockService := &mockUpdateMemberService{}
		handler := member.NewUpdateMemberHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulUpdate", func(t *testing.T) {
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		memberUUID := uuid.New()
		authUserUUID := uuid.New()

		memberUpdated := &domain.Member{
			GroupUUID: groupUUID,
			UserUUID:  memberUUID,
			Role:      domain.RoleViewer,
			CreatedAt: time.Now().UTC(),
		}

		mockService.On("UpdateMemberByUser", mock.Anything, authUserUUID, groupUUID, memberUUID, domain.RoleViewer).
			Return(memberUpdated, nil)

		body, err := json.Marshal(requests.UpdateMemberRequest{
			Role: domain.RoleViewer,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch,
			"/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
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
		assert.Equal(t, domain.RoleViewer, response["role"])
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		_, handler := setup(t)

		req := httptest.NewRequest(http.MethodPatch, "/groups/invalid/members/"+uuid.NewString(), bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
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

		req := httptest.NewRequest(http.MethodPatch, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
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

		req := httptest.NewRequest(http.MethodPatch, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
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

		req := httptest.NewRequest(http.MethodPatch, "/groups/"+groupUUID.String()+"/members/invalid", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
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

	t.Run("InvalidJSON", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()
		memberUUID := uuid.New()

		req := httptest.NewRequest(http.MethodPatch, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), bytes.NewBufferString(`invalid`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{
			{Key: "uuid", Value: groupUUID.String()},
			{Key: "user_uuid", Value: memberUUID.String()},
		}
		c.Set("user_uid", uuid.New())

		handler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "invalid request format", response["error"])
	})

	t.Run("ValidationFailure", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()
		memberUUID := uuid.New()

		req := httptest.NewRequest(http.MethodPatch, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(),
			bytes.NewBufferString(`{"role":"invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{
			{Key: "uuid", Value: groupUUID.String()},
			{Key: "user_uuid", Value: memberUUID.String()},
		}
		c.Set("user_uid", uuid.New())

		handler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "validation failed", response["error"])
		assert.NotEmpty(t, response["details"])
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
		{"OnlyAuthor", domain.ErrOnlyAuthor, http.StatusBadRequest, "there must be an author"},
		{"Internal", errors.New("db error"), http.StatusInternalServerError, "failed to update member"},
	}

	for _, tc := range serviceErrorTests {
		tc := tc
		t.Run("ServiceError_"+tc.name, func(t *testing.T) {
			mockService, handler := setup(t)

			groupUUID := uuid.New()
			memberUUID := uuid.New()
			authUserUUID := uuid.New()

			body, err := json.Marshal(requests.UpdateMemberRequest{
				Role: domain.RoleEditor,
			})
			assert.NoError(t, err)

			mockService.On("UpdateMemberByUser", mock.Anything, authUserUUID, groupUUID, memberUUID, domain.RoleEditor).
				Return((*domain.Member)(nil), tc.serviceErr).Once()

			req := httptest.NewRequest(http.MethodPatch, "/groups/"+groupUUID.String()+"/members/"+memberUUID.String(), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
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
