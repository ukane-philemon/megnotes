package mongodb

import (
	"fmt"
	"time"

	"github.com/ukane-philemon/megtask/db"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CreateAccount creates a new user with the provided username and password.
// An ErrorInvalidRequest will be returned is the username already exists.
func (mdb *MongoDB) CreateAccount(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("%w: missing username or password", db.ErrorInvalidRequest)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword error: %w", err)
	}

	userInfo := &dbUser{
		Username:  username,
		Password:  string(passwordHash),
		CreatedAt: time.Now().Unix(),
	}

	_, err = mdb.usersCollection.InsertOne(mdb.ctx, userInfo)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: please try another username", db.ErrorInvalidRequest)
		}
	}

	return nil
}

// Login checks that the provided username and password matches a record in
// the database and are correct. Returns ErrorInvalidRequest if the password
// or username does not match any record.
func (mdb *MongoDB) Login(username, password string) (*db.User, error) {
	return nil, nil
}
