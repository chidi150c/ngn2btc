package tools

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

var (
	IndexTmpl     = newappTemplate("index.html")
	ListusersTmpl = newappTemplate("listusers.html")
	// SupportTmpl   = newappTemplate("support.html")
	WalletTmpl  = newappTemplate("wallet.html")
	MarketTmpl  = newappTemplate("market.html")
	AccountTmpl = newappTemplate("account.html")
	//RateChan    = make(chan string)
)

// appTemplate is a user login-aware wrapper for a html/template.
type appTemplate struct {
	t *template.Template
}

// parseTemplate applies a given file to the body of the base template.
func newappTemplate(filename string) *appTemplate {
	var tmpl *template.Template
	path := strings.Join([]string{"tools/templates", filename}, "/")

	if filename == "index.html" {
		b, err := ioutil.ReadFile("tools/templates/base.html")
		if err != nil {
			panic(fmt.Errorf("could not read template: %v", err))
		}
		ba := strings.Replace(string(b), "{{template \"sidebar\" .User}}", "", 1)
		ba = strings.Replace(string(ba), "w3-hide-small\" id=\"forIndex\"", "\" id=\"forIndexA\"", 3)
		tmpl = template.Must(template.New("base").Parse(ba))
		tmpl.ParseFiles("tools/templates/chat.html", "tools/templates/index.html")
	} else if filename == "market.html" {
		tmpl = template.Must(template.ParseFiles("tools/templates/base.html", "tools/templates/chat.html", "tools/templates/sidebar.html", path))
	} else {
		tmpl = template.Must(template.ParseFiles("tools/templates/base.html", "tools/templates/chat.html", "tools/templates/sidebar.html", path))
	}

	return &appTemplate{t: tmpl}
}

func (tmpl *appTemplate) appTemplate() *template.Template {
	return tmpl.t
}

// Execute writes the template using the provided data, adding login and user
// information to the base template...   usr interface{}, noFooter bool
func (tmpl *appTemplate) Execute(w http.ResponseWriter, r *http.Request, dat interface{}, usr interface{}, msg interface{}) error {
	d := struct {
		Data        interface{}
		AuthEnabled bool
		LoginURL    string
		LogoutURL   string
		//AddFooter   bool
		SignupURL string
		User      interface{}
		Msg       interface{}
		//Rate      string
	}{
		Data:        dat,
		AuthEnabled: true,
		LoginURL:    "/login?redirect=" + r.URL.RequestURI(),
		LogoutURL:   "/logout?redirect=" + r.URL.RequestURI(),
		SignupURL:   "/signup?redirect=" + r.URL.RequestURI(),
		//AddFooter:   noFooter,
		User: usr,
		Msg:  msg,
		//Rate: <-RateChan,
	}
	if err := tmpl.t.Execute(w, d); err != nil {
		return errors.Wrapf(err, "could not write template: %+v", err)
	}
	return nil
}

// func (tmpl *appTemplate) parseTmpl(filename string) *appTemplate {
// 	var tm *template.Template
// 	path := strings.Join([]string{"tools/templates", filename}, "/")
// 	tm = template.Must(tmpl.t.ParseFiles(path))
// 	return &appTemplate{t: tm}
// }

// func baseWithOutSidebar(filename string) *appTemplate {
// 	var tm *template.Template
// 	path := strings.Join([]string{"tools/templates", filename}, "/")
// 	b, err := ioutil.ReadFile("tools/templates/base.html")
// 	if err != nil {
// 		panic(fmt.Errorf("could not read template: %v", err))
// 	}
// 	ba := strings.Replace(string(b), "{{template \"sidebar\" .User}}", "", 1)
// 	tm = template.Must(template.New("base").Parse(ba))
// 	tm.ParseFiles("tools/templates/chat.html", path)
// 	return &appTemplate{t: tm}
// }
// func BaseWithSidebarAndChart(filename string) *appTemplate {
// 	var tmpl *template.Template
// 	path := strings.Join([]string{"tools/templates", filename}, "/")
// 	tmpl = template.Must(template.ParseFiles("tools/templates/base.html", "tools/templates/sidebar.html", "tools/templates/chat.html", path))
// 	return &appTemplate{t: tmpl}
// }
// func (tmpl *appTemplate) AddChart(graph string) *appTemplate {
// 	var tm *template.Template
// 	path := "tools/templates/market.html"
// 	b, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		panic(fmt.Errorf("could not read chart template: %v", err))
// 	}
// 	ba := strings.Replace(string(b), "<!-- ParseGraphHere -->", graph, 1)
// 	tm = template.Must(tmpl.t.Lookup("base.html").Parse(ba))
// 	return &appTemplate{t: tm}
// }
