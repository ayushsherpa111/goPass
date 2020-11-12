package rpcCon

import (
	"net/rpc"

	"github.com/ayushsherpa111/goPass/kit"
)

type Vault_RPC interface {
	Init_User(kit.User) (kit.User, error)
	Get_User_Info() (kit.User, error)
	Add_to_Vault()
	Get_from_Vault()
}

type Client_RPC struct {
	comms *rpc.Client
}

func (u *Client_RPC) Init_User(newUser kit.User) (kit.User, error) {
	var reply kit.User
	err := u.comms.Call("API.RegisterUser", newUser, &reply)
	return reply, err
}

func (u *Client_RPC) Get_User_Info(path string) (kit.User, error) {
	var reply kit.User
	err := u.comms.Call("API.GetInfo", path, &reply)
	return reply, err
}

func (u *Client_RPC) ReEnterPassword(p kit.ReEnter) error {
	var reply bool
	return u.comms.Call("API.ReEnterPassword", p, &reply)
}

func (u *Client_RPC) AddPassword(p kit.Password) error {
	var reply bool
	err := u.comms.Call("API.AddPassword", p, &reply)
	return err
}

func (u *Client_RPC) GetAll(key string, result *[]kit.Password) error {
	e := u.comms.Call("API.GetAll", key, &result)
	return e
}

func (u *Client_RPC) GetPass(matchData kit.MatchPass, result *[]kit.Password) error {
	return u.comms.Call("API.GetPass", matchData, result)
}

func (u *Client_RPC) GetOne(match kit.MatchPass, reply *[]kit.Password) error {
	return u.comms.Call("API.GetMatch", match, reply)
}

func (u *Client_RPC) AddCSV(data kit.CSVData) error {
	var reply bool
	return u.comms.Call("API.AddCSV", data, &reply)
}

func (u *Client_RPC) DeleteItem(data kit.MatchPass) (int64, error) {
	var result int64
	e := u.comms.Call("API.DeletePassword", data, &result)
	return result, e
}

const (
	PORT  = ":5555"
	PROTO = "tcp"
)

var (
	RPC_CLIENT = &Client_RPC{}
)

func init() {
	RPC_CLIENT.comms, _ = rpc.DialHTTP(PROTO, PORT)
}
