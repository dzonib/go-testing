package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

// we declared variable usable by test files, test files are ignored when we build application
// they are only used when we run tests

var app application

// this function will be executed before tests run
func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	// now we can use all db methods
	app.DB = &dbrepo.TestDBRepo{}

	// this runs all tests
	os.Exit(m.Run())
}
