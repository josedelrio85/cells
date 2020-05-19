package leads

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
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
			Description: "when DuplicatedTime hook is successfully activated",
			Lead: Lead{
				SouID: 70,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedTime hook is successfully activated",
			Lead: Lead{
				SouID: 74,
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
		{
			Description: "when DuplicatedTime hook is successfully activated",
			Lead: Lead{
				SouID: 63,
			},
			Active: true,
		},
		{
			Description: "when DuplicatedTime hook is successfully activated",
			Lead: Lead{
				SouID: 69,
			},
			Active: true,
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

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	redis := redisclient.Redis{
		Pool: &redis.Pool{
			MaxIdle:     5,
			IdleTimeout: 60 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", s.Addr())
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

			if !test.ExpectedResult {
				expirationtime := duplicated.getExpirationTime(test.Lead.SouID)
				phone := *test.Lead.LeaPhone
				key := fmt.Sprintf("%s-%d-%d", phone, test.Lead.SouID, test.Lead.LeatypeID)
				cont.Redis.Set(key, phone, expirationtime)
			}

			var response HookResponse
			if test.Sleep {
				keyalt := fmt.Sprintf("%s-%d-%d", *test.Lead.LeaPhone, test.Lead.SouID, test.Lead.LeatypeID)

				exptime := time.Duration((duplicated.getExpirationTime(test.Lead.SouID))) * time.Second
				exptime2 := time.Duration((duplicated.getExpirationTime(test.Lead.SouID) + 1)) * time.Second

				s.SetTTL(keyalt, exptime)
				s.FastForward(exptime2 * time.Second)

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
