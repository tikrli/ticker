package api_test

import (
	"testing"
	"github.com/appleboy/gofight"
	"strings"
	"git.codecoop.org/systemli/ticker/internal/api"
	"github.com/stretchr/testify/assert"
	"git.codecoop.org/systemli/ticker/internal/model"
	"encoding/json"
	"git.codecoop.org/systemli/ticker/internal/storage"
)

func TestGetUsersHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		var response struct {
			Data   map[string][]model.UserResponse `json:"data"`
			Status string                          `json:"status"`
			Error  interface{}                     `json:"error"`
		}

		err := json.Unmarshal(r.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, response.Status)
		assert.Equal(t, nil, response.Error)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, 2, len(response.Data["users"]))
		assert.Equal(t, "louis@systemli.org", response.Data["users"][0].Email)
		assert.Equal(t, "admin@systemli.org", response.Data["users"][1].Email)
	})

	r.GET("/v1/admin/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestGetUserHandler(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/users/2000").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`, strings.TrimSpace(r.Body.String()))
	})

	r.GET("/v1/admin/users/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)
	})

	r.GET("/v1/admin/users/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestPostUserHandler(t *testing.T) {
	r := setup()

	body := `{
		"email": "user@systemli.org",
		"password": "password12",
		"is_super_admin": true
	}`

	r.POST("/v1/admin/users").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		var response struct {
			Data   map[string]model.UserResponse `json:"data"`
			Status string                        `json:"status"`
			Error  interface{}                   `json:"error"`
		}

		err := json.Unmarshal(r.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, response.Status)
		assert.Equal(t, nil, response.Error)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, "user@systemli.org", response.Data["user"].Email)
		assert.True(t, response.Data["user"].IsSuperAdmin)
	})

	r.POST("/v1/admin/users").
		SetBody(body).
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})
}

func TestPutUserHandler(t *testing.T) {
	r := setup()

	body := `{
		"email": "new@systemli.org",
		"password": "password13",
		"role": "user",
		"is_super_admin": true,
		"tickers": [1,2,3]
	}`

	r.PUT("/v1/admin/users/2").
		SetBody(body).
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})

	r.PUT(`/v1/admin/users/2`).
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		var response struct {
			Data   map[string]model.UserResponse `json:"data"`
			Status string                        `json:"status"`
			Error  interface{}                   `json:"error"`
		}

		err := json.Unmarshal(r.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, response.Status)
		assert.Equal(t, nil, response.Error)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, 2, response.Data["user"].ID)
		assert.Equal(t, "new@systemli.org", response.Data["user"].Email)
		assert.Equal(t, "user", response.Data["user"].Role)
		assert.True(t, response.Data["user"].IsSuperAdmin)
		assert.Equal(t, []int{1, 2, 3}, response.Data["user"].Tickers)

		var user model.User
		err = storage.DB.One("ID", 2, &user)
		if err != nil {
			t.Fail()
		}

		assert.NotEmpty(t, user.EncryptedPassword)
		assert.Equal(t, true, user.IsSuperAdmin)
		assert.Equal(t, []int{1, 2, 3}, user.Tickers)
	})
}

func TestDeleteUserHandler(t *testing.T) {
	r := setup()

	r.DELETE("/v1/admin/users/3").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.DELETE("/v1/admin/users/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + UserToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 403, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1003,"message":"insufficient permissions"}}`, strings.TrimSpace(r.Body.String()))
	})

	r.DELETE("/v1/admin/users/2").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		var jres struct {
			Data   map[string]interface{} `json:"data"`
			Status string                 `json:"status"`
			Error  interface{}            `json:"error"`
		}

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Nil(t, jres.Data)
		assert.Nil(t, jres.Error)
	})

	r.DELETE("/v1/admin/users/1").
		SetHeader(map[string]string{"Authorization": "Bearer " + AdminToken}).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 400, r.Code)
		assert.Equal(t, `{"data":{},"status":"error","error":{"code":1000,"message":"self deletion is forbidden"}}`, strings.TrimSpace(r.Body.String()))

	})
}
