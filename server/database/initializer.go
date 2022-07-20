package database

import (
	_ "github.com/lib/pq"
	"os"
)

func createInitialUser() error {
	firstUserName := os.Getenv("DB_FIRST_USER_NAME")
	firstPassword, err := hashPassword(os.Getenv("DB_FIRST_USER_PASSWORD"))
	if err != nil {
		return err
	}

	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO \"users\" (name, password) VALUES ($1, $2)", firstUserName, firstPassword)

	return err
}

func InitializeDatabase() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE \"users\" (id uuid primary key default uuid_generate_v4(), name text unique not null, password text not null)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE \"api_key\" (id uuid primary key default uuid_generate_v4(), token text not null unique, user_id uuid references \"users\"(id))")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE \"backup_job\" (id uuid primary key default uuid_generate_v4(), name text not null, host text not null, port int not null default 21, type text not null default 'ftp', username text not null default '', password text not null default '', remote_path text not null, local_path text not null)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE \"stored_backup\" (id uuid primary key default uuid_generate_v4(), full_path text not null, name text not null, backup_date timestamp not null default now(), backup_job_id uuid references \"backup_job\"(id))")
	if err != nil {
		return err
	}

	return createInitialUser()
}
