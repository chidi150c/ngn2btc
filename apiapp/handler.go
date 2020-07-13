package apiapp

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"user-apiv2/apichat"
	"user-apiv2/apiexch"
	"user-apiv2/apiuser"
	"user-apiv2/apicrypto"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	UserHandler *apiuser.UserHandler
	ExchHandler *apiexch.ExchHandler
	ChatHandler *apichat.ChatHandler
	//ChartHandler *chart.ChartHandler
	TransactionHandler *apicrypto.TransactionHandler
}

//initializies the Handler struct a *chart.ChartHandler
func NewHandler(u *apiuser.UserHandler, g *apiexch.ExchHandler, c *apichat.ChatHandler) *Handler {
	return &Handler{
		UserHandler: u,
		ExchHandler: g,
		ChatHandler: c,
		//ChartHandler: a,
	}
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("*********Started Handler.ServeHTTP********")
	defer log.Println("*********End Handler.ServeHTTP********")
	fmt.Println(r.URL.Path)
	if strings.HasPrefix(r.URL.Path, "/tools/asset/") {
		fmt.Println()
		http.StripPrefix("/tools/asset/", http.FileServer(http.Dir("./tools/asset/"))).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/crypto") {
		h.TransactionHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/exch") {
		h.ExchHandler.ServeHTTP(w, r)
	}else if strings.HasPrefix(r.URL.Path, "/chat") {
		h.ChatHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/") {
		log.Println("*******************In AppHandler2*******************")
		fmt.Println()
		h.UserHandler.ServeHTTP(w, r)
	} else {
		fmt.Println()
		http.NotFound(w, r)
	}
}

// else if strings.HasPrefix(r.URL.Path, "/chart") {
// 	h.ChartHandler.ServeHTTP(w, r)
// }
