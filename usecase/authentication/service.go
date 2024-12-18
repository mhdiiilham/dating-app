package authentication

//go:generate mockgen -destination=mock/service.go -package=authenticatormock . UserRepository,JwtGenerator,Hasher

import (
	"context"
	"net/mail"

	"github.com/mhdiiilham/dating-app/entity"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// UserRepository holds the repository contracts for user entity.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (user *entity.User, err error)
	Save(ctx context.Context, user *entity.User) (ID string, err error)
}

// JwtGenerator interface is the interface contracts between service and the Jwt Generator client.
type JwtGenerator interface {
	CreateAccessToken(userID, email string) (accessToken string, err error)
}

// Hasher interface is the interface contracts between authenticator service and password encryption helper.
type Hasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hashedPassword string) bool
}

// Service struct holds the implementation of Authenticator service.
type Service struct {
	userRepository UserRepository
	jwtClient      JwtGenerator
	passwordHaser  Hasher
}

// NewService function return new instance of Authenticator service.
func NewService(userRepository UserRepository, jwtClient JwtGenerator, passwordHaser Hasher) *Service {
	return &Service{
		userRepository: userRepository,
		jwtClient:      jwtClient,
		passwordHaser:  passwordHaser,
	}
}

// SignUp function
func (s *Service) SignUp(ctx context.Context, request SignUpRequest) (credential *SignUpResponse, err error) {
	if _, err = mail.ParseAddress(request.Email); err != nil {
		return nil, entity.ErrInvalidEmailAddress
	}

	existingUser, err := s.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		logrus.Warnf("[Authenticaion.SignUp] unexpected error from userRepository.GetByEmail: %v", err)
		return nil, entity.ErrInternalServerError
	}

	if existingUser != nil {
		return nil, entity.ErrUserAlreadyExist
	}

	hashedPassword, err := s.passwordHaser.HashPassword(request.Password)
	if err != nil {
		logrus.Warnf("[Authenticaion.SignUp] failed to hash password: %v", err)
		return nil, entity.ErrInternalServerError
	}

	user := entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  hashedPassword,
	}

	userID, err := s.userRepository.Save(ctx, &user)
	if err != nil {
		log.Warnf("[Authenticaion.SignUp] failed to save user: %v", err)
		return nil, entity.ErrInternalServerError
	}

	user.ID = userID
	accessToken, err := s.jwtClient.CreateAccessToken(user.ID, user.Email)
	if err != nil {
		log.Warnf("[Authenticaion.SignUp] failed to save user: %v", err)
		return nil, entity.ErrInternalServerError
	}

	return &SignUpResponse{
		ID:          user.ID,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}
