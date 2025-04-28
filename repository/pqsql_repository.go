package repository

import (
	"context"

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

type psqlRepository struct {
	pool *pgxpool.Pool
}

func NewPsqlRepository(pool *pgxpool.Pool) UserRepository {
	return &psqlRepository{pool: pool}
}

func (r *psqlRepository) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {

	user.ID = uuid.New()

	_, err := r.pool.Exec(ctx, insertUserQuery, user.ID, user.Name, user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *psqlRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	rows, err := r.pool.Query(ctx, selectAllUsersQuery)
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
