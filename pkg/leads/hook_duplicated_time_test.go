package leads

import (
	"net/http"
	"testing"
	"time"

	redisclient "github.com/bysidecar/leads/pkg/leads/redis"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestActiveDuplicated(t *testing.T) {
	assert := assert.New(t)

	var duplicated DuplicatedTime

	tests := []struct {
		Description string
		Lead        Lead
		Active      bool
	}{
		{
			Description: "when DuplicatedTime hook is successfully activated",
			Lead: Lead{
				SouID: 64,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedTime hook is successfully activated",
			Lead: Lead{
				SouID: 15,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedTime hook is not activated",
			Lead: Lead{
				SouID: 1,
			},
			Active: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			active := duplicated.Active(test.Lead)

			assert.Equal(test.Active, active)
		})
	}
}

func TestPerformDuplicated(t *testing.T) {
	assert := assert.New(t)

	var duplicated DuplicatedTime
	redis := redisclient.Redis{
		Pool: &redis.Pool{
			MaxIdle:     5,
			IdleTimeout: 60 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", GetSetting("CHECK_LEAD_REDIS")+":6379")
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	phone1 := HelperRandstring(9)
	phone2 := HelperRandstring(9)

	tests := []struct {
		Description    string
		Lead           Lead
		Response       HookResponse
		Sleep          bool
		ExpectedResult bool
	}{
		{
			Description: "When a lead is not duplicated",
			Lead: Lead{
				LeaPhone:  &phone1,
				SouID:     64,
				LeatypeID: 1,
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			Sleep:          false,
			ExpectedResult: true,
		},
		{
			Description: "When a lead is duplicated because reached the limit",
			Lead: Lead{
				LeaPhone:  &phone1,
				SouID:     64,
				LeatypeID: 1,
			},
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New("Error"),
			},
			Sleep:          false,
			ExpectedResult: false,
		},
		{
			Description: "When another lead is not duplicated",
			Lead: Lead{
				LeaPhone:  &phone2,
				SouID:     15,
				LeatypeID: 8,
			},
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			Sleep:          false,
			ExpectedResult: true,
		},
		{
			Description: "When a lead of the same class arrives at a period of time minor than the expiration time",
			Lead: Lead{
				LeaPhone:  &phone2,
				SouID:     15,
				LeatypeID: 8,
			},
			Sleep: false,
			Response: HookResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Err:        errors.New("Error"),
			},
			ExpectedResult: false,
		},
		{
			Description: "When a lead of the same class arrives at a period of time greater than the expiration time",
			Lead: Lead{
				LeaPhone:  &phone2,
				SouID:     15,
				LeatypeID: 8,
			},
			Sleep: true,
			Response: HookResponse{
				StatusCode: http.StatusOK,
				Err:        nil,
			},
			ExpectedResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			cont := Handler{
				Lead:  test.Lead,
				Redis: redis,
			}
			var response HookResponse
			if test.Sleep {
				// exptime := time.Duration(duplicated.getExpirationTime(test.Lead.SouID)) * time.Second
				exptime := time.Duration((duplicated.getExpirationTime(test.Lead.SouID) + 1)) * time.Second
				time.Sleep(exptime)
				response = duplicated.Perform(&cont)
			} else {
				response = duplicated.Perform(&cont)
			}

			assert.Equal(test.Response.StatusCode, response.StatusCode)
			if test.ExpectedResult {
				assert.Nil(response.Err)
			} else {
				assert.NotNil(response.Err)
			}
		})
	}
}
