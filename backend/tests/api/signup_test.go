//go:build func_test

package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukma-cs-ssdm-2025/team-circus/tests/pkg/seeder"
	"github.com/ukma-cs-ssdm-2025/team-circus/tests/pkg/testapp"
	"github.com/ukma-cs-ssdm-2025/team-circus/tests/pkg/testdb"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUpHandler(main *testing.T) {
	setup := func(t *testing.T) (*seeder.Seeder, error) {
		db, err := testdb.NewDB()
		if err != nil {
			return nil, err
		}
		err = testdb.ResetDB(db)
		if err != nil {
			return nil, err
		}
		seeder := seeder.NewSeeder(db)

		app := testapp.NewApp()
		ctx, cancel := context.WithCancel(context.Background())
		go app.Run(ctx)
		time.Sleep(100 * time.Millisecond)

		t.Cleanup(func() {
			db.Close()
			cancel()
		})
		return seeder, nil
	}

	main.Run("SuccessfulRegistration", func(t *testing.T) {
		seeder, err := setup(t)
		require.NoError(t, err)

		// Arrange
		requestBody := map[string]string{
			"login":    "testuser",
			"email":    "test@example.com",
			"password": "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		// Act
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
			"application/json",
			bytes.NewBuffer(jsonBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "testuser", response["login"])
		assert.Equal(t, "test@example.com", response["email"])

		// Assert user in db
		user, err := seeder.GetUserByLogin("testuser")
		require.NoError(t, err)
		assert.Equal(t, "testuser", user.Login)
		assert.Equal(t, "test@example.com", user.Email)
		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte("testpassword123"))
		require.NoError(t, err)
		assert.False(t, user.CreatedAt.IsZero())
	})

	main.Run("InvalidJSON", func(t *testing.T) {
		_, err := setup(t)
		require.NoError(t, err)

		// Arrange
		invalidJSON := `{"login": "testuser", "email": "test@example.com", "password": "testpassword123"` // Missing closing brace

		// Act
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
			"application/json",
			bytes.NewBufferString(invalidJSON),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "invalid request format", response["error"])
	})

	main.Run("ValidationErrors", func(t *testing.T) {
		_, err := setup(t)
		require.NoError(t, err)

		t.Run("EmptyLogin", func(t *testing.T) {
			// Arrange
			requestBody := map[string]string{
				"login":    "",
				"email":    "test@example.com",
				"password": "testpassword123",
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Act
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			assert.Equal(t, "validation failed", response["error"])
			assert.Contains(t, response, "details")
		})

		t.Run("EmptyEmail", func(t *testing.T) {
			// Arrange
			requestBody := map[string]string{
				"login":    "testuser",
				"email":    "",
				"password": "testpassword123",
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Act
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			assert.Equal(t, "validation failed", response["error"])
			assert.Contains(t, response, "details")
		})

		t.Run("EmptyPassword", func(t *testing.T) {
			// Arrange
			requestBody := map[string]string{
				"login":    "testuser",
				"email":    "test@example.com",
				"password": "",
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Act
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			assert.Equal(t, "validation failed", response["error"])
			assert.Contains(t, response, "details")
		})
	})

	main.Run("DuplicateUser", func(t *testing.T) {
		seeder, err := setup(t)
		require.NoError(t, err)

		// Arrange
		requestBody := map[string]string{
			"login":    "duplicateuser",
			"email":    "duplicate@example.com",
			"password": "testpassword123",
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		// Act
		resp1, err := http.Post(
			fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
			"application/json",
			bytes.NewBuffer(jsonBody),
		)
		require.NoError(t, err)
		resp1.Body.Close()
		assert.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Assert
		user, err := seeder.GetUserByLogin("duplicateuser")
		require.NoError(t, err)
		assert.Equal(t, "duplicateuser", user.Login)
		assert.Equal(t, "duplicate@example.com", user.Email)
		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte("testpassword123"))
		require.NoError(t, err)
		assert.False(t, user.CreatedAt.IsZero())

		// Act
		resp2, err := http.Post(
			fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
			"application/json",
			bytes.NewBuffer(jsonBody),
		)
		require.NoError(t, err)
		defer resp2.Body.Close()

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp2.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp2.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Equal(t, "failed to register", response["error"])
	})

	main.Run("EdgeCases", func(t *testing.T) {
		_, err := setup(t)
		require.NoError(t, err)

		t.Run("VeryLongFields", func(t *testing.T) {
			// Arrange
			longString := string(make([]byte, 300)) // 300 character string
			requestBody := map[string]string{
				"login":    longString,
				"email":    longString,
				"password": longString,
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Act
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			assert.Equal(t, "validation failed", response["error"])
		})

		t.Run("SpecialCharacters", func(t *testing.T) {
			// Arrange
			requestBody := map[string]string{
				"login":    "user@#$%^&*()",
				"email":    "test@example.com",
				"password": "password123!@#$%",
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Act
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert - should succeed with special characters
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, "user@#$%^&*()", response["login"])
		})

		t.Run("UnicodeCharacters", func(t *testing.T) {
			// Arrange
			requestBody := map[string]string{
				"login":    "用户123",
				"email":    "тест@example.com",
				"password": "пароль123",
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Act
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert - should succeed with unicode characters
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, "用户123", response["login"])
		})

		t.Run("MissingContentType", func(t *testing.T) {
			// Arrange
			requestBody := map[string]string{
				"login":    "testuser",
				"email":    "test@example.com",
				"password": "testpassword123",
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			// Act - don't set Content-Type header
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"", // No content type
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert - The handler processes the request but fails at service level
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			assert.Equal(t, "failed to register", response["error"])
		})
	})

	main.Run("DatabaseConstraints", func(t *testing.T) {
		_, err := setup(t)
		require.NoError(t, err)

		t.Run("DuplicateEmail", func(t *testing.T) {
			// First user
			user1 := map[string]string{
				"login":    "user1",
				"email":    "same@example.com",
				"password": "password123",
			}

			jsonBody1, err := json.Marshal(user1)
			require.NoError(t, err)

			resp1, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody1),
			)
			require.NoError(t, err)
			resp1.Body.Close()
			assert.Equal(t, http.StatusCreated, resp1.StatusCode)

			// Second user with same email
			user2 := map[string]string{
				"login":    "user2",
				"email":    "same@example.com",
				"password": "password123",
			}

			jsonBody2, err := json.Marshal(user2)
			require.NoError(t, err)

			resp2, err := http.Post(
				fmt.Sprintf("%s/api/v1/signup", testapp.Addr),
				"application/json",
				bytes.NewBuffer(jsonBody2),
			)
			require.NoError(t, err)
			defer resp2.Body.Close()

			// Should fail due to duplicate email
			assert.Equal(t, http.StatusInternalServerError, resp2.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp2.Body).Decode(&response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			assert.Equal(t, "failed to register", response["error"])
		})
	})
}
