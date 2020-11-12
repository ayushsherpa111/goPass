package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ayushsherpa111/goPassd/schema"
	util "github.com/ayushsherpa111/goPassd/utils"
)

type DbService interface {
	Insert(string, schema.Item) error
	CreateDb(string) error
	GetItem(string, string) (*sql.Rows, error)
	GetItemLike(string, string, map[string]interface{}) (*sql.Rows, error)
	GetOne(string, string, string, string) *sql.Row
	DeleteItem(string, map[string]interface{}) (sql.Result, error)
}

type SQLService struct {
	connection *sql.DB
}

func (s SQLService) Insert(query string, i schema.Item) error {
	// Accept a schema to generate
	log.Println("Generated Query: ", query)
	if stmt, err := s.connection.Prepare(query); err != nil {
		return err
	} else {
		log.Println(i.InsertVals())
		_, err = stmt.Exec(i.InsertVals()...)
		return err
	}
}

func (s SQLService) CreateDb(tbl_schema string) error {
	log.Println("Creating DB", tbl_schema)
	if stmt, err := s.connection.Prepare(tbl_schema); err != nil {
		return err
	} else {
		_, err = stmt.Exec()
		return err
	}
}

func (s SQLService) GetItem(tbl_name string, fields string) (*sql.Rows, error) {
	query := util.GenerateSelectQuery(tbl_name, fields)
	return s.connection.Query(query)
}

func (s SQLService) GetItemLike(tbl_name string, fields string, keyVal map[string]interface{}) (*sql.Rows, error) {
	query := util.GenerateSelectLikeQuery(tbl_name, fields, keyVal)
	log.Println(query)
	return s.connection.Query(query)
}

func (s SQLService) GetOne(tbl_name string, fields string, findField string, eqVal string) *sql.Row {
	query := util.GenerateSelectQuery(tbl_name, fields)
	query += fmt.Sprintf(" WHERE %s = ?", findField)
	return s.connection.QueryRow(query, eqVal)
}

func (s SQLService) DeleteItem(tbl_name string, keyVal map[string]interface{}) (sql.Result, error) {
	query := util.GenerateDelete(tbl_name, keyVal)
	log.Println(query)
	return s.connection.Exec(query)
}
