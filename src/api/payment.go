package api

import (
	"net/http"
	"strconv"

	"github.com/GrayFinance/mint/src/services"
	"github.com/GrayFinance/mint/src/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func ListPayments(w http.ResponseWriter, r *http.Request) {
	permission := context.Get(r, "permission").(string)
	if permission == "admin" || permission == "read" {
		query, _ := r.URL.Query()["offset"]
		payment := services.Payment{
			UserID:   context.Get(r, "user_id").(string),
			WalletID: mux.Vars(r)["wallet_id"],
		}

		offset, _ := strconv.Atoi(query[0])
		list_transactions, err := payment.ListPayments(offset)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}
		utils.SendJSONResponse(w, list_transactions)
		return
	} else {
		utils.SendJSONError(w, 500, "The key permission is not read.")
		return
	}
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	permission := context.Get(r, "permission").(string)
	if permission == "admin" || permission == "read" {
		payment := services.Payment{
			UserID:   context.Get(r, "user_id").(string),
			WalletID: mux.Vars(r)["wallet_id"],
		}

		data, err := payment.GetPayment(mux.Vars(r)["hash_id"])
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}
		utils.SendJSONResponse(w, data)
		return
	} else {
		utils.SendJSONError(w, 500, "The key permission is not read.")
		return
	}
}
