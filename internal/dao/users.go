package dao

import (
	"time"
)

type Users struct {
	Name          string
	Passwd        string
	LastLoginTime time.Time
}

func QueryUser(name, passwd string) (*Users, error) {
	user := &Users{}
	result := baseServerDB.Raw(
		"select name, passwd, last_login_time from users where name=? and passwd=?",
		name,
		passwd).Scan(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func UpdateLastLoginTime(user *Users) error {
	result := baseServerDB.Exec(
		"update users set last_login_time=? where name=? and passwd=?",
		user.LastLoginTime,
		user.Name,
		user.Passwd,
	)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
