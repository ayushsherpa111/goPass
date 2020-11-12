package db

import (
	"database/sql"
)

type DB_Factory interface {
	Create(string) DbService
}

type SQLiteService struct{}

func (fac SQLiteService) Create(connection string) DbService {
	conn, _ := sql.Open("sqlite3", connection)
	return SQLService{connection: conn}
}
