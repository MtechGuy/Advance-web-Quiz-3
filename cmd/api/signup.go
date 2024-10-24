package main

import (
	"errors"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment
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
	// create a struct to hold a comment
	// we use struct tags to make the names display in lowercase
	var incomingData struct {
		Email string `json:"email"`
		FName string `json:"fname"`
		MName string `json:"mname"`
		LName string `json:"lname"`
	}
	// perform the decoding
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
	// Initialize a Validator instance
	v := validator.New()

	data.ValidateSignup(v, signup)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors) // implemented later
		return
	}
	err = a.signupModel.Insert(signup)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Set a Location header. The path to the newly created comment
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
	// Get the id from the URL /v1/comments/:id so that we
	// can use it to query teh comments table. We will
	// implement the readIDParam() function later
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call Get() to retrieve the comment with the specified id
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

	// display the comment
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
	// Get the ID from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Retrieve the comment from the database
	signup, err := a.signupModel.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			a.notFoundResponse(w, r)
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Decode the incoming JSON
	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Update the comment fields based on the incoming data
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
	// Validate the updated comment
	v := validator.New()
	data.ValidateSignup(v, signup)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Perform the update in the database
	err = a.signupModel.Update(signup)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the updated comment
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

	err = a.signupModel.Delete(id) // Removed the '&' operator here
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.IDnotFound(w, r, id) // Pass the ID to the custom message handler
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
