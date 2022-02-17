package api

import (
	"log"
	"net/http"
	"time"

	"github.com/GrayFinance/mint/src/config"
	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter().StrictSlash(true)

	routerAPI := router.PathPrefix("/api").Subrouter()
	routerAPI.Path("/user").HandlerFunc(IsAuthorized(GetUser)).Methods("GET")
	routerAPI.Path("/wallets").HandlerFunc(IsAuthorized(ListWallets)).Methods("GET")
	routerAPI.Path("/wallet/create").HandlerFunc(IsAuthorized(CreateWallet)).Methods("POST")
	routerAPI.Path("/wallet/{wallet_id}").HandlerFunc(WalletMiddleware(GetWallet)).Methods("GET")

	userRouter := routerAPI.PathPrefix("/user").Subrouter()
	userRouter.Path("/create").HandlerFunc(CreateUser).Methods("POST")
	userRouter.Path("/auth").HandlerFunc(AuthUser).Methods("POST")
	userRouter.Path("/change/password").HandlerFunc(IsAuthorized(ChangePassword)).Methods("PUT")

	walletRouter := routerAPI.PathPrefix("/wallet/{wallet_id}").Subrouter()
	walletRouter.Path("/delete").HandlerFunc(IsAuthorized(DeleteWallet)).Methods("DELETE")
	walletRouter.Path("/rename").HandlerFunc(IsAuthorized(RenameWallet)).Methods("PUT")
	walletRouter.Path("/receive").Queries("network", "{network}").HandlerFunc(WalletMiddleware(Receive)).Methods("GET")
	walletRouter.Path("/transfer").HandlerFunc(WalletMiddleware(Transfer)).Methods("POST")
	walletRouter.Path("/withdraw/{network}").HandlerFunc(WalletMiddleware(Withdraw)).Methods("POST")
	walletRouter.Path("/payinvoice").HandlerFunc(WalletMiddleware(PayInvoice)).Methods("POST")
	walletRouter.Path("/payments").Queries("offset", "{offset}").HandlerFunc(WalletMiddleware(ListPayments)).Methods("GET")
	walletRouter.Path("/payment/{hash_id}").HandlerFunc(WalletMiddleware(GetPayment)).Methods("GET")

	server := &http.Server{
		Addr:         config.Config.API_HOST + ":" + config.Config.API_PORT,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	log.Println("Server Listen API: ", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
