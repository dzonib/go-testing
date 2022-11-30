package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webapp/pkg/data"
)

func Test_application_addIPToContext(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false},
		{"", "", "hello:world", false},
	}
	// we are getting this from setup_test.go
	// var app application

	// create dummy handler that we'll use to check te context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure that the value exists in the context
		val := r.Context().Value(contextUserKey)

		if val == nil {
			t.Error(contextUserKey, "not present")
		}

		// make sure we got a string back
		ip, ok := val.(string)

		if !ok {
			t.Error("Not string")
		}

		t.Log(ip)

	})

	for _, e := range tests {
		// create a handler to test
		handlerToTest := app.addIPToContext(nextHandler)

		// we need a request

		req := httptest.NewRequest("GET", "http://testing", nil)

		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headerName) > 0 {
			req.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}

}

func Test_application_IPFromContext(t *testing.T) {
	//create a app variable of type application

	// we are getting this from setup_test.go
	// var app application

	// get a context
	ctx := context.Background()

	// put something in a context
	ctx = context.WithValue(ctx, contextUserKey, "whatever")

	// call a function

	ip := app.ipFromContext(ctx)
	// preform the test
	if !strings.EqualFold("whatever", ip) {
		t.Errorf("Expected %s to be whatever", ip)
	}
}

func Test_app_auth(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	var tests = []struct {
		name   string
		isAuth bool
	}{
		{"logged in", true},
		{"not logged in", false},
	}

	for _, e := range tests {
		handlerToTest := app.auth(nextHandler)
		req, _ := http.NewRequest(http.MethodGet, "http://testing", nil)
		req = addContextAndSessionToRequest(req, app)

		if e.isAuth {
			app.Session.Put(req.Context(), "user", data.User{ID: 1})
		}

		rr := httptest.NewRecorder()

		handlerToTest.ServeHTTP(rr, req)

		if e.isAuth && rr.Code != http.StatusOK {
			t.Errorf("%s: expected status code of 200, but got %d", e.name, rr.Code)
		}

		if !e.isAuth && rr.Code != http.StatusTemporaryRedirect {
			t.Errorf("%s: expected status %d, but got %d", e.name, http.StatusTemporaryRedirect, rr.Code)
		}
	}
}
