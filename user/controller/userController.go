package controller

import (
	"go_clean_arc/domain"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserController interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	GetProfile(c *gin.Context)
	AuthMiddleware() gin.HandlerFunc
}

type userController struct {
	userUsecase domain.UserUsecase
	authUsecase domain.AuthUsecase
}

func NewUserController(userUsecase domain.UserUsecase, authUsecase domain.AuthUsecase) UserController {
	return &userController{
		userUsecase: userUsecase,
		authUsecase: authUsecase,
	}
}

func (uc *userController) SignUp(c *gin.Context) {
	user := &domain.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg": "invalid format",
			},
		)
		return
	}

	if msg, ok := uc.userUsecase.ValidateUserSignup(user); !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg": msg,
			},
		)
		return
	}

	if err := uc.userUsecase.SaveUser(user); err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"email": "email already registered",
				},
			)
		} else if strings.Contains(err.Error(), "users_username_key") {
			c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"username": "username already registered",
				},
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"msg": "Internal Server Error, Try Again Later.",
				},
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"msg": uc.userUsecase.ToPublic(user),
			// "msg": user,
		},
	)
	return
}

func (uc *userController) SignIn(c *gin.Context) {
	user := &domain.User{}

	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg": "invalid input format",
			},
		)
		return
	}

	if msg, ok := uc.userUsecase.ValidateUserSignin(user); !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg": msg,
			},
		)
		return
	}

	var (
		puser *domain.PublicUser
		err   error
	)

	if puser, err = uc.userUsecase.GetUserByEmailAndPassword(user.Email, user.Password); err != nil {
		if err.Error() == "invalid password" {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"password": "invalid password",
				},
			)
			return
		} else if err == gorm.ErrRecordNotFound {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"email": "Email Not registered",
				},
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"msg": "Internal Server Error",
				},
			)
		}
		return
	}

	token := &domain.Token{}
	if token, err = uc.authUsecase.CreateToken(puser.ID); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"msg": "Internal Server Error",
			},
		)
	}

	if err = uc.authUsecase.SaveToken(puser.ID, token); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"msg": "Internal Server Error",
			},
		)
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"msg": uc.authUsecase.ToPublic(token),
		},
	)
	return
}

func (uc *userController) GetProfile(c *gin.Context) {

	var (
		accessDetail *domain.AccessDetails
		err          error
		user         *domain.User
	)

	if accessDetail, err = uc.authUsecase.ExtractTokenMetadata(c.Request); err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"msg": err.Error(),
			},
		)
		return
	}

	if user, err = uc.userUsecase.GetUserProfile(accessDetail.UserId); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg": err.Error(),
			},
		)
		return
	}

	user.Password = ""

	c.JSON(
		http.StatusOK,
		gin.H{
			"msg": user,
		},
	)

}
