// internal/user/repo.go
package user

import (
	"context"

	"sumunar-pos-core/pkg/db"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByFullname(ctx context.Context, fullname string) (*User, error)
	Create(ctx context.Context, user *User) error
	FindAll(ctx context.Context, limit, offset int) ([]*User, int, error)
	UpdateLastLogin(ctx context.Context, id string) error
}

type userRepo struct {
	db db.DBTX
}

func NewUserRepository(db db.DBTX) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, fullname, email, password, role, last_login, updated_at, is_active FROM users WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)

	var u User
	err := row.Scan(
		&u.ID, &u.Fullname, &u.Email, &u.Password,
		&u.Role, &u.LastLogin, &u.UpdatedAt, &u.IsActive,
	)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &u, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, fullname, email, password, role, last_login, updated_at, is_active FROM users WHERE email=$1`
	row := r.db.QueryRow(ctx, query, email)

	var u User
	err := row.Scan(
		&u.ID, &u.Fullname, &u.Email, &u.Password,
		&u.Role, &u.LastLogin, &u.UpdatedAt, &u.IsActive,
	)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &u, nil
}

func (r *userRepo) FindByFullname(ctx context.Context, fullname string) (*User, error) {
	query := `SELECT id, fullname, email, password, role, last_login, updated_at, is_active FROM users WHERE fullname=$1`
	row := r.db.QueryRow(ctx, query, fullname)

	var u User
	err := row.Scan(
		&u.ID, &u.Fullname, &u.Email, &u.Password,
		&u.Role, &u.LastLogin, &u.UpdatedAt, &u.IsActive,
	)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &u, nil
}

func (r *userRepo) FindAll(ctx context.Context, limit, offset int) ([]*User, int, error) {
	const queryUsers = `
		SELECT id, fullname, email, google_id, password, picture, provider,
		       last_login, role, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	const queryCount = `SELECT COUNT(*) FROM users`

	// Query users
	rows, err := r.db.Query(ctx, queryUsers, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.Fullname, &u.Email, &u.GoogleID, &u.Password, &u.Picture,
			&u.Provider, &u.LastLogin, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, &u)
	}

	// Query count
	var total int
	if err := r.db.QueryRow(ctx, queryCount).Scan(&total); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO users (id, fullname, email, password, role, provider, last_login, created_by, created_at, updated_by, updated_at, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Exec(ctx, query,
		u.ID, u.Fullname, u.Email, u.Password,
		u.Role, u.Provider, u.LastLogin, u.CreatedBy, u.CreatedAt, u.UpdatedBy, u.UpdatedAt, u.IsActive,
	)
	return err
}

func (r *userRepo) UpdateLastLogin(ctx context.Context, id string) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE users SET last_login = NOW(), updated_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}
