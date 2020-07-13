package apicrypto

import (
	"log"
	"net/http"
	"os"
	"user-apiv2/model"
	"user-apiv2/tools"

	"github.com/go-chi/chi"
)

// TransactionHandler represents an HTTP API handler for exchs.
type TransactionHandler struct {
	mux *chi.Mux
	//redirectURL string
	session *Session
	//exch        *Transaction
	Logger *log.Logger
}
// NewTransactionHandler returns a new instance of TransactionHandler.
//mv TransactionService
func NewTransactionHandler(TransactionUserDB UserDBType, TransactionDataDB TransDBType, msess model.Sessioner) *TransactionHandler {
	h := &TransactionHandler{
		mux:     chi.NewRouter(),
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		session: NewSession(TransactionUserDB, TransactionDataDB, msess),
	}
	h.mux.Get("/crypto/makeprofit", h.session.session.Validate(h.MakeprofitHandler))
	//h.mux.Get("/crypto/market", h.session.session.Validate(h.MarketHandler))
	//h.mux.Get("/crypto/wallet", h.session.session.Validate(h.WalletHandler))
	// h.mux.Get("/exch/support", h.session.session.Validate(h.SupportHandler))
	return h
}

		//UserApiKey: "44c07067d435d717f6e968a319b8b783",
		//UserApiSecret: "539faf2ff2ba6ce1ffefeccb72b340e1",
		//yHost: "https://api.hitbtc.com",

func (h *TransactionHandler) MakeprofitHandler(w http.ResponseWriter, r *http.Request) {
	_, user, err := h.session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	td := *Transaction{
		SymbolID: r.FormValue("symbolid"),
		SoldOrBoughtPrice: r.FormValue("soldorboughtprice"),
		ProfitTargetRate: r.FormValue("profittargetrate"),
		QtyMultiplier: r.FormValue("qtymultiplier"),
		TransType: r.FormValue("transtype"),
		Side: r.FormValue("side"),		
	}

	td := model.Transaction{
		SymbolID: "XMOETH",
		SoldOrBoughtPrice: 0.0,
		ProfitTargetRate: 0.02,
		QtyMultiplier: 2.0,
		TransType: "A",
		Side: "buy",
		UserApiKey: "44c07067d435d717f6e968a319b8b783",
		UserApiSecret: "539faf2ff2ba6ce1ffefeccb72b340e1",
		Host: "https://api.hitbtc.com",
	}

	var worker model.TransactionService

	if strings.Contains(td.Host, "hitbtc"){
		cl := hitbtc.NewClient(td.Host, td.UserApiKey, td.UserApiSecret)
		worker = hitbtc.NewWorker (cl)
	}else if strings.Contains(td.Host, "binance"){
		cl := binance.NewClient(td.Host, td.UserApiKey, td.UserApiSecret)
		worker = binance.NewWorker (cl)
	}
	aT := bot.NewAutoTrade(worker)
	aT.AutoTrade(td)

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


func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}