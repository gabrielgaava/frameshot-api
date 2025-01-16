package repository

import (
	"context"
	"example/web-service-gin/src/adapters/storage/postgres"
	"example/web-service-gin/src/core"
	"example/web-service-gin/src/core/entity"
	"log"

	sq "github.com/Masterminds/squirrel"
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
		Columns("user_id", "user_email", "video_size", "video_key", "zip_output_key", "status", "created_at").
		Values(request.UserId, request.UserEmail, request.VideoSize, request.VideoKey, request.ZipOutputKey, request.Status, request.CreatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = repository.db.QueryRow(ctx, sql, args...).Scan(
		&request.ID,
		&request.UserId,
		&request.UserEmail,
		&request.VideoSize,
		&request.VideoKey,
		&request.ZipOutputKey,
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
			&request.UserEmail,
			&request.VideoSize,
			&request.VideoKey,
			&request.ZipOutputKey,
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

func (repository *PGRequestRepository) UpdateStatusByVideoKey(ctx context.Context, status string, videoKey string) (*entity.Request, error) {

	var request = entity.Request{}
	condition := sq.Eq{"video_key": videoKey}
	updatedData := map[string]interface{}{
		"status": status,
	}

	query := repository.db.QueryBuilder.Update("requests").
		SetMap(updatedData).
		Where(condition).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	repository.db.QueryRow(ctx, sql, args...).Scan(
		&request.ID,
		&request.UserId,
		&request.UserEmail,
		&request.VideoSize,
		&request.VideoKey,
		&request.ZipOutputKey,
		&request.Status,
		&request.CreatedAt,
		&request.FinishedAt,
	)

	return &request, nil
}
