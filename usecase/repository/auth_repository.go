package repository

import (
	"clean_arc/domain/authtoken"
)

type AuthRepository interface {
	SaveToken(uint, *authtoken.Token) error
}
