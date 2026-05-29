package users

import (
	"context"
	"errors"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/EdgarPFernandes/SecureVitalV2_backend/internal/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	ListUsers(ctx context.Context) ([]repo.ListUsersRow, error)

	RegisterUser(
		ctx context.Context,
		payload RegisterUserPayload,
	) (repo.User, string, error)

	Login(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
}
type svc struct {
	repoQuerier repo.Querier
	repoQueries *repo.Queries
	db          *pgx.Conn
}

func NewService(repoQuerier repo.Querier, repoQueries *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repoQuerier: repoQuerier,
		repoQueries: repoQueries,
		db:          db,
	}
}

func (s *svc) ListUsers(ctx context.Context) ([]repo.ListUsersRow, error) {
	return s.repoQuerier.ListUsers(ctx)
}

func (s *svc) RegisterUser(
	ctx context.Context,
	payload RegisterUserPayload,
) (repo.User, string, error) {

	hash, err := auth.HashPassword(payload.Password)
	if err != nil {
		return repo.User{}, "", err
	}

	user, err := s.repoQueries.CreateUser(
		ctx,
		repo.CreateUserParams{
			Name:         payload.Name,
			Email:        payload.Email,
			PhoneNumber:  pgtype.Text{String: payload.PhoneNumber, Valid: payload.PhoneNumber != ""},
			PasswordHash: hash,
			Role:         payload.Role,
			Photo:        pgtype.Text{String: payload.Photo, Valid: payload.Photo != ""},
		},
	)

	if err != nil {
		return repo.User{}, "", err
	}

	token, err := auth.GenerateToken(
		user.ID,
		user.Role,
	)

	if err != nil {
		return repo.User{}, "", err
	}

	return user, token, nil
}

func (s *svc) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := s.repoQueries.GetUserByEmail(
		ctx,
		email,
	)

	if err != nil {
		return "", err
	}

	err = auth.CheckPasswordHash(
		password,
		user.PasswordHash,
	)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = s.repoQueries.UpdateLastAccess(
		ctx,
		user.ID,
	)

	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(
		user.ID,
		user.Role,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
