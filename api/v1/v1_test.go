package v1_api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	v1_auth "github.com/adrieljss/golighter/api/v1/auth"
	"github.com/adrieljss/golighter/internal/testutils"
	"github.com/adrieljss/golighter/models"
	"github.com/adrieljss/golighter/platform"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

var app *platform.Application

func TestMain(m *testing.M) {
	a := platform.App(true)
	app = &a
	defer app.CloseApp()
	SetupApiRoutes(app.FiberApp, app)
	os.Exit(m.Run())
}

func TestHealth(t *testing.T) {
	t.Run("1. Health Check", func(t *testing.T) {
		testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
			{
				Method:         fiber.MethodGet,
				Path:           "/v1/health",
				ExpectedStatus: 200,
				ExpectedBody:   "OK",
			},
		})
	})
}

func TestAuthFlow(t *testing.T) {
	var currentUserCreds *v1_auth.UserRegister
	var currentUserRes *v1_auth.UserResponse
	t.Run("1. Register Flow", func(t *testing.T) {
		t.Run("1.1. User Register Invalid Fields", func(t *testing.T) {
			body, err := json.Marshal(v1_auth.UserRegister{
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

		t.Run("1.2. User Register", func(t *testing.T) {
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

		t.Run("1.3. User Register Duplicated Email", func(t *testing.T) {
			body, err := json.Marshal(v1_auth.UserRegister{
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

	t.Run("2. Login Flow", func(t *testing.T) {
		t.Run("2.1. User Login Invalid Fields", func(t *testing.T) {
			body, err := json.Marshal(v1_auth.UserLogin{
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

		t.Run("2.2. User Login Invalid Credentials", func(t *testing.T) {
			body, err := json.Marshal(v1_auth.UserLogin{
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

		t.Run("2.3. User Login Success", func(t *testing.T) {
			body, err := json.Marshal(v1_auth.UserLogin{
				Email:    currentUserCreds.Email,
				Password: currentUserCreds.Password, // Default password from gofakeit
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

			var loginResponse v1_auth.UserResponse
			json.Unmarshal(resBody, &loginResponse)
			assert.Equal(t, currentUserRes.User, loginResponse.User)
		})

		t.Run("2.4. @me Endpoint", func(t *testing.T) {
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

	t.Run("3. Token Flow", func(t *testing.T) {
		t.Run("3.1. Refresh Token Invalid Fields", func(t *testing.T) {
			body, err := json.Marshal(v1_auth.UserRefreshToken{
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

		var refreshResponse v1_auth.AccTokenResponse
		t.Run("3.2. Refresh Access Token", func(t *testing.T) {
			body, err := json.Marshal(v1_auth.UserRefreshToken{
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

		t.Run("3.3. @me Endpoint With New Access Token", func(t *testing.T) {
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
