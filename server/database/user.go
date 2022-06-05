package database

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

type User struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	password string `db:"password"`
}

func hashPassword(password string) (string, error) {
	sha := sha512.New()
	_, err := sha.Write([]byte(password))
	if err != nil {
		return "", err
	}

	hashed := sha.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashed), nil
}

func FindAllUsers() ([]User, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	users := make([]User, 0)

	err = db.Select(&users, "SELECT id, name, password FROM users")

	return users, err
}

func FindUserByName(username string) (*User, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	user := new(User)

	err = db.Get(user, "SELECT id, name, password FROM users WHERE name = $1", username)

	return user, err
}

func FindUserById(id string) (*User, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	user := new(User)

	err = db.Get(user, "SELECT id, name, password FROM users WHERE id = $1", id)

	return user, err
}

func Login(username, password string) (*ApiKey, *User, error) {
	hashedPw, err := hashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user, err := FindUserByName(username)
	if err != nil {
		return nil, nil, err
	}

	if user.password != hashedPw {
		return nil, nil, fmt.Errorf("hashes don't match")
	}

	apiKey := new(ApiKey)
	apiKey.UserId = user.Id
	err = apiKey.Create()
	if err != nil {
		return nil, nil, err
	}

	return apiKey, user, nil
}

func (user *User) SetPassword(password string) error {
	pw, err := hashPassword(password)
	if err != nil {
		return err
	}

	user.password = pw

	return nil
}

func (user *User) Create() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO users (name, password) VALUES ($1, $2)", user.Name, user.password)

	return err
}
func (user *User) Update() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE users SET name = $1, password = $2 WHERE id = $3", user.Name, user.password, user.Id)

	return err
}

func (user *User) Delete() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = $1", user.Id)

	return err
}

func (user *User) GetApiKeys() ([]ApiKey, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	apiKeys := make([]ApiKey, 0)
	err = db.Get(&apiKeys, "SELECT id, token, user_id FROM api_key WHERE user_id = $1", user.Id)

	return apiKeys, err
}
