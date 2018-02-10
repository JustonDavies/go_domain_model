package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Category struct {
	//-- User Variables ----------
	Name string `gorm:"not null;unique"`
	Tags postgres.Jsonb

	//-- System Variables ----------

	//-- Relations ----------
	Services []Service

	//-- Automated fields (ID, Timestamps) ----------
	gorm.Model
}
