package api

import (
	"log"
	"net/http"
	"time"

	"github.com/GrayFinance/mint/src/config"
	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Path("/user").HandlerFunc(IsAuthorized(GetUser)).Methods("GET")
	router.Path("/wallets").HandlerFunc(IsAuthorized(ListWallets)).Methods("GET")
	router.Path("/wallet/{wallet_id}").HandlerFunc(WalletMiddleware(GetWallet)).Methods("GET")

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Path("/create").HandlerFunc(CreateUser).Methods("POST")
	userRouter.Path("/auth").HandlerFunc(AuthUser).Methods("POST")

	walletRouter := router.PathPrefix("/wallet").Subrouter()
	walletRouter.Path("/create").HandlerFunc(IsAuthorized(CreateWallet)).Methods("POST")
	walletRouter.Path("/{wallet_id}/delete").HandlerFunc(IsAuthorized(DeleteWallet)).Methods("DELETE")
	walletRouter.Path("/{wallet_id}/rename").HandlerFunc(IsAuthorized(RenameWallet)).Methods("PUT")
	walletRouter.Path("/{wallet_id}/receive").Queries("network", "{network}").HandlerFunc(WalletMiddleware(Receive)).Methods("GET")

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
