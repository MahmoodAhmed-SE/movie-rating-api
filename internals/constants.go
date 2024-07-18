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
