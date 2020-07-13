package apiexch

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"user-apiv2/model"
	"user-apiv2/tools"

	"github.com/go-chi/chi"
)

// ExchHandler represents an HTTP API handler for exchs.
type ExchHandler struct {
	mux *chi.Mux
	//redirectURL string
	session *Session
	//exch        *Exch
	Logger *log.Logger
}

// NewExchHandler returns a new instance of ExchHandler.
//mv ExchService
func NewExchHandler(dbExchUser UDBType, dbExchData ExDBType, msess model.Sessioner) *ExchHandler {
	h := &ExchHandler{
		mux:     chi.NewRouter(),
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		session: NewSession(dbExchUser, dbExchData, msess),
	}
	h.mux.Get("/exch/account", h.session.session.Validate(h.AccountHandler))
	h.mux.Get("/exch/market", h.session.session.Validate(h.MarketHandler))
	h.mux.Get("/exch/wallet", h.session.session.Validate(h.WalletHandler))
	// h.mux.Get("/exch/support", h.session.session.Validate(h.SupportHandler))
	return h
}

func (h *ExchHandler) AccountHandler(w http.ResponseWriter, r *http.Request) {
	_, user, err := h.session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	dat := "" //struct {
	// 	Letter  string
	// 	Exch_id ExchID
	// 	UserGui *ExchGuiService
	// 	//UserGui *apipage.PageGuiService
	// }{
	// 	Letter:  Gletter,
	// 	Exch_id: exchID,
	// 	UserGui: h.session.ExchGuiService,
	// 	//UserGui: h.session.PageGuiService,
	// }
	if err := tools.AccountTmpl.Execute(w, r, dat, user, nil); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *ExchHandler) MarketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\n\n********************* In Market *************************  \n\n\n\n")
	_, user, err := h.session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login?redirect=/exch/market", http.StatusFound)
		return
	}
	if err := tools.MarketTmpl.Execute(w, r, nil, user, nil); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

type walletData struct {
	UsrAddress string
}

func (h *ExchHandler) WalletHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\n\n********************* In WalletHandler *************************  \n\n\n\n")
	_, user, err := h.session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login?redirect=/exch/market", http.StatusFound)
		return
	}

	d := walletData{
		UsrAddress: "15PiHsSMMVrV5bpe28uHrY1E25pYTWWBWu", //h.session.Userservice.GetUserAddress(context.Background(), user),
	}

	log.Println("\n********************* In WalletHandler ***d.UsrAddress**********************  ", d.UsrAddress)

	if err := tools.WalletTmpl.Execute(w, r, d, user, nil); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return

}

// func (h *ExchHandler) SupportHandler(w http.ResponseWriter, r *http.Request) {
// 	user, err := h.session.UserFromRequest(r)
// 	if err != nil {
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 	}
// 	ctx := context.WithValue(context.Background(), AddExchKey, user)
// 	exch, err := h.exchFromForm(w, r)
// 	if err != nil {
// 		tools.Error(w, errors.Wrapf(err, "could not save exch: %v", err), http.StatusInternalServerError, h.Logger)
// 		return
// 	}
// 	exch.JustStarted = "false"
// 	_, err = h.session.Exchservice.AddExch(ctx, exch)

// 	dat := struct {
// 		Exch    *Exch
// 		UserGui *ExchGuiService
// 	}{
// 		Exch:    exch,
// 		UserGui: h.session.ExchGuiService,
// 	}

// 	if err := tools.ExchTmpl.Execute(w, r, dat, user, false); err != nil {
// 		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
// 		return
// 	}
// 	return
// }

func (h *ExchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// // handleGetExch handles requests to create a new exch.

// // exchFromRequest retrieves a exch from the database given a exch ID in the
// func (h *ExchHandler) exchFromRequest(r *http.Request) (*Exch, error) {
// 	ids := chi.URLParam(r, "id")
// 	id, err := strconv.Atoi(ids)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "bad exch id: %v", err)
// 	}
// 	exch, err := h.session.Exchservice.GetExch(ExchID(id))
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "could not find exch: %v", err)
// 	}
// 	return exch, nil
// }

// func (h *ExchHandler) exchFromForm(w http.ResponseWriter, r *http.Request) (*Exch, error) {
// 	ids, _ := strconv.Atoi(r.FormValue("id"))
// 	wgs := []byte(r.FormValue("wronguessesr"))
// 	exch := &Exch{
// 		PlayerID:       r.FormValue("playerid"),
// 		GuessedLetterR: r.FormValue("guessedletterr"),
// 		WordSoFarR:     r.FormValue("wordsofarr"),
// 		WrongGuessesR:  wgs,
// 		Invalid:        r.FormValue("invalid"),
// 		ID:             ExchID(ids),
// 		JustStarted:    r.FormValue("juststarted"),
// 	}
// 	fmt.Printf("\n\n\n In FromForm exch =  %v \n\n\n\n", *exch)
// 	return exch, nil
// }
