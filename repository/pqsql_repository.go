package repository

import (
	"context"
	"fmt"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertUserQuery = `
    INSERT INTO users (id, name, email)
    VALUES ($1, $2, $3)
    ON CONFLICT (id) DO UPDATE
      SET name       = EXCLUDED.name,
          email      = EXCLUDED.email;
    `

	selectAllUsersQuery = `SELECT id, name, email FROM users`
)

// PsqlRepository provides methods for interacting with the users table in a PostgreSQL database.
type PsqlRepository struct {
	pool *pgxpool.Pool
}

// NewPsqlRepository creates a new instance of PsqlRepository with the given pgxpool.Pool.
func NewPsqlRepository(pool *pgxpool.Pool) *PsqlRepository {
	return &PsqlRepository{pool: pool}
}

// AddUser inserts a new user into the database and returns the created user.
func (r *PsqlRepository) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	user.ID = uuid.New()

	_, err := r.pool.Exec(ctx, insertUserQuery, user.ID, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert user query: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *PsqlRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0)
	rows, err := r.pool.Query(ctx, selectAllUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select all users query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
