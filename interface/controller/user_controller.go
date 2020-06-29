package controller

import (
	"clean_arc/domain/entity"
	"clean_arc/usecase/interactor"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type userController struct {
	interactor.TokenInteractor
	interactor.AuthInteractor
	userInteractor interactor.UserInteractor
}

type UserController interface {
	SaveUser(c *gin.Context)
	GetUserByEmailAndPassowrd(c *gin.Context)
}

func NewUserController(ui interactor.UserInteractor, ai interactor.AuthInteractor, ti interactor.TokenInteractor) UserController {
	return &userController{
		AuthInteractor:  ai,
		TokenInteractor: ti,
		userInteractor:  ui,
	}
}

func (uc *userController) SaveUser(c *gin.Context) {
	tempUser := &entity.User{}

	if err := c.ShouldBindJSON(tempUser); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid input format",
			},
		)
		return
	}

	if resp, ok := tempUser.ValidSave(); !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": resp,
			},
		)
		return
	}

	user, err := uc.userInteractor.Save(tempUser)
	if err != nil {

		if strings.Contains(err.Error(), "users_email_key") {
			c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"invalid_email": "email already registered",
				},
			)
		} else if strings.Contains(err.Error(), "users_username_key") {
			c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"invalid_username": "username already registered",
				},
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"internal_server_error": "please contact admin",
				},
			)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": user,
	})
	return
}

func (uc *userController) GetUserByEmailAndPassowrd(c *gin.Context) {
	tempUser := &entity.User{}

	if err := c.ShouldBindJSON(tempUser); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid input format",
			},
		)
		return
	}

	if resp, ok := tempUser.ValidGetByEmailAndPassword(); !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": resp,
			},
		)
		return
	}

	user, err := uc.userInteractor.GetUserByEmailAndPassword(tempUser.Email, tempUser.Password)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"invalid_email": "Email Not registered",
				},
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"internal_server_error": "please contact admin",
				},
			)
		}
		return
	}

	ok := user.ComparePassword(tempUser.Password)
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"invalid_password": "invalid password",
			},
		)
		return
	}

	token, err := uc.TokenInteractor.CreateToken(user.ID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Message": err.Error(),
			},
		)
	}

	err = uc.AuthInteractor.SaveAuthToken(user.ID, token)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Message": err.Error(),
			},
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": token})

}
