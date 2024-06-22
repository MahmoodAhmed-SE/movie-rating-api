package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	var data RequestBody
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	// saving user data into postgresdb
	hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		switch hashErr.(type) {
		case bcrypt.InvalidCostError:
			http.Error(w, "Invalid cost parameter", http.StatusBadRequest)
		case bcrypt.InvalidHashPrefixError:
			http.Error(w, "Invalid hash prefix", http.StatusBadRequest)
		case bcrypt.HashVersionTooNewError:
			http.Error(w, "Hash version too new", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Printf("Error generating hash: %v", hashErr)
		return
	}
	fmt.Println(string(hashedPasswordBytes))
	dbConfigeration := pgx.ConnConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     "localhost",
		Database: "recordings",
		Port:     5432,
	}
	conn, connErr := pgx.Connect(dbConfigeration)

	if connErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var albums []Album

	albumsQuery, albumErrQuery := conn.Query("select * from album;")

	if albumErrQuery != nil {
		fmt.Printf("Error selecting album rows: %s", albumErrQuery.Error())
	}

	defer albumsQuery.Close()

	for albumsQuery.Next() {
		var alb Album
		if err := albumsQuery.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			fmt.Println(err.Error())
			return
		}
		albums = append(albums, alb)
		fmt.Println(alb.Title)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/api/register-user", handleUserRegistration)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
