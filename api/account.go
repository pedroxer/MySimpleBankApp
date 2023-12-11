package api

import (
	db "MelBank/db/sqlc"
	"MelBank/token"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
	Username string `json:"username"`
	FullName string `json:"full_Name"`
	Balance  int    `json:"balance"`
}

// createAccount godoc
// @Summary createAccount
// @Security ApiKeyAuth
// @Tags account
// @Description create account
// @Accept json
// @Produce json
// @Param input body createAccountRequest true "account info"
// @Success 200 {string} string "Request send. Please come back later"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /accounts [post]
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	req.Username = authPayload.Username
	req.FullName = user.FullName
	accReq, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.AddRequestToQueue(ctx, accReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Request send. Please come back later")
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount godoc
// @Summary getAccount
// @Security ApiKeyAuth
// @Tags account
// @Description get account
// @Accept json
// @Produce json
// @Param input body getAccountRequest true "account info"
// @Success 200 {object} db.Account
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 401 {string} string "account doesn't belong to the user"
// @Failure 500 {object} error
// @Router /accounts/:id [get]
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccount godoc
// @Summary listAccount
// @Security ApiKeyAuth
// @Tags account
// @Description list accounts
// @Accept json
// @Produce json
// @Param input body listAccountRequest true "account info"
// @Success 200 {object} []db.Account
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 401 {string} string "account doesn't belong to the user"
// @Failure 500 {object} error
// @Router /accounts [get]
func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	account, err := server.store.ListAccounts(ctx, arg)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)
}

type updateAccountRequest struct {
	ID     int64 `form:"id" binding:"required,min=1"`
	Amount int64 `form:"amount" binding:"required,min=1"`
}

// updateAccount godoc
// @Summary updateAccount
// @Security ApiKeyAuth
// @Tags account
// @Description update balance of account
// @Accept json
// @Produce json
// @Param input body updateAccountRequest true "account info"
// @Success 200 {object} db.Account
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 401 {string} string "account doesn't belong to the user"
// @Failure 500 {object} error
// @Router /accounts/update [put]
func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Amount,
	}
	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)

}
