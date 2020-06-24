package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserPrepare(t *testing.T) {
	tmpUser := &User{
		Email: "    test@example.com 		",
		Username: "\t\t<myusername\n",
	}

	tmpUser.Prepare()

	assert.Equal(t, "test@example.com", tmpUser.Email)
	assert.Equal(t, "&lt;myusername", tmpUser.Username)
}
