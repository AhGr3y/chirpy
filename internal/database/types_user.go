package database

import (
	"errors"
	"os"

	"github.com/ahgr3y/chirpy/internal/auth"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

// CreateUser creates a User and saves it in the database
func (db *DB) CreateUser(email string, password string) (User, error) {

	// Load database
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Ensure no duplicate email
	if isDuplicate := hasDuplicateEmail(dbStructure, email); isDuplicate {
		return User{}, errors.New("cannot create user: duplicate email")
	}

	// Get unique id of new User
	id := len(dbStructure.Users) + 1

	// Create a new User
	user := User{
		ID:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}

	// Save user to database
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// hasDuplicateEmail checks the database for duplicate email
func hasDuplicateEmail(dbStructure DBStructure, email string) bool {

	users := dbStructure.Users

	for _, user := range users {
		if user.Email == email {
			return true
		}
	}

	return false
}

// GetUser retrieves a single user by id
func (db *DB) GetUser(id int) (User, error) {

	db.mux.RLock()
	defer db.mux.RUnlock()

	// Retrieve dbStructure from database
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	users := dbStructure.Users
	user, exist := users[id]
	if !exist {
		return User{}, os.ErrNotExist
	}

	return user, nil
}

// AuthenticateUser compares given password and saved password
// and return the User upon successful authentication
func (db *DB) AuthenticateUser(email string, password string) (User, error) {

	// Retrieve dbStructure from database
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Loop through users to get user that matches email
	users := dbStructure.Users
	for _, user := range users {
		if user.Email == email {
			// Check if password matches
			err := auth.AuthenticatePassword(user.Password, password)
			if err != nil {
				return User{}, err
			}

			return user, nil
		}
	}

	return User{}, os.ErrNotExist
}

// UpdateUser updates user's email and/or password
func (db *DB) UpdateUserEmailPassword(id int, email string, password string, isChirpyRed bool) (User, error) {

	// Updated user
	user := User{
		ID:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: isChirpyRed,
	}

	// Retrieve database
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Upload user to database
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	// Return updated user
	return user, nil
}

func (db *DB) UpdateUserToDatabase(user User) error {

	// Retrieve database
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Update user to database
	dbStructure.Users[user.ID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// UpgradeUser promotes user with userID to a Chirpy Red user.
func (db *DB) UpgradeUser(userID int) error {

	// Load database.
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Retrieve user from database.
	user, err := db.GetUser(userID)
	if err != nil {
		return err
	}

	// Upgrade user to Chirpy Red status
	user.IsChirpyRed = true

	// Update database.
	dbStructure.Users[user.ID] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
