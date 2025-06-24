package service

import "golang-tutorial/contract"

func New(repo *contract.Repository) *contract.Service {
	return &contract.Service{
		Perpus: implPerpusService(repo),
	}
}
