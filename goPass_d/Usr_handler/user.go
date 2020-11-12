package pwuser

import (
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/ayushsherpa111/goPassd/schema"
	util "github.com/ayushsherpa111/goPassd/utils"
	"github.com/fatih/structs"
)

type UserGenerator struct{}

func (u UserGenerator) Generate(i schema.Item) schema.Item {
	var newUser = new(User)
	*newUser = *i.(*User)
	return newUser
}

type User struct {
	Salt     []byte
	Username string // Sent from client
	HomePath string // Sent from client
	Email    string // Sent from client
	Hash     []byte
	Key      []byte // Sent from client, [MUST BE REMOVED FROM THE OBJECT BEFORE SAVING TO DB]
	tblName  string // should be removed. use plain text as argument when inserting
}

func (u *User) GetId() string {
	return u.HomePath
}

func (u *User) SetTbl(tbl string) {
	u.tblName = tbl
}

func (u *User) GetTbl() string {
	return u.tblName
}

func (u *User) InsertKeys() []string {
	customFunc := func(k string) bool {
		return util.IsCapitalized(k) && !RemoveKey(k)
	}
	return util.Filter(structs.Names(u), customFunc)
}

func (u *User) InsertVals() []interface{} {
	vals := make([]interface{}, 0)
	vals = append(vals, u.Salt, u.Username, u.HomePath, u.Email, u.Hash)
	log.Println("INSERTING ", vals)
	return vals
}

func (u *User) GenSalt() {
	Salt := make([]byte, 12)
	rand.Read(Salt)
	u.Salt = Salt
}

func (u *User) GenHash() {
	hasher := sha256.New()
	u.Key = append(u.Key, u.Salt...)
	hasher.Write(u.Key)
	u.Hash = hasher.Sum(nil)
}
