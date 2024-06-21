package database

import (
	"sync"
	"testing"
)

func TestCreateUser(t *testing.T) {

	const databaseFilepath = "../../database.json"

	db := DB{
		path: databaseFilepath,
		mux:  &sync.RWMutex{},
	}

	user, err := db.CreateUser("luffy@onepiece.com", "ilovemeat") // id = 1
	if err != nil {
		t.Error("Failed to create user")
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		t.Error("Failed to load database")
	}

	if dbStructure.Users[1] != user {
		t.Errorf("%v != %v", dbStructure.Users[1], user)
	}

	_, err = db.CreateUser("luffy@onepiece.com", "ilovemeat")
	if err == nil {
		t.Error("Cannot create duplicate users")
	}

}

func TestGetUser(t *testing.T) {

	const databaseFilepath = "../../database.json"

	db := DB{
		path: databaseFilepath,
		mux:  &sync.RWMutex{},
	}

	user, err := db.CreateUser("lane@bootdev.com", "password") // id = 2
	if err != nil {
		t.Error("Failed to create user")
	}

	userFromDB, err := db.GetUser(2)
	if err != nil {
		t.Error("Failed to get user")
	}

	if userFromDB != user {
		t.Errorf("%v != %v", userFromDB, user)
	}
}

func TestUpdateUser(t *testing.T) {

	const databaseFilepath = "../../database.json"

	db := DB{
		path: databaseFilepath,
		mux:  &sync.RWMutex{},
	}

	_, err := db.CreateUser("harry@wizards.com", "ilovevoldemort") // id = 3
	if err != nil {
		t.Error("Failed to create user")
	}

	user, err := db.UpdateUser(3, "ron@wizards.com", "iloveclowns")
	if err != nil {
		t.Error("Failed to update user")
	}

	dbUser, err := db.GetUser(3)
	if err != nil {
		t.Error("Failed to get user")
	}

	if dbUser != user {
		t.Errorf("%v != %v", user, dbUser)
	}

}
