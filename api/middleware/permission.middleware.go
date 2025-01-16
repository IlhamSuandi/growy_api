package middleware

import (
	"net/http"
	"slices"
	"strings"

	permission "github.com/ilhamSuandi/business_assistant/constant"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/repository"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/utils"
	"gorm.io/gorm"
)

// hasPermission checks if the user has any of the required permissions
func hasPermission(userPermissions []string, requiredPermissions []string) bool {
	for _, required := range requiredPermissions {
		for _, userPermission := range userPermissions {
			if userPermission == required || userPermission == "*" {
				return true
			}
		}
	}
	return false
}

func Permission(requiredPermissions []string, db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(requiredPermissions) == 0 || slices.Contains(requiredPermissions, "all") {
			next.ServeHTTP(w, r)
			return
		}

		log := utils.Log
		userInfo := r.Context().Value("userInfo").(*model.User)
		userRepo := repository.NewUserRepository(db)
		requestMethod := r.Method

		// admin can access all resources
		if userInfo.Role == "admin" {
			next.ServeHTTP(w, r)
			return
		}

		// admin restricted routes
		if slices.Contains(requiredPermissions, permission.Admin) {
			log.Error("[middleware] not enough permission")
			response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
				Message: "you don't have permission to access this resource",
				Error:   "you don't have permission to access this resource",
				Status:  http.StatusForbidden,
			})
			return
		}

		// owner can access all resources except admin restricted routes
		if userInfo.Role == "owner" {
			next.ServeHTTP(w, r)
			return
		}

		// getting user permissions
		userPermission, err := userRepo.GetUserPermission(userInfo.Id)
		if err != nil {
			log.Errorf("[middleware] error getting user permission %s", err)

			response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
				Message: "error getting user permission",
				Error:   err.Error(),
				Status:  http.StatusForbidden,
			})
			return
		}

		// verify if user has required permission
		for _, permission := range *userPermission {
			permissionActions := strings.Split(permission.Action, ",")

			hasRequiredPermission := hasPermission([]string{permission.Resource}, requiredPermissions)

			isAllowedAction := slices.Contains(permissionActions, "all") || slices.Contains(permissionActions, strings.ToLower(requestMethod))

			if hasRequiredPermission && isAllowedAction {
				// log.Info("[middleware] permission is valid")
				next.ServeHTTP(w, r)
				return
			}
		}

		log.Error("[middleware] not enough permission")
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "you don't have permission to access this resource",
			Error:   "you don't have permission to access this resource",
			Status:  http.StatusForbidden,
		})
	})
}
