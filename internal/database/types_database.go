package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist.
func NewDB(path string) (*DB, error) {

	// Check if path is empty
	if len(path) == 0 {
		return &DB{}, errors.New("path is empty")
	}

	// Create a new DB
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	// Create database file if it doesn't exist
	err := db.ensureDB()

	return db, err
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

// CreateUser creates a User and saves it in the database
func (db *DB) CreateUser(email string) (User, error) {

	// Load database
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Get unique id of new User
	id := len(dbStructure.Users) + 1

	// Create a new User
	user := User{
		ID:    id,
		Email: email,
	}

	// Save chirp to database
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
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

// createDB creates an empty dbStructure and writes it to disk
func (db *DB) createDB() error {

	// Create an empty dbStructure
	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}

	// Create a new database file
	return db.writeDB(dbStructure)
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {

	// Check if database.json exist
	_, err := os.ReadFile(db.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Create a new database file
			return db.createDB()
		}
	}

	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

	// Make database.json safe for reading
	db.mux.RLock()
	defer db.mux.RUnlock()

	// Read database.json
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return DBStructure{}, err
	}

	// Parse the JSON to DBStructure
	dbStructure := DBStructure{}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {

	// Make sure file is safe to read/write
	db.mux.Lock()
	defer db.mux.Unlock()

	// Parse dbStructure to JSON
	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	// Write dat to path
	err = os.WriteFile(db.path, dat, 0o600)
	if err != nil {
		return err
	}

	return nil
}
