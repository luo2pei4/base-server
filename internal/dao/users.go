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
	if result := baseDB.
		Table("users").
		Where("name = ? and passwd = ?", name, passwd).
		Take(user); result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func UpdateLastLoginTime(user *Users) error {
	if result := baseDB.
		Table("users").
		Where("name = ?", user.Name).
		Update("last_login_time", user.LastLoginTime); result.Error != nil {
		return result.Error
	}
	return nil
}
