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

type UserSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserSignUp) Validate() linesHttp.HttpError {
	httpErr := linesHttp.HttpError{}
	if u.Email == "" {
		httpErr.Message = append(httpErr.Message, "Email is required.")
	}
	if u.Password == "" {
		httpErr.Message = append(httpErr.Message, "Password is required.")
	}
	if u.Name == "" {
		httpErr.Message = append(httpErr.Message, "Name is required.")
	}
	return httpErr
}

type UserReadDTO struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
