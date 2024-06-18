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

	// Create a Chirp
	chirpOne, err := db.CreateChirp("This is the first chirp.")
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
	if chirps[0].ID != chirpOne.ID {
		t.Errorf("%v != %v: Expecting both ID to be 1.", chirps[0].ID, chirpOne.ID)
	}
	if chirps[0].Body != chirpOne.Body {
		t.Errorf("%v != %v: Expecting both body to be 'This is the first chirp.'", chirps[0].Body, chirpOne.Body)
	}

}
