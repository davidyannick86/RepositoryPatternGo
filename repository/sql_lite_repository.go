package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/google/uuid"
)

const (
	insertUserQuery2 = `
    INSERT INTO users(id, name, email)
    VALUES(?, ?, ?)
    ON CONFLICT(id) DO UPDATE SET
      name       = excluded.name,
      email      = excluded.email;
`

	selectAllUsersQuery2 = `
    SELECT id, name, email
      FROM users;
`
)

// SqlliteRepository provides methods for user data operations using SQLite.
type SqlliteRepository struct {
	db *sql.DB
}

// NewSQLLiteRepository creates a new SQLite repository for user data.
func NewSQLLiteRepository(db *sql.DB) *SqlliteRepository {
	return &SqlliteRepository{db: db}
}

// AddUser adds a new user to the SQLite database and returns the created user.
func (r *SqlliteRepository) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	user.ID = uuid.New()
	_, err := r.db.ExecContext(ctx, insertUserQuery2, user.ID, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the SQLite database.
func (r *SqlliteRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	// Preallocate users slice with a reasonable capacity (e.g., 10)
	users := make([]domain.User, 0, 10)
	rows, err := r.db.QueryContext(ctx, selectAllUsersQuery2)
	if err != nil {
		return nil, fmt.Errorf("failed to query all users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, nil
}
