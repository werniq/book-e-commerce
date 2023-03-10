package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	StringMap            map[string]string
	IntMap               map[string]int
	FloatMap             map[string]float32
	Data                 map[string]interface{}
	IsAuthenticated      int
	MainnetAddress       string
	ErrorData            []string
	StripeSecretKey      string
	StripePublishableKey string
	CSRFToken            string
	Flash                string
	Warning              string
	Error                string
	CSSVersion           string
}

//go:embed admin/templates
var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.MainnetAddress = "0x59ABDFCc700DfB6fFf671B2198B26107f6AFE036"
	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("admin/templates/%s.gohtml", page)

	_, templateInMap := app.templateCache[templateToRender]

	if app.cfg.env == "production" && templateInMap {
		t = app.templateCache[templateToRender]
	} else {
		t, err = app.parseTemplate(page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = app.addDefaultData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	t, err = template.New(fmt.Sprintf("%s.gohtml", page)).ParseFS(templateFS, "admin/templates/base.gohtml", templateToRender)

	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t
	return t, nil
}
