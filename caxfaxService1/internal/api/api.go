package api

import "caxfaxService1/internal/entity"

type APIClient interface {
	GetFunFact() (entity.Fact, error)
}
