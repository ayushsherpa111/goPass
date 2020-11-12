package db

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"os"
	"path"

	encr "github.com/ayushsherpa111/goPassd/Encr"
	reader "github.com/ayushsherpa111/goPassd/Reader"
	pwuser "github.com/ayushsherpa111/goPassd/Usr_handler"
	"github.com/ayushsherpa111/goPassd/schema"
	util "github.com/ayushsherpa111/goPassd/utils"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseFactory struct {
	collections map[string]DbService
	passMap     map[string]encr.Encpr
}

func (db *DatabaseFactory) SetPassMap(key string, val []byte, salt []byte) {
	gcm := &encr.GCM_Encr{}
	gcm.Init(val, salt)
	db.passMap[key] = gcm
}

func (db *DatabaseFactory) GetPassMap(key string) (encr.Encpr, error) {
	if v, ok := db.passMap[key]; ok {
		return v, nil
	}
	return nil, errors.New("The user does not have their password registered")
}

func (db *DatabaseFactory) SetUpDB(userPath string, factory DB_Factory) {
	connectionString := path.Join(userPath, util.ConfigFolder, util.MAIN_DB)
	db.collections[userPath] = factory.Create(connectionString)
}

func (db *DatabaseFactory) GetPasswords(matchPass schema.MatchPass,
	fields string,
	valuePtrs []interface{},
	value *pwuser.Password,
	cryp encr.Encpr,
) ([]pwuser.Password, error) {
	// used for getting passwords, all and filtered
	rows, e := db.GetItem(matchPass, fields) // get me a pointer to sql rows
	results := make([]pwuser.Password, 0)
	if e != nil {
		return results, e
	}
	for rows.Next() {
		rows.Scan(valuePtrs...)
		db.DecryptData(cryp, value)
		results = append(results, *value)
	}
	return results, e
}

func (db *DatabaseFactory) GetDB(key string) DbService {
	con, ok := db.collections[key]
	if !ok {
		db.SetUpDB(key, SQLiteService{})
		con = db.collections[key]
	}
	return con
}

func (db *DatabaseFactory) GetItem(item schema.MatchPass,
	fields string,
) (*sql.Rows, error) {
	con := db.GetDB(item.Key)
	log.Println("Got DB", con)
	return con.GetItemLike(item.Tbl, fields, item.KeyVal)
}

func (db *DatabaseFactory) DecryptData(dcr encr.Encpr, val *pwuser.Password) {
	plain, err := dcr.Decrypt(val.Password, val.Nonce)
	if err != nil {
		val.Password = []byte("Failed to Decrypt")
		return
	}
	val.Password = plain
}

func (db *DatabaseFactory) CreateConfig(factory DB_Factory, userPath string) {
	// Creates the .goPass folder
	util.SetUpConfigDir(userPath)
	// creates the config file for the DB
	util.SetUpDb(userPath)
	// create a connection
	db.SetUpDB(userPath, factory)
}

func (db *DatabaseFactory) CheckPasswordMatch(newP []byte, key string) (pwuser.User, bool) {
	log.Println(string(newP), key)
	conn := db.GetDB(key)
	res := conn.GetOne("user", "Hash, Salt", "HomePath", key)
	var targetVal pwuser.User
	res.Scan(&targetVal.Hash, &targetVal.Salt)
	newUser := new(pwuser.User)
	newUser.Salt = targetVal.Salt
	newUser.Key = newP
	newUser.GenHash()
	if bytes.Compare(newUser.Hash, targetVal.Hash) == 0 {
		return targetVal, true
	}
	return targetVal, false
}

func (db *DatabaseFactory) MakeDB(key string, schema string) {
	// Call Create DB on the map
	conn := db.GetDB(key)
	conn.CreateDb(schema)
}

func (db *DatabaseFactory) InsertCSV(csvPayload schema.CSVData) (chan schema.Item, error) {
	var err error
	file, _ := os.Open(csvPayload.CsvPath)
	encptr, err := db.GetPassMap(csvPayload.Key)
	if err != nil {
		return nil, err
	}
	return reader.ReadCSVFile(file, db.prepPassword(encptr, csvPayload.Key)), err
}

func (db *DatabaseFactory) Insert(conn DbService, i schema.Item) (schema.Item, error) {
	query := util.GenerateInsertQuery(i.GetTbl(), i.InsertKeys()) // Provide the DB with a string of prepared statement to execute
	return i, conn.Insert(query, i)
}

func (db *DatabaseFactory) SetUpMap() {
	db.collections = make(map[string]DbService)
	db.passMap = make(map[string]encr.Encpr)
}

func (db *DatabaseFactory) DeleteItem(crit schema.MatchPass, conn DbService) (int64, error) {
	res, e := conn.DeleteItem(crit.Tbl, crit.KeyVal)
	if e != nil {
		log.Println(e.Error())
	}
	return res.RowsAffected()
}

func (db *DatabaseFactory) prepPassword(encr encr.Encpr, key string) func([]string, []string) schema.Item {
	return func(header []string, data []string) schema.Item {
		mapData, _ := util.Zip(header, data)
		newPass := &pwuser.Password{}
		cpr, nonce, _ := encr.Encrypt([]byte(mapData["password"]))
		newPass.Username = mapData["username"]
		newPass.Email = mapData["username"]
		newPass.Password = cpr
		newPass.Nonce = nonce
		newPass.Site = mapData["name"]
		newPass.SetTbl("passwords")
		newPass.Path = key
		return newPass
	}
}
