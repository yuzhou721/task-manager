package models

import "github.com/jinzhu/gorm"

type file struct {
	gorm.Model
	name string
	ext  string
}
