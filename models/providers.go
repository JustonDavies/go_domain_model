package models

import (
	"github.com/jinzhu/gorm"
)

type Provider struct {
	//-- User Variables ----------
	Name   string `gorm:"not null;unique"valid:"required"`
	Domain string `gorm:"not null;unique"valid:"required,dns"`

	//-- System Variables ----------

	//-- Relations ----------
	Services []Service

	//-- Automated fields (ID, Timestamps) ----------
	gorm.Model
}

func (provider Provider) Validate(db *gorm.DB) {
	//if provider.Domain != "SATAN" {
	//	db.AddError(errors.New("all domains must be satan"))
	//}
}
