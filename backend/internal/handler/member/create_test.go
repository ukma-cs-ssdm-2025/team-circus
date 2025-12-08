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

type mockCreateMemberService struct {
	mock.Mock
}

func (m *mockCreateMemberService) CreateMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID,
	role string) (*domain.Member, error) {
	args := m.Called(ctx, userUUID, groupUUID, memberUUID, role)
	member := args.Get(0).(*domain.Member) //nolint:errcheck
	return member, args.Error(1)
}

func TestNewCreateMemberHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockCreateMemberService, gin.HandlerFunc) {
		mockService := &mockCreateMemberService{}
		handler := member.NewCreateMemberHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulCreate", func(t *testing.T) {
		mockService, handler := setup(t)
		groupUUID := uuid.New()
		requestedUserUUID := uuid.New()
		authUserUUID := uuid.New()

		memberCreatedAt := time.Now().UTC()
		memberDomain := &domain.Member{
			GroupUUID: groupUUID,
			UserUUID:  requestedUserUUID,
			Role:      domain.RoleEditor,
			CreatedAt: memberCreatedAt,
		}

		mockService.On("CreateMemberByUser", mock.Anything, authUserUUID, groupUUID, requestedUserUUID, domain.RoleEditor).
			Return(memberDomain, nil)

		body, err := json.Marshal(requests.CreateMemberRequest{
			UserUUID: requestedUserUUID,
			Role:     domain.RoleEditor,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/groups/"+groupUUID.String()+"/members", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
		c.Set("user_uid", authUserUUID)

		handler(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, groupUUID.String(), response["group_uuid"])
		assert.Equal(t, requestedUserUUID.String(), response["user_uuid"])
		assert.Equal(t, domain.RoleEditor, response["role"])
		assert.NotEmpty(t, response["created_at"])
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		_, handler := setup(t)

		req := httptest.NewRequest(http.MethodPost, "/groups/invalid/members", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: "invalid"}}

		handler(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, "invalid group uuid format", response["error"])
	})

	t.Run("MissingUserContext", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()

		req := httptest.NewRequest(http.MethodPost, "/groups/"+groupUUID.String()+"/members", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
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

		req := httptest.NewRequest(http.MethodPost, "/groups/"+groupUUID.String()+"/members", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
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

	t.Run("InvalidJSON", func(t *testing.T) {
		_, handler := setup(t)
		groupUUID := uuid.New()

		req := httptest.NewRequest(http.MethodPost, "/groups/"+groupUUID.String()+"/members", bytes.NewBufferString(`invalid`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
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

		req := httptest.NewRequest(http.MethodPost, "/groups/"+groupUUID.String()+"/members",
			bytes.NewBufferString(`{"role":"invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: groupUUID.String()}}
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
		{"UserNotFound", domain.ErrUserNotFound, http.StatusNotFound, "user not found"},
		{"AlreadyExists", domain.ErrAlreadyExists, http.StatusConflict, "member already exists"},
		{"OnlyAuthor", domain.ErrOnlyAuthor, http.StatusBadRequest, "there must be only one author"},
		{"Internal", errors.New("db failed"), http.StatusInternalServerError, "failed to create member"},
	}

	for _, tc := range serviceErrorTests {
		tc := tc
		t.Run("ServiceError_"+tc.name, func(t *testing.T) {
			mockService, handler := setup(t)
			groupUUID := uuid.New()
			requestedUserUUID := uuid.New()
			authUserUUID := uuid.New()

			body, err := json.Marshal(requests.CreateMemberRequest{
				UserUUID: requestedUserUUID,
				Role:     domain.RoleViewer,
			})
			assert.NoError(t, err)

			mockService.On("CreateMemberByUser", mock.Anything, authUserUUID, groupUUID, requestedUserUUID, domain.RoleViewer).
				Return((*domain.Member)(nil), tc.serviceErr).Once()

			req := httptest.NewRequest(http.MethodPost, "/groups/"+groupUUID.String()+"/members", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
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
