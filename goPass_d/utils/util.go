package util

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"

	// pwuser "github.com/ayushsherpa111/goPassd/Usr_handler"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

const (
	ConfigFolder = ".goPass"
	MAIN_DB      = "store.db"
	ITER         = 100000
	KEY_LEN      = 32 // bytes (AES-256)
)

func CreateIfNotExist(filePath string) error {
	if _, e := os.Stat(filePath); os.IsNotExist(e) {
		os.Create(filePath)
	}
	return nil
}

func CheckIfExists(filePath string) error {
	_, e := os.Stat(filePath)
	return e
}

func Filter(arr []string, check func(string) bool) []string {
	var newAR []string
	for _, v := range arr {
		if check(v) {
			newAR = append(newAR, v)
		}
	}
	return newAR
}

func IsCapitalized(val string) bool {
	return val[0] >= 65 && val[0] <= 90
}

func SetUpConfigDir(userHomePath string) {
	// creates .goPass directory
	configPath := path.Join(userHomePath, ConfigFolder)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err = os.Mkdir(configPath, 0777); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func GenerateInsertQuery(tbl_name string, fields []string) string {
	// log.Println(tbl_name)
	// log.Println(fields)
	insertStmt := fmt.Sprintf("INSERT INTO %s", tbl_name)
	insertStmt += "(" + strings.Join(fields, ",") + ") VALUES (" + strings.Repeat("?,", len(fields))[:2*len(fields)-1] + ")"
	log.Println(insertStmt)
	return insertStmt
}

func GenUid() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

func SetUpDb(homepath string) {
	// creates the db file for the user
	if e := CreateIfNotExist(path.Join(homepath, ConfigFolder, MAIN_DB)); e != nil {
		log.Fatal(e.Error())
	}
}

func GenerateSelectQuery(tbl_name string, fields string) string {
	return fmt.Sprintf("SELECT %s FROM %s", fields, tbl_name)
}

func GenerateSelectLikeQuery(tbl_name string, fields string, colms map[string]interface{}) string {
	baseQuery := GenerateSelectQuery(tbl_name, fields)
	count := 0
	for key, val := range colms {
		if count == 0 {
			baseQuery = fmt.Sprintf("%s WHERE", baseQuery)
		} else {
			baseQuery += " OR "
		}
		count++
		if key == "pid" {
			baseQuery += fmt.Sprintf(" %s = %v", key, val)
			continue
		}
		baseQuery += fmt.Sprintf(" %s LIKE '%%%v%%'", key, val)
	}
	return baseQuery
}

func GetAllHomePaths() map[string]string {
	var startID = 1000
	var list = make(map[string]string)
	for ; ; startID += 1 {
		if u, e := user.LookupId(fmt.Sprint(startID)); e == nil {
			if exist, _ := regexp.Match("home", []byte(u.HomeDir)); exist {
				list[u.Name] = u.HomeDir
			}
		} else {
			break
		}
	}
	return list
}

func GenPass(pass []byte, salt []byte) []byte {
	return pbkdf2.Key(pass, salt, ITER, KEY_LEN, sha256.New)
}

func Zip(header []string, data []string) (map[string]string, error) {
	result := make(map[string]string)
	if len(header) != len(data) {
		return nil, errors.New("Arugment Len mismatch")
	}
	for i, v := range header {
		result[v] = data[i]
	}
	return result, nil
}

func GenerateDelete(tbl_name string, fieldsData map[string]interface{}) string {
	query := fmt.Sprintf("DELETE FROM %s WHERE ", tbl_name)
	count := 0
	for i, v := range fieldsData {
		query += fmt.Sprintf("%s = %v ", i, v)
		if count < len(fieldsData)-1 {
			query += " AND "
		}
		count++
	}
	return query
}
