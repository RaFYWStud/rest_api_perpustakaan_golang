package repository

import (
	"golang-tutorial/contract"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *contract.Repository {
	return &contract.Repository{
		Perpus: implPerpusRepository(db),
	}
}
