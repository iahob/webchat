package user

import (
	"fmt"

	"webchat/db"
)

type UserModel struct {
	Token string `db:"token"`
	Name  string `db:"name"`
	Pwd   string `db:"pwd"`
}

func (u *UserModel) Add() {

}
func (u *UserModel) Del() {

}
func (u *UserModel) Update() {
	db.Mysql.Exec("update t_user set token = ? where name = ?", u.Token, u.Name)
}
func (u *UserModel) Query() {

}

func Login(name, pwd string) error {
	if len(name) == 0 || len(pwd) == 0 {
		return fmt.Errorf("empty name or pwd")
	}
	var dbPwd string
	err := db.Mysql.Get(&dbPwd, "select pwd from t_user where name = ? limit 1", name)
	if err != nil {
		return fmt.Errorf("error name")
	}
	if dbPwd != pwd {
		return fmt.Errorf("error pwd")
	}
	return nil
}

func Oauth(token string) (*UserModel, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("empty token")
	}
	var data UserModel
	err := db.Mysql.Get(&data, "select name,pwd,token from t_user where token = ? limit 1", token)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
