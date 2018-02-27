package models

import (
	"github.com/jinzhu/gorm"
)

type Domain struct {
	//-- User Variables ----------
	DNSName string `gorm:"not null;unique"valid:"required,dns"`

	//-- System Variables ----------

	//-- Relations ----------
	//CategoryID uint `gorm:"not null; index"`

	//-- Automated fields (ID, Timestamps) ----------
	gorm.Model
}
