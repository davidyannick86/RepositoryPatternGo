package repository

import (
	"context"
	"database/sql"
	"time"

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

type sqlliteRepository struct {
	db *sql.DB
}

func NewSqlLiteRepository(db *sql.DB) UserRepository {
	return &sqlliteRepository{db: db}
}

func (r *sqlliteRepository) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	user.ID = uuid.New()

	_, err := r.db.ExecContext(ctx, insertUserQuery2, user.ID, user.Name, user.Email, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *sqlliteRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	rows, err := r.db.QueryContext(ctx, selectAllUsersQuery2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
