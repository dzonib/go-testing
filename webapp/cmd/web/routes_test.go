package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func Test_application_routes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{route: "/", method: "GET"},
		{route: "/login", method: "POST"},
		{route: "/user/profile", method: "GET"},
		{route: "/static/*", method: "GET"},
	}

	// we are getting this from setup_test.go
	// var app application

	mux := app.routes()

	chiRoutes := mux.(chi.Routes)

	for _, route := range registered {
		// check and see if a route exists
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})
	return found
}

func Test_app_Auth(t *testing.T) {

}
