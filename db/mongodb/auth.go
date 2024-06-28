package mongodb

import "github.com/ukane-philemon/megtask/db"

// CreateAccount creates a new user with the provided username and password.
// An ErrorInvalidRequest will be returned is the username already exists.
func (mdb *MongoDB) CreateAccount(username, password string) error {
	return nil
}

// Login checks that the provided username and password matches a record in
// the database and are correct. Returns ErrorInvalidRequest if the password
// or username does not match any record.
func (mdb *MongoDB) Login(username, password string) (*db.User, error) {
	return nil, nil
}
