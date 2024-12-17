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

// creates a new user in the database, hashing the password before storing it
func (r *userRepository) Create(ctx context.Context, username string, email string, password string) (*models.User, error) {
	// hashing is done in the app, not sent to the db
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING *`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, err
	}

	row, err := r.Db.Query(ctx, query, username, email, string(hashedPassword))

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
