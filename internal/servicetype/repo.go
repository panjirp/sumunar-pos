package servicetype

import (
	"context"
	"sumunar-pos-core/pkg/db"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *ServiceType) error
	FindByID(ctx context.Context, id string) (*ServiceType, error)
	FindAll(ctx context.Context, limit, offset int) ([]*ServiceType, int, error)
	Update(ctx context.Context, service *ServiceType) error
	Delete(ctx context.Context, id string) error
}

type serviceRepo struct {
	db db.DBTX
}

func NewServiceRepository(db db.DBTX) ServiceRepository {
	return &serviceRepo{db}
}

func (r *serviceRepo) Create(ctx context.Context, serviceType *ServiceType) error {
	query := `
		INSERT INTO service_types (id, name, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		serviceType.ID,
		serviceType.Name,
		serviceType.IsActive,
		serviceType.CreatedAt,
		serviceType.CreatedBy,
	)
	return err
}

func (r *serviceRepo) FindByID(ctx context.Context, id string) (*ServiceType, error) {
	query := `SELECT id, name, is_active, created_at, created_by, updated_at, updated_by FROM service_types WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var service ServiceType
	err := row.Scan(
		&service.ID,
		&service.Name,
		&service.IsActive,
		&service.CreatedAt,
		&service.CreatedBy,
		&service.UpdatedAt,
		&service.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepo) FindAll(ctx context.Context, limit, offset int) ([]*ServiceType, int, error) {
	query := `SELECT id, name, is_active, created_at, created_by, updated_at, updated_by FROM service_types LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var services []*ServiceType
	for rows.Next() {
		var s ServiceType
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.IsActive,
			&s.CreatedAt,
			&s.CreatedBy,
			&s.UpdatedAt,
			&s.UpdatedBy,
		); err != nil {
			return nil, 0, err
		}
		services = append(services, &s)
	}

	// total count
	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM service_types`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *serviceRepo) Update(ctx context.Context, service *ServiceType) error {
	query := `
		UPDATE service_types SET name = $1, is_active = $2,
		updated_at = $3, updated_by = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query,
		service.Name,
		service.IsActive,
		service.UpdatedAt,
		service.UpdatedBy,
		service.ID,
	)
	return err
}

func (r *serviceRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM service_types WHERE id = $1`, id)
	return err
}
