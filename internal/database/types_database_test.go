package database

import "testing"

// Delete database.json before test
func TestWriteGetChirps(t *testing.T) {

	// Test NewDB and ensureDB
	db, err := NewDB("../../database.json")
	if err != nil {
		t.Error(err)
	}

	// Test GetChirps on empty db, and also test loadDB
	chirps, err := db.GetChirps()
	if err != nil {
		t.Error(err)
	}
	if len(chirps) != 0 {
		t.Errorf("Expecting an empty slice of Chirps")
	}

	// Load database.
	dbStructure, err := db.loadDB()
	if err != nil {
		t.Error(err)
	}

	// Generate unique id for Chirp.
	chirpID := len(dbStructure.Chirps) + 1

	// Initialize Chirp.
	chirp := Chirp{
		ID:   chirpID,
		Body: "This is the first chirp.",
	}

	// Save chirp to database.
	dbStructure.Chirps[chirpID] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		t.Error(err)
	}

	// Get Chirps from database
	chirps, err = db.GetChirps()
	if err != nil {
		t.Error(err)
	}

	// Test result of CreateChirp, writeChirp, GetChirps
	if len(chirps) != 1 {
		t.Errorf("Expecting length of 1.")
	}
	if chirps[0].ID != chirp.ID {
		t.Errorf("%v != %v: Expecting both ID to be 1.", chirps[0].ID, chirp.ID)
	}
	if chirps[0].Body != chirp.Body {
		t.Errorf("%v != %v: Expecting both body to be 'This is the first chirp.'", chirps[0].Body, chirp.Body)
	}

}
