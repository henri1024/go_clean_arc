package interactor

import (
	"clean_arc/domain/authtoken"
	"clean_arc/usecase/repository"
)

type tokenInteractor struct {
	TokenRepository repository.TokenRepository
}

type TokenInteractor interface {
	CreateToken(uint) (*authtoken.Token, error)
}

func NewTokenInteractor(r repository.TokenRepository) TokenInteractor {
	return &tokenInteractor{
		TokenRepository: r,
	}
}

func (ti *tokenInteractor) CreateToken(userid uint) (*authtoken.Token, error) {
	token, err := ti.TokenRepository.CreateToken(userid)
	if err != nil {
		return nil, err
	}
	return token, nil
}
