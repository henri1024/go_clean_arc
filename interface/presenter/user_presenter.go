package presenter

import "clean_arc/domain/entity"

type userPresenter struct{}

type UserPresenter interface {
	ResponseSave(u *entity.User) *entity.PublicUser
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) ResponseSave(user *entity.User) *entity.PublicUser {
	return user.ToPublic()
}
