package main

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var tests = []struct {
		name                    string
		url                     string
		expectedStatusCode      int
		expectedUrl             string
		expectedFirstStatusCode int
	}{
		{name: "home", url: "/", expectedStatusCode: http.StatusOK, expectedUrl: "/", expectedFirstStatusCode: http.StatusOK},
		{name: "404", url: "/fish", expectedStatusCode: http.StatusNotFound, expectedUrl: "/fish", expectedFirstStatusCode: http.StatusNotFound},
		{name: "profile", url: "/user/profile", expectedStatusCode: http.StatusOK, expectedUrl: "/", expectedFirstStatusCode: http.StatusTemporaryRedirect},
	}

	// we are getting this from setup_test.go
	// var app application

	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			// accept invalid ssl certificate
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// changing path to templates so our test can run, it is different then it would be in prod
	// pathToTemplates = "./../../templates/"
	// we do this in setup_test.go

	for _, e := range tests {
		res, err := ts.Client().Get(ts.URL + e.url)

		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if res.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected status %d, but got %d", e.name, e.expectedStatusCode, res.StatusCode)
		}

		if res.Request.URL.Path != e.expectedUrl {
			t.Errorf("%s: expected final url of %s but got %s", e.name, e.expectedUrl, res.Request.URL.Path)
		}

		res2, _ := client.Get(ts.URL + e.url)

		if res2.StatusCode != e.expectedFirstStatusCode {
			t.Errorf("%s: expected first return status code of %d, but got %d", e.name, e.expectedFirstStatusCode, res2.StatusCode)
		}
	}
}

func TestApp_Login(t *testing.T) {
	var tests = []struct {
		name               string
		postedData         url.Values
		expectedStatusCode int
		expectedLoc        string
	}{
		{
			name: "valid login",
			postedData: url.Values{
				"email": {
					"admin@example.com",
				},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/user/profile",
		},
		{
			name: "missing form data",
			postedData: url.Values{
				"email": {
					"",
				},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "Bad credentials",
			postedData: url.Values{
				"email": {
					"admin@example.com",
				},
				"password": {"32332"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "user not found",
			postedData: url.Values{
				"email": {
					"asdsda@asd.com",
				},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
	}

	for _, e := range tests {
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(e.postedData.Encode()))

		req = addContextAndSessionToRequest(req, app)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Login)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: returned wrong status code; expected %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		actualLoc, err := rr.Result().Location()

		if err == nil {
			if actualLoc.String() != e.expectedLoc {
				t.Errorf("%s: expected location %s,  but got %s", e.name, e.expectedLoc, actualLoc.String())
			}
		} else {
			t.Errorf("%s: no location header set", e.name)
		}

	}
}

// func TestAppHomeOld(t *testing.T) {
// 	// create a request
// 	req, _ := http.NewRequest("GET", "/", nil)

// 	req = addContextAndSessionToRequest(req, app)

// 	// response writer
// 	rr := httptest.NewRecorder()

// 	handler := http.HandlerFunc(app.Home)

// 	handler.ServeHTTP(rr, req)

// 	// check status code

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("TestAppHome returned wrong status code, expected 200, but got %d", rr.Code)
// 	}

// 	body, _ := io.ReadAll(rr.Body)

// 	if !strings.Contains(string(body), `<small>From Session:`) {
// 		t.Error("Did not found correct text in html")
// 	}
// }

func TestApp_Home(t *testing.T) {
	tests := []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
		{"first visit", "hello, world!", "<small>From Session: hello, world!"},
	}

	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/", nil)

		req = addContextAndSessionToRequest(req, app)

		// make sure nothing is in a session
		app.Session.Destroy(req.Context())

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

		// response writer
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)

		handler.ServeHTTP(rr, req)

		// check status code

		if rr.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code, expected 200, but got %d", rr.Code)
		}

		body, _ := io.ReadAll(rr.Body)

		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s did not find %s in the response body", e.name, e.expectedHTML)
		}
	}
}

func TestApp_RenderWithBadTemplate(t *testing.T) {
	// set template file to a location with a bac template
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)

	req = addContextAndSessionToRequest(req, app)

	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})

	if err == nil {
		t.Error("Expected error from bad template but did not get it")
	}

	pathToTemplates = "./../../templates/"
}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")

	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}
