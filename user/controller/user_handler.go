package controller

import (
	"fmt"
	"net/http"
	"userauth/domain"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserController struct {
	UserUsecase  domain.UserUsecase
	TokenUsecase domain.TokenUsecase
}

func NewUserController(userUsecase domain.UserUsecase, tokenUsecase domain.TokenUsecase) UserController {
	handler := UserController{
		UserUsecase:  userUsecase,
		TokenUsecase: tokenUsecase,
	}
	return handler
}

func (uc *UserController) SignUp(c *gin.Context) {
	tmpUser := &domain.User{}

	if err := c.ShouldBindJSON(tmpUser); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid input format",
			},
		)
		return
	}

	if resp, ok := tmpUser.ValidSave(); !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": resp,
			},
		)
		return
	}

	err := uc.UserUsecase.SaveUser(tmpUser)
	if err != nil {
		if strings.Contains(err.Error(), "failed to hash password") {
			c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"invalid_password": "cant hash password, try another",
				},
			)
		} else if strings.Contains(err.Error(), "users_email_key") {
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
		"message": "user created",
	})
	return

}

func (uc *UserController) Login(c *gin.Context) {
	tmpUser := &domain.User{}

	if err := c.ShouldBindJSON(tmpUser); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid input format",
			},
		)
		return
	}

	if resp, ok := tmpUser.ValidGetByEmailAndPassword(); !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": resp,
			},
		)
		return
	}

	user, err := uc.UserUsecase.GetUserByEmailAndPassword(tmpUser.Email, tmpUser.Password)

	if err != nil {

		if err.Error() == "invalid password" {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"invalid_password": "invalid password",
				},
			)
			return
		}

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

	token, err := uc.TokenUsecase.CreateToken(user.ID)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"internal_server_error": "user created but failed to create token",
			},
		)
		return
	}

	err = uc.TokenUsecase.SaveToken(user.ID, token)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"internal_server_error": "please contact admin",
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": token.ToPublic()})

}

func (uc *UserController) LogOut(c *gin.Context) {
	metadata, err := uc.TokenUsecase.ExtractTokenMetadata(c.Request)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = uc.TokenUsecase.DeleteToken(metadata.AccessUuid)

	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")

}

func (uc *UserController) Refresh(c *gin.Context) {

	mapToken := map[string]string{}

	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	token, err := uc.TokenUsecase.VerifyToken(refreshToken, "refresh")
	if err != nil {
		fmt.Println("its here 1")
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	err = uc.TokenUsecase.IsValidToken(token)
	if err != nil {
		fmt.Println("its here 2")
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	// Valid token
	claims, ok := token.Claims.(jwt.MapClaims)

	refreshUUID := claims["refresh_uuid"].(string)
	err = uc.TokenUsecase.DeleteToken(refreshUUID)
	if err != nil {
		fmt.Println("its here 3")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if ok {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		userID := uint(uid)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}

		token, err := uc.TokenUsecase.CreateToken(userID)

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"internal_server_error": "user created but failed to create token",
				},
			)
			return
		}

		err = uc.TokenUsecase.SaveToken(userID, token)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"internal_server_error": "please contact admin",
				},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": token.ToPublic()})
		return
	}
	c.JSON(http.StatusUnauthorized, "refresh expired")
}
