package main

import (
	"errors"
	"log"
	"net/http"
	"net/rpc"

	db "github.com/ayushsherpa111/goPassd/DB_handler"
	pwuser "github.com/ayushsherpa111/goPassd/Usr_handler"
	"github.com/ayushsherpa111/goPassd/schema"
	util "github.com/ayushsherpa111/goPassd/utils"
)

type API int

var (
	Api_db = new(db.DatabaseFactory)
)

func (a *API) RegisterUser(newUser pwuser.User, reply *pwuser.User) error {
	// Generate a Salt for the user
	(&newUser).GenSalt()
	(&newUser).GenHash()
	newUser.SetTbl("user")
	Api_db.CreateConfig(db.SQLiteService{}, newUser.HomePath) // Create the .goPass folder and the db file
	// set the password string for the user BTS so the user doesnt have to keep entering the password
	Api_db.SetPassMap(newUser.HomePath, newUser.Key, newUser.Salt)
	Api_db.MakeDB(newUser.HomePath, schema.User_DB_SCHEMA)     // make the user table
	Api_db.MakeDB(newUser.HomePath, schema.Password_DB_Schema) // make the password table
	conn := Api_db.GetDB(newUser.GetId())
	if ret, err := Api_db.Insert(conn, &newUser); err != nil {
		log.Fatalln(err)
		return err
	} else {
		reply = ret.(*pwuser.User)
		return nil
	}
}

// Path is the PK for the users table
func (a *API) GetInfo(path string, reply *pwuser.User) error {
	log.Println(path)

	conn := Api_db.GetDB(path)

	result := conn.GetOne("user", "Username, Email, HomePath", "HomePath", path)
	log.Println(result)
	log.Println(result.Err())
	if result.Err() != nil {
		return result.Err()
	}
	var user pwuser.User
	result.Scan(&user.Username, &user.Email, &user.HomePath)
	*reply = user
	log.Println(user)
	return nil
}

type ReEnter struct {
	Path     string
	Password []byte
}

func (a *API) ReEnterPassword(p ReEnter, reply *bool) error {
	log.Println("password reentered", string(p.Password))
	if usr, e := Api_db.CheckPasswordMatch(p.Password, p.Path); e {
		Api_db.SetPassMap(p.Path, p.Password, usr.Salt)
		*reply = true
		return nil
	}
	return errors.New("Passwords dont match. Try Again")
}

func (a *API) AddPassword(p pwuser.Password, reply *bool) error {
	encr, err := Api_db.GetPassMap(p.Path)
	if err != nil {
		return err
	}
	// safe to start adding password
	cpr, nonce, err := encr.Encrypt(p.Password)
	if err != nil {
		return err
	}
	p.SetTbl("passwords")
	p.Nonce = nonce
	p.Password = cpr
	conn := Api_db.GetDB(p.GetId())
	_, err = Api_db.Insert(conn, &p)
	return err
}

func (a *API) GetPass(pass schema.MatchPass, reply *[]pwuser.Password) error {
	encr, err := Api_db.GetPassMap(pass.Key)
	if err != nil {
		log.Println("ReEnter PASS")
		return err
	}

	var pwd pwuser.Password
	var ValuePointer = []interface{}{&pwd.Pid, &pwd.Username, &pwd.Email, &pwd.Password, &pwd.Nonce, &pwd.Site}
	var e error
	*reply, e = Api_db.GetPasswords(pass, "pid,Username, Email,Password, Nonce,Site", ValuePointer, &pwd, encr)
	if e != nil {
		log.Println(e.Error())
	}
	return e
}

func (a *API) DeletePassword(criteria schema.MatchPass, reply *int64) error {
	*reply, _ = Api_db.DeleteItem(criteria, Api_db.GetDB(criteria.Key))
	return nil
}

func (a *API) AddCSV(csvPayload schema.CSVData, reply *bool) error {
	log.Println(csvPayload)
	if e := util.CheckIfExists(csvPayload.CsvPath); e != nil {
		return e
	}
	if channel, e := Api_db.InsertCSV(csvPayload); e != nil {
		return e
	} else {
		conn := Api_db.GetDB(csvPayload.Key)
		for v := range channel {
			Api_db.Insert(conn, v)
		}
	}
	// csv file exists
	*reply = true
	return nil
}

type Test struct {
	Key  string
	Pass []byte
}

func startUp() {
	keys := util.GetAllHomePaths()
	for u, v := range keys {
		log.Println("Setting up ", u)
		Api_db.SetUpDB(v, db.SQLiteService{})
	}
}

func main() {
	api := new(API)
	log.SetFlags(log.Lshortfile)
	Api_db.SetUpMap()
	go startUp()
	rpc.Register(api)

	rpc.HandleHTTP()
	if e := http.ListenAndServe(":5555", nil); e != nil {
		log.Fatalln(e.Error())
	}
}
