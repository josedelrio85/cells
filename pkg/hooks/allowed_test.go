package leads

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	container "github.com/bysidecar/leads/pkg/container"
	model "github.com/bysidecar/leads/pkg/model"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestActiveAllowed(t *testing.T) {
	assert := assert.New(t)

	var allowed Allowed

	tests := []struct {
		Description string
		Lead        model.Lead
		Active      bool
	}{
		{
			Description: "when Allowed hook is successfully activated",
			Lead: model.Lead{
				SouID: 64,
			},
			Active: true,
		},
		{
			Description: "when Allowed hook is not activated",
			Lead: model.Lead{
				SouID: 1,
			},
			Active: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := allowed.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerformAllowed(t *testing.T) {
	assert := assert.New(t)

	var allowed Allowed
	redis := model.Redis{
		Pool: &redis.Pool{
			MaxIdle:     5,
			IdleTimeout: 60 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", getSettingAllowed("CHECK_LEAD_REDIS")+":6379")
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	phone1 := "666666666"
	phone2 := "666666667"

	tests := []struct {
		Description    string
		Lead           model.Lead
		Response       HookResponse
		ExpectedResult bool
	}{
		{
			Description: "When a lead is allowed",
			Lead: model.Lead{
				LeaPhone:  &phone1,
				SouID:     64,
				LeatypeID: 1,
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: true,
		},
		{
			Description: "When a lead is not allowed because reached the limit",
			Lead: model.Lead{
				LeaPhone:  &phone1,
				SouID:     64,
				LeatypeID: 1,
			},
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New("Error"),
			},
			ExpectedResult: false,
		},
		{
			Description: "When another lead is allowed",
			Lead: model.Lead{
				LeaPhone:  &phone2,
				SouID:     66,
				LeatypeID: 8,
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			cont := container.Container{
				Lead:  test.Lead,
				Redis: redis,
			}
			response := allowed.Perform(cont)

			assert.Equal(test.Response.StatusCode, response.StatusCode)
			if test.ExpectedResult {
				assert.Nil(response.Err)
			} else {
				assert.NotNil(response.Err)
			}
		})
	}
}

func getSettingAllowed(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}
