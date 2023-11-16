package api

import (
	"net/http"
	"time"

	db "github.com/Xinnan-Alex/simplebank/db/sqlc"
	"github.com/Xinnan-Alex/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUserAPI(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorReponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}

	response := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, response)
}
