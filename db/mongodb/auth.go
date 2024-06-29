package mongodb

import (
	"errors"
	"fmt"
	"time"

	"github.com/ukane-philemon/megtask/db"
	"go.mongodb.org/mongo-driver/bson"
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
		return fmt.Errorf("usersCollection.InsertOne error: %w", err)
	}

	return nil
}

// Login checks that the provided username and password matches a record in
// the database and are correct. Returns ErrorInvalidRequest if the password
// or username does not match any record.
func (mdb *MongoDB) Login(username, password string) (*db.User, error) {
	var dbUser *dbUser
	err := mdb.usersCollection.FindOne(mdb.ctx, bson.M{usernameKey: username}).Decode(&dbUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: username or password is incorrect", db.ErrorInvalidRequest)
		}
		return nil, fmt.Errorf("usersCollection.FindOne error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w: username or password is incorrect", db.ErrorInvalidRequest)
	}

	userID := dbUser.ID.Hex()
	tasks, err := mdb.userTasks(userID)
	if err != nil {
		mdb.log.Error("failed to retrieve user tasks: ", " error", err)
	}

	return &db.User{
		ID:       userID,
		Username: username,
		Tasks:    tasks,
	}, nil
}
