package api

import (
	db "MelBank/db/sqlc"
	"MelBank/token"
	"MelBank/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type createManagerRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
}

func (server *Server) createManager(ctx *gin.Context) {
	var req createManagerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.AddManagerParams{
		FullName:       req.FullName,
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}
	man, err := server.store.AddManager(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, man)
}

func (server *Server) listAllRequests(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if !authPayload.IsManager {
		ctx.JSON(http.StatusForbidden, "You are not the manager")
		return
	}
	list, err := server.store.ListRequests(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, list)
}

type checkRequestStruct struct {
	ReqId int `json:"req_id"`
	ManId int `json:"man_id"`
}

func (server *Server) checkRequest(ctx *gin.Context) {
	var req checkRequestStruct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if !authPayload.IsManager {
		ctx.JSON(http.StatusForbidden, "You are not the manager")
		return
	}
	request, err := server.store.GetRequest(ctx, int64(req.ReqId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := struct {
		Currency string `json:"currency"`
		Username string `json:"username"`
		FullName string `json:"full_Name"`
	}{}
	err = json.Unmarshal([]byte(request.Req), &arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	secRes := util.SecurityCheck()
	finRes := util.FinReport()
	arg1 := db.AddDecisionParams{}

	if secRes && finRes > 0.5 {
		arg1 = db.AddDecisionParams{
			ManID:    int64(req.ManId),
			Decision: true,
			Message:  sql.NullString{},
		}
		_, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
			Owner:    arg.Username,
			Balance:  0,
			Currency: arg.Currency,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "foreign_key_violation", "unique_violation":
					ctx.JSON(http.StatusForbidden, errorResponse(err))
					return
				}
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		arg1 = db.AddDecisionParams{
			ManID:    int64(req.ManId),
			Decision: false,
			Message: sql.NullString{
				String: fmt.Sprintf(`Request denied 'cause user has low fin. rating %d 
					or didn't pass security check'`, finRes),
			},
		}
	}
	_, err = server.store.AddDecision(ctx, arg1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = server.store.DeleteFromQueue(ctx, int64(req.ReqId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, "Decision have been made")

}

type loginManagerRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginManagerResponse struct {
	AccessToken string
	Manager     db.Manager
}

func (server *Server) loginManager(ctx *gin.Context) {
	var req loginManagerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	man, err := server.store.GetManager(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, man.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(man.Username, true, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := loginManagerResponse{
		AccessToken: accessToken,
		Manager:     man,
	}
	ctx.JSON(http.StatusOK, rsp)
}
