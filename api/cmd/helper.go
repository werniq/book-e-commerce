package main

import "net/http"

func (app *application) HelperParseForm(r *http.Request, formName string) bool {
	r.ParseForm()

	if r.FormValue(formName) != "" {
		return true
	}
	return false
}
