package model

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type User struct {
	gorm.Model
	UserName     string
	Password     []byte
	Name         string
	LastName     string
	Email        string
	Applications []Application
}

type Application struct {
	gorm.Model
	ApplicationName string
	Roles           pq.StringArray `gorm:"type:varchar(100)[]"`
	UserID          uint
}

type UserDTO struct {
	ID                uint                 `json:"id"`
	UserName          string               `json:"userName"`
	Password          string               `json:"password,omitempty"`
	Name              string               `json:"name"`
	LastName          string               `json:"lastName"`
	Email             string               `json:"eMail"`
	Applications      []ApplicationRoleDTO `json:"applicationRoleDTO"`
	ClearApplications bool                 `json:"clearApplications,omitempty"`
}

type ApplicationRoleDTO struct {
	ApplicationName string   `json:"applicationName"`
	Roles           []string `json:"roles"`
}
