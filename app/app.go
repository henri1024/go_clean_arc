package app

import (
	auth_repository "go_clean_arc/auth/repository"
	auth_usecase "go_clean_arc/auth/usecase"
	"go_clean_arc/domain"
	"go_clean_arc/infrastructure/hash"
	"go_clean_arc/infrastructure/jwtAuth"
	uuid "go_clean_arc/infrastructure/uuid"
	user_controller "go_clean_arc/user/controller"
	user_repository "go_clean_arc/user/repository"
	user_usecase "go_clean_arc/user/usecase"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type App struct {
	UserController user_controller.UserController
}

func NewControllers(db *gorm.DB, hasher *hash.Hasher, idGen *uuid.UuidGenerator, rdb *redis.Client, jwtauth *jwtAuth.JwtWidget) *App {
	return &App{
		UserController: CreateUserController(db, hasher, CreateAuthUsecase(idGen, rdb, jwtauth)),
	}
}

func CreateUserController(db *gorm.DB, hasher *hash.Hasher, authusecase domain.AuthUsecase) user_controller.UserController {
	userRepo := user_repository.NewUserRepository(db)
	userUsecase := user_usecase.NewUserUsecase(userRepo, hasher)
	return user_controller.NewUserController(userUsecase, authusecase)
}

func CreateAuthUsecase(idGenerator *uuid.UuidGenerator, rdb *redis.Client, jwtauth *jwtAuth.JwtWidget) domain.AuthUsecase {
	AuthRepo := auth_repository.NewAuthRepository(rdb)
	return auth_usecase.NewAuthUsecase(AuthRepo, idGenerator, jwtauth)
}
