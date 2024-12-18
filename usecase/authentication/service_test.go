package authentication

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt"
	"github.com/mhdiiilham/dating-app/entity"
	authenticatormock "github.com/mhdiiilham/dating-app/usecase/authentication/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var mockAccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

type authenticationTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	mockUserRepository *authenticatormock.MockUserRepository
	mockJwtClient      *authenticatormock.MockJwtGenerator
	mockPasswordHasher *authenticatormock.MockHasher
}

func TestSignUp(t *testing.T) {
	suite.Run(t, new(authenticationTestSuite))
}

func (suite *authenticationTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockUserRepository = authenticatormock.NewMockUserRepository(suite.ctrl)
	suite.mockJwtClient = authenticatormock.NewMockJwtGenerator(suite.ctrl)
	suite.mockPasswordHasher = authenticatormock.NewMockHasher(suite.ctrl)
}

func (suite *authenticationTestSuite) TestSignUp() {
	testCases := []struct {
		condition   string
		request     SignUpRequest
		doMock      func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher)
		expected    *SignUpResponse
		expectedErr error
	}{
		{
			condition: "failed to validate email",
			request: SignUpRequest{
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Email:     "notvalidemail",
				Password:  faker.Password(),
			},
			doMock: func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher) {
			},
			expected:    nil,
			expectedErr: entity.ErrInvalidEmailAddress,
		},
		{
			condition: "email already registered",
			request: SignUpRequest{
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Email:     "hi@muhammadilham.xyz",
				Password:  faker.Password(),
			},
			doMock: func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher) {
				mockUserRepository.
					EXPECT().
					GetByEmail(gomock.Any(), "hi@muhammadilham.xyz").
					Return(&entity.User{
						FirstName: faker.FirstName(),
						LastName:  faker.LastName(),
						Email:     "hi@muhammadilham.xyz",
						Password:  faker.Password(),
					}, nil).
					Times(1)
			},
			expected:    nil,
			expectedErr: entity.ErrUserAlreadyExist,
		},
		{
			condition: "get user by email return unexpected error",
			request: SignUpRequest{
				FirstName: faker.FirstName(),
				LastName:  faker.LastName(),
				Email:     "hi@muhammadilham.xyz",
				Password:  faker.Password(),
			},
			doMock: func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher) {
				mockUserRepository.
					EXPECT().
					GetByEmail(gomock.Any(), "hi@muhammadilham.xyz").
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			expected:    nil,
			expectedErr: entity.ErrInternalServerError,
		},
		{
			condition: "success registered an user",
			request: SignUpRequest{
				FirstName: "Muhammad",
				LastName:  "Ilham",
				Email:     "hi@muhammadilham.xyz",
				Password:  "HelloWorld!",
			},
			doMock: func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "hi@muhammadilham.xyz").Return(nil, nil).Times(1)

				mockHasher.EXPECT().HashPassword("HelloWorld!").Return("HashedPassword", nil).Times(1)

				mockUserRepository.EXPECT().Save(gomock.Any(), &entity.User{
					FirstName: "Muhammad",
					LastName:  "Ilham",
					Email:     "hi@muhammadilham.xyz",
					Password:  "HashedPassword",
				}).Return("1", nil).Times(1)
				mockJwtClient.EXPECT().CreateAccessToken("1", "hi@muhammadilham.xyz").Return(mockAccessToken, nil).Times(1)
			},
			expected: &SignUpResponse{
				ID:          "1",
				Email:       "hi@muhammadilham.xyz",
				AccessToken: mockAccessToken,
			},
			expectedErr: nil,
		},
		{
			condition: "user repository save returned an unexpectec error",
			request: SignUpRequest{
				FirstName: "Muhammad",
				LastName:  "Ilham",
				Email:     "hi@muhammadilham.xyz",
				Password:  "HelloWorld!",
			},
			doMock: func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "hi@muhammadilham.xyz").Return(nil, nil).Times(1)

				mockHasher.EXPECT().HashPassword("HelloWorld!").Return("HashedPassword", nil).Times(1)

				mockUserRepository.EXPECT().Save(gomock.Any(), &entity.User{
					FirstName: "Muhammad",
					LastName:  "Ilham",
					Email:     "hi@muhammadilham.xyz",
					Password:  "HashedPassword",
				}).Return("", sql.ErrConnDone).Times(1)
			},
			expected:    nil,
			expectedErr: entity.ErrInternalServerError,
		},
		{
			condition: "hashing password failed",
			request: SignUpRequest{
				FirstName: "Muhammad",
				LastName:  "Ilham",
				Email:     "hi@muhammadilham.xyz",
				Password:  "HelloWorld!",
			},
			doMock: func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "hi@muhammadilham.xyz").Return(nil, nil).Times(1)

				mockHasher.EXPECT().HashPassword("HelloWorld!").Return("", errors.New("some error")).Times(1)
			},
			expected:    nil,
			expectedErr: entity.ErrInternalServerError,
		},
		{
			condition: "failed create access token",
			request: SignUpRequest{
				FirstName: "Muhammad",
				LastName:  "Ilham",
				Email:     "hi@muhammadilham.xyz",
				Password:  "HelloWorld!",
			},
			doMock: func(mockUserRepository *authenticatormock.MockUserRepository, mockJwtClient *authenticatormock.MockJwtGenerator, mockHasher *authenticatormock.MockHasher) {
				mockUserRepository.EXPECT().GetByEmail(gomock.Any(), "hi@muhammadilham.xyz").Return(nil, nil).Times(1)

				mockHasher.EXPECT().HashPassword("HelloWorld!").Return("HashedPassword", nil).Times(1)

				mockUserRepository.EXPECT().Save(gomock.Any(), &entity.User{
					FirstName: "Muhammad",
					LastName:  "Ilham",
					Email:     "hi@muhammadilham.xyz",
					Password:  "HashedPassword",
				}).Return("1", nil).Times(1)
				mockJwtClient.EXPECT().CreateAccessToken("1", "hi@muhammadilham.xyz").Return("", jwt.ErrInvalidKeyType).Times(1)
			},
			expected:    nil,
			expectedErr: entity.ErrInternalServerError,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.condition, func(t *testing.T) {
			assertion := assert.New(t)
			authenticator := NewService(suite.mockUserRepository, suite.mockJwtClient, suite.mockPasswordHasher)
			tc.doMock(suite.mockUserRepository, suite.mockJwtClient, suite.mockPasswordHasher)

			actual, actualErr := authenticator.SignUp(context.Background(), tc.request)
			assertion.Equal(tc.expected, actual)
			assertion.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *authenticationTestSuite) TearDownTest() {
	defer suite.ctrl.Finish()
}
