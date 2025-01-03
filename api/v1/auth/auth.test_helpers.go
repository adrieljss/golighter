package v1_auth

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/adrieljss/golighter/internal/testutils"
	"github.com/adrieljss/golighter/models"
	"github.com/adrieljss/golighter/platform"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T, app *platform.Application) {
	t.Helper()
	var currentUserCreds *UserRegister
	var currentUserRes *UserResponse
	t.Run("Register Flow", func(t *testing.T) {
		t.Run("User Register Invalid Fields", func(t *testing.T) {
			body, err := json.Marshal(UserRegister{
				Username: "1invalid-user",
				Email:    "invalid-email",
				Password: "",
			})
			assert.NoError(t, err)
			testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/register",
					Body:           body,
					ExpectedStatus: 400,
				},
			})
		})

		t.Run("User Register", func(t *testing.T) {
			gofakeit.Struct(&currentUserCreds)
			body, err := json.Marshal(currentUserCreds)

			assert.NoError(t, err)
			res := testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/register",
					Body:           body,
					ExpectedStatus: 200,
				},
			})[0]

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			json.Unmarshal(resBody, &currentUserRes)
		})

		t.Run("User Register Duplicated Email", func(t *testing.T) {
			body, err := json.Marshal(UserRegister{
				Username: "duplicated_user",
				Email:    currentUserCreds.Email,
				Password: "password",
			})

			assert.NoError(t, err)
			testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/register",
					Body:           body,
					ExpectedStatus: 400,
				},
			})
		})
	})

	t.Run("Login Flow", func(t *testing.T) {
		t.Run("User Login Invalid Fields", func(t *testing.T) {
			body, err := json.Marshal(UserLogin{
				Email:    "invalid-email",
				Password: "",
			})
			assert.NoError(t, err)
			testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/login",
					Body:           body,
					ExpectedStatus: 400,
				},
			})
		})

		t.Run("User Login Invalid Credentials", func(t *testing.T) {
			body, err := json.Marshal(UserLogin{
				Email:    "nonexistent@example.com",
				Password: "wrongpassword",
			})
			assert.NoError(t, err)
			testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/login",
					Body:           body,
					ExpectedStatus: 401,
				},
			})
		})

		t.Run("User Login Success", func(t *testing.T) {
			body, err := json.Marshal(UserLogin{
				Email:    currentUserCreds.Email,
				Password: currentUserCreds.Password,
			})
			assert.NoError(t, err)
			res := testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/login",
					Body:           body,
					ExpectedStatus: 200,
				},
			})[0]

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			var loginResponse UserResponse
			json.Unmarshal(resBody, &loginResponse)
			assert.Equal(t, currentUserRes.User, loginResponse.User)
		})

		t.Run("Me Endpoint", func(t *testing.T) {
			res := testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method: fiber.MethodGet,
					Path:   "/v1/users/@me",
					Headers: map[string]string{
						"Authorization": "Bearer " + currentUserRes.Token.AccessToken,
					},
					ExpectedStatus: 200,
				},
			})[0]

			resBody, err := io.ReadAll(res.Body)

			assert.NoError(t, err)

			var meResponse models.User
			json.Unmarshal(resBody, &meResponse)

			assert.Equal(t, currentUserRes.User, &meResponse)
		})
	})

	fmt.Printf("%+v\n", currentUserRes)

	t.Run("Token Flow", func(t *testing.T) {
		t.Run("Refresh Token Invalid Fields", func(t *testing.T) {
			body, err := json.Marshal(UserRefreshToken{
				RefreshToken: "",
			})
			assert.NoError(t, err)
			testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/refresh",
					Body:           body,
					ExpectedStatus: 400,
				},
			})
		})

		var refreshResponse AccTokenResponse
		t.Run("Refresh Access Token", func(t *testing.T) {
			body, err := json.Marshal(UserRefreshToken{
				RefreshToken: currentUserRes.Token.RefreshToken,
			})
			assert.NoError(t, err)
			res := testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method:         fiber.MethodPost,
					Path:           "/v1/auth/refresh",
					Body:           body,
					ExpectedStatus: 200,
				},
			})[0]

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			json.Unmarshal(resBody, &refreshResponse)
		})

		t.Run("Me Endpoint With New Access Token", func(t *testing.T) {
			res := testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
				{
					Method: fiber.MethodGet,
					Path:   "/v1/users/@me",
					Headers: map[string]string{
						"Authorization": "Bearer " + refreshResponse.AccessToken,
					},
					ExpectedStatus: 200,
				},
			})[0]

			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			var meResponse models.User
			json.Unmarshal(resBody, &meResponse)
			assert.Equal(t, currentUserRes.User, &meResponse)
		})
	})
}
