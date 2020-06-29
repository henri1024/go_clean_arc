package interactor

import (
	"clean_arc/domain/authtoken"
	"clean_arc/usecase/repository"
)

type authInteractor struct {
	AuthRepository repository.AuthRepository
}

type AuthInteractor interface {
	SaveAuthToken(uint, *authtoken.Token) error
}

func NewAuthInteractor(r repository.AuthRepository) AuthInteractor {
	return &authInteractor{
		AuthRepository: r,
	}
}

func (ai authInteractor) SaveAuthToken(userid uint, token *authtoken.Token) error {
	return ai.AuthRepository.SaveToken(userid, token)
}
