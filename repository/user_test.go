package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/mhdiiilham/dating-app/entity"
	"github.com/mhdiiilham/dating-app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type userRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *sql.DB
}

func (suite *userRepositoryTestSuite) SetupTest() {
	var err error

	suite.db, suite.mock, err = sqlmock.New()
	suite.NoError(err)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(userRepositoryTestSuite))
}

func (suite *userRepositoryTestSuite) TestGetByEmail() {
	dbTime, err := time.Parse(time.RFC3339, "2019-10-05T14:45:05+07:00")
	suite.NoError(err)

	testCases := []struct {
		condition    string
		doMock       func(db *sql.DB, mock sqlmock.Sqlmock)
		email        string
		expectedUser *entity.User
		expecterErr  error
	}{
		{
			condition: "user is found",
			doMock: func(db *sql.DB, mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "first_name", "last_name", "email", "password", "created_at", "updated_at", "delete_at"}).
					AddRow("1", "Muhammad", "Ilham", "hi@muhammadilham.xyz", "$2a$12$N6DyTiEfloOITL5MAecrMOk2DTV/ejLtuHpjmlYY1NTOVSDovF1Ay", dbTime, dbTime, nil)

				mock.
					ExpectQuery(regexp.QuoteMeta(repository.UserFindByEmail)).
					WithArgs("hi@muhammadilham.xyz").
					WillReturnRows(rows)
			},
			email: "hi@muhammadilham.xyz",
			expectedUser: &entity.User{
				ID:        "1",
				FirstName: "Muhammad",
				LastName:  "Ilham",
				Email:     "hi@muhammadilham.xyz",
				Password:  "$2a$12$N6DyTiEfloOITL5MAecrMOk2DTV/ejLtuHpjmlYY1NTOVSDovF1Ay",
				CreatedAt: dbTime,
				UpdatedAt: dbTime,
				DeletedAt: nil,
			},
			expecterErr: nil,
		},
		{
			condition: "user is not found",
			doMock: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(repository.UserFindByEmail)).
					WithArgs("hi@muhammadilham.xyz").
					WillReturnError(sql.ErrNoRows)
			},
			email:        "hi@muhammadilham.xyz",
			expectedUser: nil,
			expecterErr:  nil,
		},
		{
			condition: "Unexpected error",
			doMock: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(repository.UserFindByEmail)).
					WithArgs("hi@muhammadilham.xyz").
					WillReturnError(sql.ErrConnDone)
			},
			email:        "hi@muhammadilham.xyz",
			expectedUser: nil,
			expecterErr:  sql.ErrConnDone,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.condition, func(t *testing.T) {
			userRepository := repository.NewUser(suite.db)
			tc.doMock(suite.db, suite.mock)

			ctx := context.Background()
			assert := assert.New(t)
			actual, actualErr := userRepository.GetByEmail(ctx, tc.email)
			assert.Equal(tc.expectedUser, actual)
			assert.Equal(tc.expecterErr, actualErr)
		})
	}
}

func (suite *userRepositoryTestSuite) TestSaveUser() {
	testCases := []struct {
		condition   string
		user        *entity.User
		doMock      func(db *sql.DB, mock sqlmock.Sqlmock)
		expected    string
		expectedErr error
	}{
		{
			condition: "success save user",
			user:      &entity.User{FirstName: "Muhammad", LastName: "Ilham", Email: "hi@muhammadilham.xyz", Password: "hashedpassword"},
			doMock: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(repository.UserSave)).
					WithArgs("Muhammad", "Ilham", "hi@muhammadilham.xyz", "hashedpassword").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
			},
			expected:    "1",
			expectedErr: nil,
		},
		{
			condition: "unique constraint validation error",
			user:      &entity.User{FirstName: "Muhammad", LastName: "Ilham", Email: "hi@muhammadilham.xyz", Password: "hashedpassword"},
			doMock: func(db *sql.DB, mock sqlmock.Sqlmock) {

				mock.
					ExpectQuery(regexp.QuoteMeta(repository.UserSave)).
					WithArgs("Muhammad", "Ilham", "hi@muhammadilham.xyz", "hashedpassword").
					WillReturnError(&pq.Error{Code: pq.ErrorCode("23505")})
			},
			expected:    "",
			expectedErr: entity.ErrUserAlreadyExist,
		},
		{
			condition: "Unexpected error",
			user:      &entity.User{FirstName: "Muhammad", LastName: "Ilham", Email: "hi@muhammadilham.xyz", Password: "hashedpassword"},
			doMock: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(regexp.QuoteMeta(repository.UserSave)).
					WithArgs("Muhammad", "Ilham", "hi@muhammadilham.xyz", "hashedpassword").
					WillReturnError(&pq.Error{Code: pq.ErrorCode("P0000")})
			},
			expected:    "",
			expectedErr: entity.ErrInternalServerError,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.condition, func(t *testing.T) {
			assert := assert.New(t)

			ctx := context.Background()
			userRepository := repository.NewUser(suite.db)
			tc.doMock(suite.db, suite.mock)

			actual, actualErr := userRepository.Save(ctx, tc.user)
			assert.Equal(tc.expected, actual)
			assert.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *userRepositoryTestSuite) TearDownTest() {
	defer suite.db.Close()
}
