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
