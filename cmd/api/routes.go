// Filename: cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/", a.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/signup", a.listSignupHandler)
	router.HandlerFunc(http.MethodPost, "/signup", a.createSignupHandler)
	router.HandlerFunc(http.MethodGet, "/signup/:id", a.displaySignupHandler)
	router.HandlerFunc(http.MethodPatch, "/signup/:id", a.updateSignupHandler)
	router.HandlerFunc(http.MethodDelete, "/signup/:id", a.deleteSignupHandler)

	return a.recoverPanic(router)

}
