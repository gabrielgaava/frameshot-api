package repository

import (
	"context"
	"example/web-service-gin/src/adapters/storage/postgres"
	"example/web-service-gin/src/core"
	"example/web-service-gin/src/core/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const ReturnSuffix = "RETURNING *"

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
		Suffix(ReturnSuffix)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := repository.db.QueryRow(ctx, sql, args...)
	request, err = mapRowToRequest(row)

	if err != nil {
		if errCode := repository.db.ErrorCode(err); errCode == "23505" {
			return nil, core.ErrConflictingData
		}
		return nil, err
	}

	return request, nil
}

func (repository *PGRequestRepository) GetById(ctx context.Context, id uint64) (*entity.Request, error) {
	condition := sq.Eq{"id": id}
	query := repository.db.QueryBuilder.Select("*").
		From("requests").
		Where(condition).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := repository.db.QueryRow(ctx, sql, args...)
	request, _ := mapRowToRequest(row)

	return request, nil
}

func (repository *PGRequestRepository) GetAllUserRequests(ctx context.Context, userId string) ([]entity.Request, error) {
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

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	userRequests, _ = mapRowListToRequest(rows)

	return userRequests, nil
}

func (repository *PGRequestRepository) UpdateRequest(ctx context.Context, request *entity.Request) (*entity.Request, error) {
	condition := sq.Eq{"id": request.ID}
	updatedData := map[string]interface{}{
		"user_email":     request.UserEmail,
		"video_size":     request.VideoSize,
		"video_key":      request.VideoKey,
		"zip_output_key": request.ZipOutputKey,
		"status":         request.Status,
		"finished_at":    request.FinishedAt,
	}

	query := repository.db.QueryBuilder.Update("requests").
		SetMap(updatedData).
		Where(condition).
		Suffix(ReturnSuffix)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := repository.db.QueryRow(ctx, sql, args...)
	updatedRequest, err := mapRowToRequest(row)

	if err != nil {
		return nil, err
	}

	return updatedRequest, nil
}

func (repository *PGRequestRepository) UpdateStatusByVideoKey(ctx context.Context, status string, videoKey string) (*entity.Request, error) {
	condition := sq.Eq{"video_key": videoKey}
	updatedData := map[string]interface{}{
		"status": status,
	}

	query := repository.db.QueryBuilder.Update("requests").
		SetMap(updatedData).
		Where(condition).
		Suffix(ReturnSuffix)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := repository.db.QueryRow(ctx, sql, args...)
	updatedRequest, err := mapRowToRequest(row)

	if err != nil {
		return nil, err
	}

	return updatedRequest, nil
}

// Map a row of database data to domain entity Request model
func mapRowToRequest(row pgx.Row) (*entity.Request, error) {
	var request RequestModel

	err := row.Scan(
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

	return modelToEntity(request), nil
}

// Map the rows (list) to domain entity Request model
func mapRowListToRequest(rows pgx.Rows) ([]entity.Request, error) {
	var requests []entity.Request
	for rows.Next() {
		request, err := mapRowToRequest(rows)
		if err != nil {
			return nil, err
		}
		requests = append(requests, *request)
	}

	return requests, nil
}

func modelToEntity(model RequestModel) *entity.Request {

	var data = entity.Request{
		ID:        model.ID,
		UserId:    model.UserId,
		UserEmail: model.UserEmail,
		VideoSize: model.VideoSize,
		VideoKey:  model.VideoKey,
		Status:    entity.RequestStatus(model.Status),
		CreatedAt: model.CreatedAt,
	}

	if model.ZipOutputKey.Valid {
		data.ZipOutputKey = model.ZipOutputKey.String
	}

	if model.FinishedAt.Valid {
		data.FinishedAt = model.FinishedAt.Time
	}

	return &data
}
