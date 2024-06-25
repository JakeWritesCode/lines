package http

import (
	linesHttp "lines/lines/http"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserLogin) Validate() linesHttp.HttpError {
	httpErr := linesHttp.HttpError{}
	if u.Email == "" {
		httpErr.Message = append(httpErr.Message, "Email is required.")
	}
	if u.Password == "" {
		httpErr.Message = append(httpErr.Message, "Password is required.")
	}
	return httpErr
}
