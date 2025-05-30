package postgres

import (
	"context"
	"strings"

	"github.com/flexGURU/zeiba-glam/backend/internal/postgres/generated"
	"github.com/flexGURU/zeiba-glam/backend/internal/repository"
	"github.com/flexGURU/zeiba-glam/backend/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ repository.UserRepository = (*UserRepo)(nil)

// UserRepo is the repository for the user model
type UserRepo struct {
	queries generated.Querier
}

func NewUserRepo(db *Store) *UserRepo {
	return &UserRepo{
		queries: generated.New(db.pool),
	}
}

func (r *UserRepo) CreateUser(
	ctx context.Context,
	user *repository.User,
) (*repository.User, error) {
	createdUser, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Password:     *user.Password,
		IsAdmin:      user.IsAdmin,
		RefreshToken: *user.RefreshToken,
	})
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating user: %s", err.Error())
	}

	user.ID = uint32(createdUser.ID)
	user.CreatedAt = createdUser.CreatedAt
	user.Password = nil
	user.RefreshToken = nil

	return user, nil
}

func (r *UserRepo) GetUserInternal(
	ctx context.Context,
	id uint32,
	email string,
) (*repository.User, error) {
	params := generated.GetUserInternalParams{}

	if id != 0 {
		params.ID = pgtype.Int8{
			Int64: int64(id),
			Valid: true,
		}
	}

	if email != "" {
		params.Email = pgtype.Text{
			String: email,
			Valid:  true,
		}
	}

	if id == 0 && email == "" {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "id or email is required")
	}

	user, err := r.queries.GetUserInternal(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user: %s", err.Error())
	}

	return &repository.User{
		ID:           uint32(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Password:     &user.Password,
		IsAdmin:      user.IsAdmin,
		RefreshToken: &user.RefreshToken,
		CreatedAt:    user.CreatedAt,
	}, nil
}

func (r *UserRepo) GetUser(
	ctx context.Context,
	id uint32,
	email string,
) (*repository.User, error) {
	getUserParams := generated.GetUserParams{}
	if id != 0 {
		getUserParams.ID = pgtype.Int8{
			Int64: int64(id),
			Valid: true,
		}
	}

	if email != "" {
		getUserParams.Email = pgtype.Text{
			String: email,
			Valid:  true,
		}
	}

	if id == 0 && email == "" {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "id or email is required")
	}

	user, err := r.queries.GetUser(ctx, getUserParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user: %s", err.Error())
	}

	return &repository.User{
		ID:          uint32(user.ID),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsAdmin:     user.IsAdmin,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (r *UserRepo) ListUsers(
	ctx context.Context,
	filter *repository.UserFilter,
) ([]*repository.User, *pkg.Pagination, error) {
	paramListUser := generated.ListUsersParams{
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	}

	paramListUserCount := generated.ListUsersCountParams{}

	if filter.Search != nil {
		search := strings.ToLower(*filter.Search)
		paramListUser.Search = pgtype.Text{
			String: "%" + search + "%",
			Valid:  true,
		}
		paramListUserCount.Search = pgtype.Text{
			String: "%" + search + "%",
			Valid:  true,
		}
	}

	if filter.IsAdmin != nil {
		paramListUser.IsAdmin = pgtype.Bool{
			Bool:  *filter.IsAdmin,
			Valid: true,
		}
		paramListUserCount.IsAdmin = pgtype.Bool{
			Bool:  *filter.IsAdmin,
			Valid: true,
		}
	}

	users, err := r.queries.ListUsers(ctx, paramListUser)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error listing users: %s", err.Error())
	}

	totalUsers, err := r.queries.ListUsersCount(ctx, paramListUserCount)
	if err != nil {
		return nil, nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error counting users: %s", err.Error())
	}

	usersResult := make([]*repository.User, len(users))
	for idx, user := range users {
		usersResult[idx] = &repository.User{
			ID:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			IsAdmin:     user.IsAdmin,
			CreatedAt:   user.CreatedAt,
		}
	}

	return usersResult, pkg.CalculatePagination(
		uint32(totalUsers),
		filter.Pagination.PageSize,
		filter.Pagination.Page,
	), nil
}

func (r *UserRepo) UpdateUser(
	ctx context.Context,
	user *repository.UpdateUser,
) (*repository.User, error) {
	params := generated.UpdateUserParams{
		ID: int64(user.ID),
	}

	if user.Name != nil {
		params.Name = pgtype.Text{
			String: *user.Name,
			Valid:  true,
		}
	}

	if user.Email != nil {
		params.Email = pgtype.Text{
			String: *user.Email,
			Valid:  true,
		}
	}

	if user.PhoneNumber != nil {
		params.PhoneNumber = pgtype.Text{
			String: *user.PhoneNumber,
			Valid:  true,
		}
	}

	if user.Password != nil {
		params.Password = pgtype.Text{
			String: *user.Password,
			Valid:  true,
		}
	}

	if user.IsAdmin != nil {
		params.IsAdmin = pgtype.Bool{
			Bool:  *user.IsAdmin,
			Valid: true,
		}
	}

	if user.RefreshToken != nil {
		params.RefreshToken = pgtype.Text{
			String: *user.RefreshToken,
			Valid:  true,
		}
	}

	updatedUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error updating user: %s", err.Error())
	}

	return &repository.User{
		ID:          uint32(updatedUser.ID),
		Name:        updatedUser.Name,
		Email:       updatedUser.Email,
		PhoneNumber: updatedUser.PhoneNumber,
		IsAdmin:     updatedUser.IsAdmin,
		CreatedAt:   updatedUser.CreatedAt,
	}, nil
}
