package user_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/testutil"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user"
	"go.uber.org/zap"
)

type mockGetUserService struct {
	mock.Mock
}

func (m *mockGetUserService) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck
}

type MockGetAllUsersService struct {
	mock.Mock
}

const usersEndpoint = "/users"

func (m *MockGetAllUsersService) GetAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1) //nolint:errcheck
}

func TestNewGetUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockGetUserService, gin.HandlerFunc) {
		mockService := &mockGetUserService{}
		handler := user.NewGetUserHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulGetUser", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		expectedUser := &domain.User{
			UUID:      userUUID,
			Login:     "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("GetByUUID", mock.Anything, userUUID).Return(expectedUser, nil)

		c, w := testutil.NewRequestContext(t, http.MethodGet, fmt.Sprintf("%s/%s", usersEndpoint, userUUID))
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "testuser", response["login"])
		assert.Equal(t, "test@example.com", response["email"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		c, w := testutil.NewRequestContext(t, http.MethodGet, fmt.Sprintf("%s/%s", usersEndpoint, "invalid-uuid"))
		c.Params = gin.Params{{Key: "uuid", Value: "invalid-uuid"}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "invalid uuid format", response["error"])

		mockService.AssertNotCalled(t, "GetByUUID")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetByUUID", mock.Anything, userUUID).Return(nil, domain.ErrUserNotFound)

		c, w := testutil.NewRequestContext(t, http.MethodGet, fmt.Sprintf("%s/%s", usersEndpoint, userUUID))
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "user not found", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetByUUID", mock.Anything, userUUID).Return(nil, domain.ErrInternal)

		c, w := testutil.NewRequestContext(t, http.MethodGet, fmt.Sprintf("%s/%s", usersEndpoint, userUUID))
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "failed to get user", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetByUUID", mock.Anything, userUUID).Return(nil, errors.New("database connection failed"))

		c, w := testutil.NewRequestContext(t, http.MethodGet, fmt.Sprintf("%s/%s", usersEndpoint, userUUID))
		c.Params = gin.Params{{Key: "uuid", Value: userUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "failed to get user", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestNewGetAllUsersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*MockGetAllUsersService, gin.HandlerFunc) {
		mockService := &MockGetAllUsersService{}
		handler := user.NewGetAllUsersHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulGetAllUsers", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		user1 := &domain.User{
			UUID:      uuid.New(),
			Login:     "user1",
			Email:     "user1@example.com",
			Password:  "hashedpassword1",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		user2 := &domain.User{
			UUID:      uuid.New(),
			Login:     "user2",
			Email:     "user2@example.com",
			Password:  "hashedpassword2",
			CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		}

		expectedUsers := []*domain.User{user1, user2}
		mockService.On("GetAll", mock.Anything).Return(expectedUsers, nil)

		c, w := testutil.NewRequestContext(t, http.MethodGet, usersEndpoint)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Contains(t, response, "users")

		users := response["users"].([]interface{}) //nolint:errcheck
		assert.Len(t, users, 2)
		assert.Equal(t, "user1", users[0].(map[string]interface{})["login"]) //nolint:errcheck
		assert.Equal(t, "user2", users[1].(map[string]interface{})["login"]) //nolint:errcheck

		mockService.AssertExpectations(t)
	})

	t.Run("EmptyUsersList", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		expectedUsers := []*domain.User{}
		mockService.On("GetAll", mock.Anything).Return(expectedUsers, nil)

		c, w := testutil.NewRequestContext(t, http.MethodGet, usersEndpoint)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Contains(t, response, "users")

		users := response["users"].([]interface{}) //nolint:errcheck
		assert.Len(t, users, 0)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		mockService.On("GetAll", mock.Anything).Return(nil, domain.ErrInternal)

		c, w := testutil.NewRequestContext(t, http.MethodGet, usersEndpoint)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "failed to get users", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		mockService.On("GetAll", mock.Anything).Return(nil, errors.New("database connection failed"))

		c, w := testutil.NewRequestContext(t, http.MethodGet, usersEndpoint)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "failed to get users", response["error"])

		mockService.AssertExpectations(t)
	})
}
