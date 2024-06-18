package database

import "os"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {

	// Load database
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Get unique id of new Chirp
	id := len(dbStructure.Chirps) + 1

	// Create a new Chirp
	chirp := Chirp{
		ID:   id,
		Body: body,
	}

	// Save chirp to database
	dbStructure.Chirps[id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
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

// GetChirp retrieves a single Chirp by ID
func (db *DB) GetChirp(id int) (Chirp, error) {

	// Retrieve dbStructure from database
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Retrieve chirp from dbStructure
	chirp, exist := dbStructure.Chirps[id]
	if !exist {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}
