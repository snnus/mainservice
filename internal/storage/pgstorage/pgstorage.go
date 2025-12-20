package pgstorage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/snnus/mainservice/config"
	"github.com/snnus/mainservice/internal/models"
)

func NewConnection(cfg config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Addr,
		cfg.Postgres.Port,
		cfg.Postgres.DB)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres server: %w", err)
	}

	return db, nil
}

type PGStorage struct {
	db *sql.DB
}

func NewPGStorage(db *sql.DB) *PGStorage {
	return &PGStorage{db: db}
}

func (p *PGStorage) CreateServicePoint(ctx context.Context, sp models.NewServicePointRequest) (*models.ServicePoint, error) {
	query := `
		INSERT INTO service_points (name, short_name, office_number)
		VALUES ($1, $2, $3)
		RETURNING id, name, short_name, office_number, created_at, updated_at
	`

	var servicePoint models.ServicePoint

	err := p.db.QueryRowContext(ctx, query, sp.Name, sp.ShortName, sp.OfficeNumber).Scan(
		&servicePoint.ID,
		&servicePoint.Name,
		&servicePoint.ShortName,
		&servicePoint.OfficeNumber,
		&servicePoint.CreatedAt,
		&servicePoint.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create service point: %w", err)
	}
	return &servicePoint, nil
}

func (p *PGStorage) UpdateServicePoint(ctx context.Context, id string, sp models.NewServicePointRequest) (*models.ServicePoint, error) {
	query := `
		UPDATE service_points
		SET name = $1, short_name = $2, office_number = $3
		WHERE id = $4
		RETURNING id, name, short_name, office_number, created_at, updated_at
	`

	var servicePoint models.ServicePoint

	err := p.db.QueryRowContext(ctx, query, sp.Name, sp.ShortName, sp.OfficeNumber, id).Scan(
		&servicePoint.ID,
		&servicePoint.Name,
		&servicePoint.ShortName,
		&servicePoint.OfficeNumber,
		&servicePoint.CreatedAt,
		&servicePoint.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create service point: %w", err)
	}
	return &servicePoint, nil
}

func (p *PGStorage) DeleteServicePoint(ctx context.Context, id string) (*models.ServicePoint, error) {
	query := `
		DELETE FROM service_points
		WHERE id = $1
		RETURNING id, name, short_name, office_number, created_at, updated_at
	`

	var servicePoint models.ServicePoint

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&servicePoint.ID,
		&servicePoint.Name,
		&servicePoint.ShortName,
		&servicePoint.OfficeNumber,
		&servicePoint.CreatedAt,
		&servicePoint.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create service point: %w", err)
	}
	return &servicePoint, nil
}

func (p *PGStorage) GetServicePointByID(ctx context.Context, id string) (*models.ServicePoint, error) {
	query := `
		SELECT id, name, short_name, office_number, created_at, updated_at
		FROM service_points
		WHERE id = $1
	`

	var servicePoint models.ServicePoint

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&servicePoint.ID,
		&servicePoint.Name,
		&servicePoint.ShortName,
		&servicePoint.OfficeNumber,
		&servicePoint.CreatedAt,
		&servicePoint.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get service point: %w", err)
	}
	return &servicePoint, nil
}

func (p *PGStorage) GetShortNameById(ctx context.Context, id string) (string, error) {
	query := `
		SELECT short_name
		FROM service_points
		WHERE id = $1
	`

	var res string

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res,
	)

	if err != nil {
		return "", fmt.Errorf("failed to get short name: %w", err)
	}
	return res, nil
}

func (p *PGStorage) GetOfficeNumberById(ctx context.Context, id string) (string, error) {
	query := `
		SELECT office_number
		FROM service_points
		WHERE id = $1
	`

	var res string

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res,
	)

	if err != nil {
		return "", fmt.Errorf("failed to get office number: %w", err)
	}
	return res, nil
}
