package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/mhdiiilham/dating-app/entity"
	log "github.com/sirupsen/logrus"
)

// Query use in UserRepository.
var (
	UserFindByEmail = `SELECT id, first_name, last_name, email, password from users where email = $1 LIMIT 1;`
	UserSave        = `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING "id";`
)

// User struct contains db connections.
type User struct {
	dbClient *sql.DB
}

// NewUser function get a new instance of User Repository.
func NewUser(db *sql.DB) *User {
	return &User{dbClient: db}
}

// GetByEmail function retrive an user by it's email.
func (r *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	row := r.dbClient.QueryRowContext(ctx, UserFindByEmail, email)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		log.Errorf("[UserRepository] Unexpected error: %v", err)
		return nil, err
	}

	return &user, nil
}

// Save function insert the user to database.
// This function will not check for duplication, instead return conflict error from db.
func (r *User) Save(ctx context.Context, user *entity.User) (ID string, err error) {
	if err := r.dbClient.QueryRowContext(ctx, UserSave, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID); err != nil {
		if pgsqlErr, ok := err.(*pq.Error); ok {
			if pgsqlErr.Code == "23505" {
				return "", entity.ErrUserAlreadyExist
			}
		}

		log.Errorf("[UserRepository.Save] Unexpected error: %v", err)
		return "", entity.ErrInternalServerError
	}

	return user.ID, nil
}
