package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mtechguy/quiz3/internal/data"
	"github.com/mtechguy/quiz3/internal/validator"
)

var incomingData struct {
	Email    *string `json:"email"`
	FName    *string `json:"fname"`
	MName    *string `json:"mname"`
	LName    *string `json:"lname"`
	FullName *string
}

func (a *applicationDependencies) createSignupHandler(w http.ResponseWriter, r *http.Request) {

	var incomingData struct {
		Email string `json:"email"`
		FName string `json:"fname"`
		MName string `json:"mname"`
		LName string `json:"lname"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	signup := &data.Signup{
		Email:    incomingData.Email,
		FName:    incomingData.FName,
		MName:    incomingData.MName,
		LName:    incomingData.LName,
		FullName: incomingData.FName + " " + incomingData.MName + " " + incomingData.LName,
	}

	v := validator.New()

	data.ValidateSignup(v, signup)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = a.signupModel.Insert(signup)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/signup/%d", signup.ID))

	data := envelope{
		"signup": signup,
	}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}

func (a *applicationDependencies) displaySignupHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	signup, err := a.signupModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"signup": signup,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}

func (a *applicationDependencies) updateSignupHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	signup, err := a.signupModel.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			a.notFoundResponse(w, r)
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.Email != nil {
		signup.Email = *incomingData.Email
	}
	if incomingData.FName != nil {
		signup.FName = *incomingData.FName
	}

	if incomingData.MName != nil {
		signup.MName = *incomingData.MName
	}

	if incomingData.LName != nil {
		signup.LName = *incomingData.LName
	}

	signup.FullName = signup.FName + " " + signup.MName + " " + signup.LName

	v := validator.New()
	data.ValidateSignup(v, signup)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.signupModel.Update(signup)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"signup": signup,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteSignupHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.signupModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.IDnotFound(w, r, id)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"message": "signup successfully deleted",
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listSignupHandler(w http.ResponseWriter, r *http.Request) {
	signup, err := a.signupModel.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"signup": signup,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
