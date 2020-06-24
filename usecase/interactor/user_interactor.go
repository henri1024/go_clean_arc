package interactor

import (
	"clean_arc/domain/entity"
	"clean_arc/usecase/presenter"
	"clean_arc/usecase/repository"
)

type userInteractor struct {
	UserRepository repository.UserRepository
	UserPresenter  presenter.UserPresenter
}

type UserInteractor interface {
	Save(u *entity.User) (*entity.PublicUser, error)
}

func NewUserInteractor(r repository.UserRepository, p presenter.UserPresenter) UserInteractor {
	return &userInteractor{
		UserRepository: r,
		UserPresenter:  p,
	}
}

func (ui *userInteractor) Save(user *entity.User) (*entity.PublicUser, error) {
	user, err := ui.UserRepository.SaveUser(user)
	if err != nil {
		return nil, err
	}

	return ui.UserPresenter.ResponseSave(user), nil
}
