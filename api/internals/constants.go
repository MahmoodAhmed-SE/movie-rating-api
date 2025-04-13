package constants

type contextKey string

const (
	EnvDBUser string = "POSTGRES_USER"
	EnvDBPass string = "POSTGRES_PASSWORD"

	EnvJWTSecretKey string = "JWT_SECRET_TOKEN"
	EnvMovieAPI     string = "MOVIE_API_URL"

	/*
		Type definition to avoid probable context key collision.


		Context advice: Should not use built-in type string as key for value;
		define your own type to avoid collisions (SA1029)go-staticcheck
	*/
	UserIdKey contextKey = "user_id"
)

// errors
const (
	// PostgreSQL related
	UniqueConstraintViolation string = "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)"

	// Movie rating API Logic related
	ErrEmptyFilters            = "ERROR: no search filters are found for the operation of populating filters"
	ErrRatingOutOfBounds       = "ERROR: rating is out of bounds (either less than 0 or more than 10)"
	ErrIncorrectDurationFormat = "ERROR: incorrect duration format. duration format must be in the form of hh:mm exampe: 03:30"

	ErrSearchFilterTitleIsNotString     = "ERROR: unexpected type of filter title in filters map. it must be string but it isn't"
	ErrSearchFilterReleaseYearIsNotTime = "ERROR: unexpected type of filter release_year in filters map. it must be time.Time but it isn't"
	ErrSearchFilterRatingIsNotFloat64   = "ERROR: unexpected type of filter rating in filters map. it must be float64 but it isn't"
	ErrSearchFilterDirectorIsNotString
)
