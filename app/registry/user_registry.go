package registry

import (
	"userauth/domain"
	"userauth/user/controller"
	userrepo "userauth/user/repository"
	userusecase "userauth/user/usecase"

	tokenrepo "userauth/token/repository"
	tokenusecase "userauth/token/usecase"
)

func (r *registry) NewUserController() controller.UserController {
	return controller.NewUserController(
		r.NewUserUsecase(),
		r.NewTokenUsecase(),
	)
}

func (r *registry) NewUserUsecase() domain.UserUsecase {
	return userusecase.NewUserUsecase(
		r.NewUserRepository(),
	)
}

func (r *registry) NewUserRepository() domain.UserRepository {
	return userrepo.NewUserRepository(
		r.db,
	)
}

func (r *registry) NewTokenUsecase() domain.TokenUsecase {
	return tokenusecase.NewTokenUsecase(
		r.NewTokenRepository(),
	)
}

func (r *registry) NewTokenRepository() domain.TokenRepository {
	return tokenrepo.NewTokenRepository(
		r.rdb,
	)
}
