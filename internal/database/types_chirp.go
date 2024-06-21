package database

import (
	"errors"
	"os"
	"sort"
)

type Chirp struct {
	AuthorID int    `json:"author_id"`
	ID       int    `json:"id"`
	Body     string `json:"body"`
}

// CreateChirp creates a Chirp using body
// and saves it to the database.
func (db *DB) CreateChirp(userID int, body string) (Chirp, error) {

	// Load database.
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Generate unique id for Chirp.
	chirpID := len(dbStructure.Chirps) + 1

	// Initialize Chirp.
	chirp := Chirp{
		AuthorID: userID,
		ID:       chirpID,
		Body:     body,
	}

	// Save chirp to database.
	dbStructure.Chirps[chirpID] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database.
func (db *DB) GetChirps() ([]Chirp, error) {

	// Load DBStructure
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	// Empty slice to store Chirps
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	// Fill chirps with Chirps from dbStructure
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// GetChirps returns all chirps created by user with userID in the database.
func (db *DB) GetChirpsByID(userID int) ([]Chirp, error) {

	// Load DBStructure
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	// Empty slice to store Chirps
	chirps := []Chirp{}

	// Fill chirps with Chirps from dbStructure
	for _, chirp := range dbStructure.Chirps {
		if chirp.AuthorID == userID {
			chirps = append(chirps, chirp)
		}
	}

	return chirps, nil
}

// GetChirp retrieves a single Chirp by chirp ID.
func (db *DB) GetChirp(chirpID int) (Chirp, error) {

	// Retrieve dbStructure from database
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Retrieve chirp from dbStructure
	chirp, exist := dbStructure.Chirps[chirpID]
	if !exist {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}

// DeleteChirp deletes chirp with chirpID by user with userID.
// Sorts the remaining chirps to prevent ID conflict by
// future Chirps created.
func (db *DB) DeleteChirp(userID int, chirpID int) error {

	// Load database.
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Ensure Chirp can only be deleted by owner.
	chirpToDelete, err := db.GetChirp(chirpID)
	if err != nil {
		return err
	}
	if chirpToDelete.AuthorID != userID {
		return errors.New("unauthorized to delete chirp")
	}

	delete(dbStructure.Chirps, chirpID)

	// Re-assign chirp id for remaining chirps in database.
	db.ReassignChirpID(dbStructure.Chirps)

	// Update database.
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ReassignChirpID(chirps map[int]Chirp) {

	for i, chirp := range chirps {
		chirp.ID = i + 1
	}

}

// SortChirpByID sorts chirps in ascending order by ID.
func SortChirpsByID(chirps []Chirp, sortBy string) {

	// Sort in ascending order
	if sortBy == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}

	// Sort in descending order
	if sortBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

}
