package foohandler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TEST HANDLERS WITHOUT SERVER
// RR => response recorder
func TestHandleGetFooRR(t *testing.T) {
	// recorder will capture everything we write to our response writer (header, bytes....)
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "", nil)

	if err != nil {
		t.Error(err)
	}
	handleGetFoo(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected 200 but got %d", rr.Result().StatusCode)
	}

	defer rr.Result().Body.Close()

	expected := "FOO"

	b, err := ioutil.ReadAll(rr.Result().Body)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s, but we got %s", expected, string(b))
	}
}

// TEST HANDLERS WITH SERVER
func TestHandleGetFoo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleGetFoo))

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 but got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	expected := "FOO"

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s, but we got %s", expected, string(b))
	}
}
