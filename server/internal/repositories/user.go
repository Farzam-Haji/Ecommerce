package repositories

import (
	"context"
	"database/sql"

	"github.com/btalv/Ecommerce/server/internal/models"
)




type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository{
	return &PostgresUserRepository{
		db: db,
	}
}

func (ur *PostgresUserRepository) Create(ctx context.Context, user *models.User) (error) {
	row := ur.db.QueryRowContext(ctx, "INSERT INTO users(email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at",
		user.Email, user.PasswordHash)
	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return err
}

func (ur *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email=($1)",
		email)	
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return &user, err
}

func (ur *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row := ur.db.QueryRowContext(ctx, "SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id=($1)",
		id)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return &user, err
}