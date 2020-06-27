package repository

import "clean_arc/domain/authtoken"

type TokenRepository interface {
	CreateToken(uint) (*authtoken.Token, error)
}
