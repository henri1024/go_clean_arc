package presenter

import "clean_arc/domain/entity"

type UserPresenter interface {
	ResponseByPublicUserDetail(u *entity.User) *entity.PublicUser
}
