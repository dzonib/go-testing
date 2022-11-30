package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
	"webapp/pkg/data"
)

// path that would be in production
var pathToTemplates = "./templates/"

// W in our case web browser
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// template data
	var td = make(map[string]any)

	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "Hit this page at "+time.Now().UTC().String())
	}

	_ = app.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
}

type TemplateData struct {
	IP    string
	Data  map[string]any
	Error string
	Flash string
	User  data.User
}

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {

	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}

// W in our case web browser
func (app *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	// parse template from disc
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))

	// template not found, or error in template
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(r.Context())

	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	// execute the template, passing the date, if any
	err = parsedTemplate.Execute(w, td)

	if err != nil {
		return err
	}
	return nil
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {

	// A stub is a small piece of code that takes the place of another component during testing.
	//  The benefit of using a stub is that it returns consistent results,
	//   making the test easier to write.
	// And you can run tests even if the other components are not working yet
	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// validate data
	form := NewForm(r.PostForm)

	form.Required("email", "password")

	if !form.Valid() {
		// redirect to login page with error message
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)

	if err != nil {
		// redirect to login page with error message
		app.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// if not authenticated then redirect with error
	if !app.authenticate(r, user, password) {
		app.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// prevent fixation attack
	_ = app.Session.RenewToken(r.Context())

	// redirect to some other page
	app.Session.Put(r.Context(), "flash", "Successfully logged in")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *application) authenticate(r *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}

	app.Session.Put(r.Context(), "user", user)
	return true
}
