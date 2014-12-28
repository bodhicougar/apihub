package api

import (
	"encoding/json"
	"net/http"

	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/auth"
	. "github.com/backstage/backstage/errors"
	"github.com/zenazn/goji/web"
)

type UsersHandler struct {
	ApiHandler
}

func (handler *UsersHandler) CreateUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, user)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	err = user.Save()
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	return Created(user.ToString())
}

func (handler *UsersHandler) DeleteUser(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user, err := GetCurrentUser(c)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	auth.RevokeTokensFor(user)
	user.Delete()
	return OK(user.ToString())
}

func (handler *UsersHandler) Login(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	user := &User{}
	err := handler.parseBody(r.Body, user)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, err.Error())
	}
	token, err := LoginAndGetToken(user)
	if err != nil {
		return BadRequest(E_BAD_REQUEST, ErrAuthenticationFailed.Error())
	}
	payload, _ := json.Marshal(token)
	return OK(string(payload))
}
