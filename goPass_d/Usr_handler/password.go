package pwuser

import (
	"github.com/ayushsherpa111/goPassd/schema"
	util "github.com/ayushsherpa111/goPassd/utils"
	"github.com/fatih/structs"
)

type PasswordGenerator struct{}

func (p PasswordGenerator) Generate(i schema.Item) schema.Item {
	var newPass = new(Password)
	*newPass = *i.(*Password)
	return newPass
}

type Password struct {
	Pid      int
	Username string
	Email    string
	Password []byte
	Nonce    []byte
	Site     string
	tblName  string
	Path     string
}

func RemoveKey(str string) bool {
	return str == "Path" || str == "Key" || str == "ReEnter" || str == "Pid"
}

func (p *Password) InsertKeys() []string {
	customFunc := func(s string) bool {
		return util.IsCapitalized(s) && !RemoveKey(s)
	}
	return util.Filter(structs.Names(p), customFunc)
}

func (p *Password) GetTbl() string {
	return p.tblName
}

func (p *Password) SetTbl(tbl string) {
	p.tblName = tbl
}

func (p *Password) InsertVals() []interface{} {
	var a = make([]interface{}, 0)
	a = append(a, p.Username, p.Email, p.Password, p.Nonce, p.Site)
	return a
}

func (p *Password) GetId() string {
	return p.Path
}
