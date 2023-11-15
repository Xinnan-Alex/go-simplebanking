package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Xinnan-Alex/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD MYR SGD"`
}

func (server *Server) checkCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorReponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return false
	}

	return true
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	if !server.checkCurrency(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.checkCurrency(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
