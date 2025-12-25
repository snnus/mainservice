package spstorage

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"

	_ "github.com/lib/pq"
	"github.com/snnus/mainservice/config"
	"github.com/snnus/mainservice/internal/models"
)

func NewConnection(cfg *config.Config) (*sql.DB, error) {
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

type SPStorage struct {
	db      *sql.DB
	nShards uint32
}

func NewSPStorage(cfg *config.Config) (*SPStorage, func() error, error) {
	db, err := NewConnection(cfg)
	if err != nil {
		return nil, db.Close, err
	}
	return &SPStorage{db: db, nShards: cfg.Postgres.NShards}, db.Close, nil
}

func (p *SPStorage) GetHash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (p *SPStorage) GetShard(h uint32) uint32 {
	return h%p.nShards + 1
}

func (p *SPStorage) UpsertServicePoint(ctx context.Context, id string, sp models.NewServicePointRequest) (*models.ServicePoint, error) {
	shard := p.GetShard(p.GetHash(id))
	query := fmt.Sprintf(`
		INSERT INTO shard_%d.service_points (id, name, short_name, office_number)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name, 
			short_name = EXCLUDED.short_name, 
			office_number = EXCLUDED.office_number
		RETURNING id, name, short_name, office_number, created_at, updated_at
	`, shard)

	var servicePoint models.ServicePoint

	err := p.db.QueryRowContext(ctx, query, id, sp.Name, sp.ShortName, sp.OfficeNumber).Scan(
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

func (p *SPStorage) DeleteServicePoint(ctx context.Context, id string) (*models.ServicePoint, error) {
	shard := p.GetShard(p.GetHash(id))
	query := fmt.Sprintf(`
		DELETE FROM shard_%d.service_points
		WHERE id = $1
		RETURNING id, name, short_name, office_number, created_at, updated_at
	`, shard)

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

func (p *SPStorage) GetServicePointByID(ctx context.Context, id string) (*models.ServicePoint, error) {
	shard := p.GetShard(p.GetHash(id))
	query := fmt.Sprintf(`
		SELECT id, name, short_name, office_number, created_at, updated_at
		FROM shard_%d.service_points
		WHERE id = $1
	`, shard)

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

func (p *SPStorage) GetShortNameById(ctx context.Context, id string) (string, error) {
	shard := p.GetShard(p.GetHash(id))
	query := fmt.Sprintf(`
		SELECT short_name
		FROM shard_%d.service_points
		WHERE id = $1
	`, shard)

	var res string

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res,
	)

	if err != nil {
		return "", fmt.Errorf("failed to get short name: %w", err)
	}
	return res, nil
}

func (p *SPStorage) GetOfficeNumberById(ctx context.Context, id string) (string, error) {
	shard := p.GetShard(p.GetHash(id))
	query := fmt.Sprintf(`
		SELECT office_number
		FROM shard_%d.service_points
		WHERE id = $1
	`, shard)

	var res string

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&res,
	)

	if err != nil {
		return "", fmt.Errorf("failed to get office number: %w", err)
	}
	return res, nil
}
