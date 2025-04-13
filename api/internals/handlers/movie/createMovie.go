package movie

import (
	// "net/http"
	// "encoding/json"
	// "log"
	// "movie-rating-api-go/internals/services"
	"database/sql"
)

type MovieCreationReqBody struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Images         []sql.NullString `json:"images"`
	Release_date   sql.NullTime     `json:"release_date"`
	Director       sql.NullString   `json:"director"`
	Rating_average sql.NullFloat64  `json:"rating_average"`
	Duration       sql.NullInt32    `json:"duration"`
}


// func CreateMovie(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
	
// 	decoder := json.NewDecoder(r.Body)

// 	var movie MovieCreationReqBody 
// 	if err := decoder.Decode(&movie); err != nil {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		log.Printf("Error decoding body: %v", err)
// 		return
// 	}

// 	if err := services.AddMovie(&movie); err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		log.Printf("Error inserting a created movie to database: %v", err)
// 		return
// 	}
// }
