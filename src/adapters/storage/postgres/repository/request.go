package repository

import (
	"context"
	"example/web-service-gin/src/adapters/storage/postgres"
	"example/web-service-gin/src/core"
	"example/web-service-gin/src/core/entity"
	"log"
)

// PGRequestRepository implements port.RequestRepository interface
// * and provides access to the postgres database/**
type PGRequestRepository struct {
	db *postgres.DB
}

// NewPGRequestRepository creates a new request storage instance for postgres
func NewPGRequestRepository(db *postgres.DB) *PGRequestRepository {
	return &PGRequestRepository{
		db,
	}
}

// CreateRequest creates a new request register in the database
func (repository *PGRequestRepository) CreateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error) {

	query := repository.db.QueryBuilder.Insert("requests").
		Columns("user_id", "video_url", "zip_output_url", "status", "created_at").
		Values(request.UserId, request.VideoUrl, request.ZipOutputUrl, request.Status, request.CreatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = repository.db.QueryRow(ctx, sql, args...).Scan(
		&request.ID,
		&request.UserId,
		&request.VideoUrl,
		&request.ZipOutputUrl,
		&request.Status,
		&request.CreatedAt,
		&request.FinishedAt,
	)

	if err != nil {
		if errCode := repository.db.ErrorCode(err); errCode == "23505" {
			return nil, core.ErrConflictingData
		}
		return nil, err
	}

	return request, nil
}

func (repository *PGRequestRepository) GetAllUserRequests(ctx context.Context, userId string) ([]entity.Request, error) {

	log.Println("GET ALL DATABASE")

	var request entity.Request
	var userRequests []entity.Request

	query := repository.db.QueryBuilder.Select("*").
		From("requests").
		OrderBy("created_at")

	// Create SQL Statement
	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	// Runs the SQL and return in rows
	rows, err := repository.db.Query(ctx, sql, args...)

	log.Println(rows)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&request.ID,
			&request.UserId,
			&request.VideoUrl,
			&request.ZipOutputUrl,
			&request.Status,
			&request.CreatedAt,
			&request.FinishedAt,
		)
		if err != nil {
			return nil, err
		}

		userRequests = append(userRequests, request)
	}

	return userRequests, nil

}

func (repository *PGRequestRepository) UpdateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	return nil, nil
}
