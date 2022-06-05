package database

import "github.com/google/uuid"

type ApiKey struct {
	Token  string `db:"token"`
	Id     string `db:"id"`
	UserId string `db:"user_id"`
}

func FindApiKeyByToken(token string) (*ApiKey, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	apiKey := new(ApiKey)
	err = db.Get(apiKey, "SELECT id, token, user_id FROM api_key WHERE token = $1", token)

	return apiKey, err
}

func (apiKey *ApiKey) Create() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	apiKey.Token = uuid.New().String()
	_, err = db.Exec("INSERT INTO api_key (token, user_id) VALUES ($1, $2)", apiKey.Token, apiKey.UserId)

	return err
}

func (apiKey *ApiKey) Delete() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM api_key WHERE id = $1", apiKey.Id)

	return err
}

func (apiKey *ApiKey) GetUser() (*User, error) {
	return FindUserById(apiKey.UserId)
}
