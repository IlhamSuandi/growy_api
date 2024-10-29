package middleware

import (
	"net/http"
	"slices"
)

// hasPermission checks if the user has any of the required permissions
func hasPermission(userPermissions []string, requiredPermissions []string) bool {
	for _, required := range requiredPermissions {
		for _, userPermission := range userPermissions {
			if userPermission == required {
				return true // User has the required permission
			}
		}
	}
	return false // No matching permission found
}

func Permission(requiredPermissions []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims := r.Context().Value("claims").(types.JwtClaims)

		if requiredPermissions == nil || slices.Contains(requiredPermissions, "all") {
			next.ServeHTTP(w, r)
		} else if hasPermission([]string{"admin"}, requiredPermissions) {
			next.ServeHTTP(w, r)
		}
	})
}
