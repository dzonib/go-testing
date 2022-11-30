package main

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)

	// nothing there, it does not exist
	has := form.Has("whatever")

	if has {
		t.Error("Form shows has field, when should not.")
	}

	// creating data for form
	postedData := url.Values{}

	//add name and value
	postedData.Add("a", "b")
	postedData.Add("some key", "")

	// initialize form with posted data
	form = NewForm(postedData)

	has = form.Has("a")

	if !has {
		t.Error("Form should contain 'a', and it does not")
	}

	// this is separate function so it needsa separate tests
	// form.Required("some key")

	// errorsSlice := form.Errors.Get("some key")

	// if len(errorsSlice) == 0 {
	// 	t.Error("Should have validation errors, but it has not")
	// }
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := NewForm(r.PostForm)

	form.Required("a", "b", "c")
	// we added fields as required, this should be false
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}

	postedData.Add("key a", "value a")
	postedData.Add("key b", "value b")
	postedData.Add("key c", "value c")

	r = httptest.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData

	form = NewForm(r.PostForm)

	form.Required("key a", "key b", "key c")

	if !form.Valid() {
		t.Error("form shows invalid when required fields are present")
	}
}

func TestForm_Check(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")

	if form.Valid() {
		t.Error("Valid() returns true and it should be false when calling Check()")
	}
}

func TestForm_ErrorGet(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")

	// if it cant find an error, it returns an empty string
	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("should have an error returned from Get, but do not")
	}

	s = form.Errors.Get("whatever")

	if len(s) != 0 {
		t.Error("should not have an error, but got one")
	}
}
