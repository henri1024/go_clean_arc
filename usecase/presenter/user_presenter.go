package presenter

import "clean_arc/domain/entity"

type UserPresenter interface {
	ResponseSave(u *entity.User) *entity.PublicUser
}
