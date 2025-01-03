package v1_users

import (
	"context"

	"github.com/adrieljss/golighter/models"
	"github.com/adrieljss/golighter/platform"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type userRepository struct {
	*platform.Application
}

func NewUserRepository(app *platform.Application) *userRepository {
	return &userRepository{app}
}

func (r *userRepository) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// creates a new user in the database, hashing the password before storing it
func (r *userRepository) Create(ctx context.Context, username string, email string, password string) (*models.User, error) {
	// hashing is done in the app, not sent to the db
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING *`
	hashedPassword, err := r.hashPassword(password)
	if err != nil {
		return nil, err
	}

	row, err := r.Db.Query(ctx, query, username, email, hashedPassword)

	if err != nil {
		return nil, err
	}

	var newUser models.User
	newUser, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])
	return &newUser, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	row, err := r.Db.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	var user models.User
	user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])
	return &user, err
}

func (r *userRepository) GetByUID(ctx context.Context, uid string) (*models.User, error) {
	query := `SELECT * FROM users WHERE uid = $1`

	row, err := r.Db.Query(ctx, query, uid)
	if err != nil {
		return nil, err
	}

	var user models.User
	user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])
	return &user, err
}

// updates non-sensitive user information, such as username
func (r *userRepository) UpdateProfile(ctx context.Context, uid string, username, email string) error {
	query := `UPDATE users SET username = $1 WHERE uid = $2 RETURNING *`
	_, err := r.Db.Exec(ctx, query, username, email, uid)
	return err
}

func (r *userRepository) UpdateEmail(ctx context.Context, uid string, email string) error {
	query := `UPDATE users SET email = $1 WHERE uid = $2`
	_, err := r.Db.Exec(ctx, query, email, uid)
	return err
}

// updates the user's password, hashing it before storing it
func (r *userRepository) UpdatePassword(ctx context.Context, uid string, passwordRaw string) error {
	query := `UPDATE users SET password_hash = $1 WHERE uid = $2`
	hashedPassword, err := r.hashPassword(passwordRaw)
	if err != nil {
		return err
	}
	_, err = r.Db.Exec(ctx, query, hashedPassword, uid)
	return err
}

// dangerous operation, deletes the user from the database
func (r *userRepository) Delete(ctx context.Context, uid string) error {
	query := `DELETE FROM users WHERE uid = $1`
	_, err := r.Db.Exec(ctx, query, uid)
	return err
}
