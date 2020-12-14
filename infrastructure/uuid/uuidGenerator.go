package uuidgenerator

import "github.com/google/uuid"

type UuidGenerator struct {
}

func NewUuidGenerator() *UuidGenerator {
	return &UuidGenerator{}
}

func (ug *UuidGenerator) NewId() string {
	return uuid.New().String()

}
