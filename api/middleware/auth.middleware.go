package middleware

import (
	"context"
	"net/http"

	"github.com/ilhamSuandi/business_assistant/pkg/auth"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

func Auth(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := utils.Log
		ctx := context.Background()
		userRepo := repository.NewUserRepository(db)

		log.Info("getting access token")
		accessToken, err := auth.GetAccessToken(r)
		if err != nil {
			log.Errorf("error getting access token %s", err)
			response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
				Message: "Unauthorized",
				Error:   "Access Token is required",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		log.Info("parsing access token")
		claims, err := auth.ParseToken(accessToken)
		if err != nil {
			log.Errorf("error parsing access token %s", err)
			response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
				Message: "Token is invalid",
				Error:   err.Error(),
				Status:  http.StatusUnauthorized,
			})
			return
		}

		log.Info("getting user informations")
		user, err := userRepo.GetUserByUserId(claims.UserId)
		if err != nil {
			log.Errorf("error getting user informations %s", err)
			response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
				Message: "Unauthorized",
				Error:   err.Error(),
				Status:  http.StatusUnauthorized,
			})
			return
		}

		log.Info("setting context")
		ctx = context.WithValue(ctx, "accessToken", accessToken)
		ctx = context.WithValue(ctx, "claims", claims)
		ctx = context.WithValue(ctx, "userInfo", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
