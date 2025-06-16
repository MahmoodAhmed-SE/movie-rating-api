package services

import (
	"errors"
	"fmt"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateWSToken(inputToken string) (int, error) {
	_, scanErr := fmt.Sscanf(inputToken, "Bearer %s", &inputToken)
	if scanErr != nil {
		log.Printf("Error scanning header authorization token: %v", scanErr)
		return 0, scanErr
	}

	token, err := jwt.Parse(inputToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv(constants.EnvJWTSecretKey)), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return 0, err
	}

	if !token.Valid {
		log.Println("invalid token")
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("unauthorized")
	}

	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		log.Printf("expired token: %v", err)
		return 0, fmt.Errorf("token expired: %v", token)
	}

	userId, ok := claims["id"].(float64)

	if !ok {
		log.Printf("id claim assertion failed: %v", userId)
		return 0, err
	}

	return int(userId), nil
}

func GetUser(username string) (models.USER, error) {
	var user models.USER

	conn := database.GetConn()
	row := conn.QueryRow(`
		SELECT * FROM USERS 
		WHERE username = $1;`,
		username)

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return user, err
	}

	return user, nil
}
