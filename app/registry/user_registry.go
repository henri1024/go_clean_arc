package registry

import (
	"store/domain"
	"store/user/controller"
	userrepo "store/user/repository"
	userusecase "store/user/usecase"

	tokenrepo "store/token/repository"
	tokenusecase "store/token/usecase"
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
