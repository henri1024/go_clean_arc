package repository

import (
	"encoding/json"
	"testing"
	"time"
	"userauth/domain"
	"userauth/infrastructure/uuidgenerator"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RedisSuite struct {
	suite.Suite
	Host     string
	Password string
	DB       int
	Client   *redis.Client
}

func (r *RedisSuite) SetupSuite() {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.Host,
		Password: r.Password,
		DB:       r.DB,
	})
}
func (r *RedisSuite) TearDownSuite() {
	r.Client.Close()
}

type redisHandlerSuite struct {
	RedisSuite // Embed RedisSuite struct
}

func TestRedisSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test for redis repository")
	}
	redisHostTest := "localhost:5432"
	redisPassword := "satuduatiga"
	// redisDb, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
	redisHandlerSuiteTest := &redisHandlerSuite{
		RedisSuite{
			Host:     redisHostTest,
			Password: redisPassword,
			DB:       0,
		},
	}
	suite.Run(t, redisHandlerSuiteTest)
}

func getItemByKey(client *redis.Client, key string) ([]byte, error) {
	return client.Get(key).Bytes()
}

func seedItem(client *redis.Client, key string, value interface{}) error {
	jybt, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(key, jybt, time.Second*30).Err()
}

func (r *redisHandlerSuite) TestSaveToken() {

	item := struct {
		id    uint
		token *domain.Token
	}{

		1,
		&domain.Token{
			AccessExpired:  time.Now().Add(10 * time.Minute).Unix(),
			RefreshExpired: time.Now().Add(7 * time.Hour * 24).Unix(),
			TokenUuid:      uuidgenerator.NewUuid(),
			RefreshUuid:    uuidgenerator.NewUuid(),
		},
	}

	repo := NewTokenRepository(r.Client)
	err := repo.SaveToken(item.id, item.token)
	require.NoError(r.T(), err)

	jbyt, err := getItemByKey(r.Client, item.token.TokenUuid)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), jbyt)

	assert.Equal(r.T(), item.id, jbyt)

	jbyt, err = getItemByKey(r.Client, item.token.RefreshUuid)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), jbyt)
	assert.Equal(r.T(), item.id, jbyt)
}
