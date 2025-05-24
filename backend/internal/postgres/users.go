package postgres

import (
	"context"

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
		Password:     user.Password,
		IsAdmin:      user.IsAdmin,
		RefreshToken: user.RefreshToken,
	})
	if err != nil {
		if pkg.PgxErrorCode(err) == pkg.UNIQUE_VIOLATION {
			return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "%s", err.Error())
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error creating user: %s", err.Error())
	}

	user.ID = createdUser.ID

	return user, nil
}

func (r *UserRepo) GetUserByID(
	ctx context.Context,
	id int64,
) (*repository.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user: %s", err.Error())
	}

	return marshalUser(user), nil
}

func (r *UserRepo) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user: %s", err.Error())
	}

	return marshalUser(user), nil
}

func (r *UserRepo) ListUsers(
	ctx context.Context,
	filter *repository.UserFilter,
) ([]*repository.User, *pkg.Pagination, error) {
	paramListUser := generated.ListUsersParams{
		Search: pgtype.Text{
			Valid: false,
		},
		IsAdmin: pgtype.Bool{
			Valid: false,
		},
		Limit:  int32(filter.Pagination.PageSize),
		Offset: pkg.Offset(filter.Pagination.Page, filter.Pagination.PageSize),
	}

	paramListUserCount := generated.ListUsersCountParams{
		Search: pgtype.Text{
			Valid: false,
		},
		IsAdmin: pgtype.Bool{
			Valid: false,
		},
	}

	if filter.Search != nil {
		paramListUser.Search = pgtype.Text{
			String: *filter.Search,
			Valid:  true,
		}
		paramListUserCount.Search = pgtype.Text{
			String: *filter.Search,
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
			ID:          user.ID,
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
		ID: user.ID,
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

	return marshalUser(updatedUser), nil
}

func marshalUser(user generated.User) *repository.User {
	return &repository.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Password:     user.Password,
		IsAdmin:      user.IsAdmin,
		RefreshToken: user.RefreshToken,
		CreatedAt:    user.CreatedAt,
	}
}
