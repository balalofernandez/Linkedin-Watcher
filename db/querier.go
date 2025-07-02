package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	// User methods
	GetUserByID(ctx context.Context, id pgtype.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByGoogleID(ctx context.Context, googleID pgtype.Text) (User, error)
	CreateUserWithPassword(ctx context.Context, arg CreateUserWithPasswordParams) (User, error)
	CreateUserWithGoogle(ctx context.Context, arg CreateUserWithGoogleParams) (User, error)
	CreateUserWithBoth(ctx context.Context, arg CreateUserWithBothParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	UpdateUserGoogleAuth(ctx context.Context, arg UpdateUserGoogleAuthParams) error
	UpdateUserTokens(ctx context.Context, arg UpdateUserTokensParams) error
	UpdateUserActiveStatus(ctx context.Context, arg UpdateUserActiveStatusParams) error
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	ListUsers(ctx context.Context) ([]ListUsersRow, error)

	// LinkedIn Profile methods
	GetLinkedInProfileByID(ctx context.Context, id pgtype.UUID) (LinkedinProfile, error)
	GetLinkedInProfileByURL(ctx context.Context, linkedinUrl string) (LinkedinProfile, error)
	CreateLinkedInProfile(ctx context.Context, arg CreateLinkedInProfileParams) (LinkedinProfile, error)
	UpdateLinkedInProfile(ctx context.Context, arg UpdateLinkedInProfileParams) error
	ListLinkedInProfiles(ctx context.Context) ([]LinkedinProfile, error)

	// Company methods
	GetCompanyByID(ctx context.Context, id pgtype.UUID) (Company, error)
	GetCompanyByName(ctx context.Context, name string) (Company, error)
	GetCompanyByLinkedInURL(ctx context.Context, linkedinUrl pgtype.Text) (Company, error)
	CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error)
	UpdateCompany(ctx context.Context, arg UpdateCompanyParams) error
	ListCompanies(ctx context.Context) ([]Company, error)

	// Utility methods
	PingDb(ctx context.Context) (int32, error)
}
