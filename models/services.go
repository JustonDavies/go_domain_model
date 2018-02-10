package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Service struct {
	//-- User Variables ----------
	Name      string `gorm:"not null"`
	SubDomain string

	URLs         postgres.Jsonb
	CategoryTags postgres.Jsonb

	//-- System Variables ----------

	//-- Relations ----------
	CategoryID uint `gorm:"not null; index"`
	ProviderID uint `gorm:"not null; index"`

	//-- Automated fields (ID, Timestamps) ----------
	gorm.Model
}
