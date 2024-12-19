package restful

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhdiiilham/dating-app/entity"
	"github.com/mhdiiilham/dating-app/usecase/authentication"
)

// SignUpRequest struct holds the JSON format for signup.
type SignUpRequest struct {
	Fistname string `json:"fistname,omitempty"`
	Lastname string `json:"lastname,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// SignInRequest struct ...
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResponse struct is the JSON response when user is successfully signed up.
type SignUpResponse struct {
}

// Authenticator interface ...
type Authenticator interface {
	SignUp(ctx context.Context, request authentication.SignUpRequest) (credential *authentication.AccessResponse, err error)
	SignIn(ctx context.Context, email, password string) (*authentication.AccessResponse, error)
}

// HandleUserSignUp function handle REST endpoint for user registrations.
func HandleUserSignUp(authenticatorService Authenticator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var request SignUpRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse(err))
		}

		ctx := c.Request().Context()
		access, err := authenticatorService.SignUp(ctx, authentication.SignUpRequest{
			FirstName: request.Fistname,
			LastName:  request.Lastname,
			Email:     request.Email,
			Password:  request.Password,
		})
		if err != nil {
			if errors.Is(err, entity.ErrInternalServerError) {
				return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse(err))
			}

			return c.JSON(http.StatusBadRequest, BadRequestErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, SuccessResponse(access, "signup success", http.StatusCreated))
	}
}

// HandleUserSignIn function ...
func HandleUserSignIn(authenticatorService Authenticator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var signInRequest SignInRequest
		if err := c.Bind(&signInRequest); err != nil {
			return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse(err))
		}

		ctx := c.Request().Context()
		access, err := authenticatorService.SignIn(ctx, signInRequest.Email, signInRequest.Password)
		if err != nil {
			if errors.Is(err, entity.ErrInternalServerError) {
				return c.JSON(http.StatusInternalServerError, InternalServerErrorResponse(err))
			}

			return c.JSON(http.StatusBadRequest, BadRequestErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, SuccessResponse(access, "Sign-In Success", http.StatusOK))
	}
}
