package auth

import (
	"fmt"
	"net/http"

	mutations "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/mutation"
	"github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/queries"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type updatePasswordBody struct {
	// Token string `json:"token" graphql:"token"`
	OldPassword string `json:"old_password" graphql:"old_password"`
	NewPassword string `json:"new_password" graphql:"new_password"`
}

func UpdatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		var body struct {
			Input struct {
				Params updatePasswordBody `json:"params"`
			}`json:"input"`
			SessionVariables map[string]interface{} `json:"session_variables"`
		}
		if err = ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "unprocessed_request", Success: false})
			return
		}
		userId := body.SessionVariables["x-hasura-user-id"]
		filters := []string{
			fmt.Sprintf(`id:{_eq: "%s"}`, userId ),
		}
		user, err := queries.FindUser(filters)
		if err != nil {
			logrus.Error("find user error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		err = CheckUserValidity(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: err.Error(), Success: false})
			return
		}

		// Check if the old password is correct ...
		isMatch := utils.CompareHashedPassword(user.Password, body.Input.Params.OldPassword)
		if !isMatch {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "incorrect_old_password", Success: false})
			return
		}
		// Hash newly provided password
		hashedNewPassword , err := utils.HashPassword(body.Input.Params.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "process_halted", Success: false})
			return
		}
		// Exec update user mutation
		err = mutations.UpdateUserPassword(hashedNewPassword, user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{Message: "update_mutation_error", Success: false})
			return
		}


		ctx.JSON(http.StatusOK, model.Response{Message: "password_changed_successfully", Success: true})
	}
}