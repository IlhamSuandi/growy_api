package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ilhamSuandi/business_assistant/api/dto"
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/database/model"
	"github.com/ilhamSuandi/business_assistant/pkg/auth"
	"github.com/ilhamSuandi/business_assistant/pkg/auth/oauth"
	"github.com/ilhamSuandi/business_assistant/pkg/response"
	"github.com/ilhamSuandi/business_assistant/types"
	"github.com/ilhamSuandi/business_assistant/usecase"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	AuthUsecase usecase.AuthUsecase
	Logger      *logrus.Logger
}

func NewAuthController(authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
		Logger:      utils.Log,
	}
}

// @Tags Auth
// @Summary Register as user
// @Description Register new user
// @Router /auth/register [post]
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Request body"
// @Failure 400 {object} types.ErrorResponse "request body is invalid"
// @Failure 409 {object} types.ErrorResponse "user exists"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Success 201 {object} types.Response{data=dto.RegisteredUserResponse} "Successfully registered user"
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	ac.Logger.Info("parsing request body")
	var registerPayload dto.RegisterRequest
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.Header.Get("X-Forwarded-For")

	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// get the request body
	if err := utils.ParseJSON(r, &registerPayload); err != nil {
		ac.Logger.Errorf("error parsing request body %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error Parsing Request Body",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	ac.Logger.Info("checking if user exists")
	if _, exists := ac.AuthUsecase.IsUserExists(registerPayload.Email); exists {
		ac.Logger.Error("user already exists")
		response.WriteError(w, http.StatusConflict, types.ErrorResponse{
			Message: "Error Registering User",
			Error:   "User Already Exists",
			Status:  http.StatusConflict,
		})
		return
	}

	ac.Logger.Info("creating user")
	newUser := model.User{
		Email:        registerPayload.Email,
		Username:     registerPayload.Username,
		Password:     registerPayload.Password,
		AuthProvider: "custom",
	}

	if err := ac.AuthUsecase.CreateUser(&newUser); err != nil {
		ac.Logger.Errorf("error creating user %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Registering User",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Generating Access Token and Refresh Token
	ac.Logger.Info("generating access token")
	sessionId := uuid.New()
	accessTokenExpiration := time.Minute * time.Duration(config.JWT_EXPIRATION)
	accessToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    newUser.UUID,
		SessionId: sessionId,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Duration:  accessTokenExpiration,
	})
	if err != nil {
		ac.Logger.Errorf("error creating access token %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ac.Logger.Info("generating refresh token")
	refreshTokenExpiration := time.Hour * time.Duration(config.JWT_REFRESH_EXPIRATION)
	refreshToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    newUser.UUID,
		SessionId: sessionId,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Duration:  refreshTokenExpiration,
	})
	if err != nil {
		ac.Logger.Errorf("error creating refresh token %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Storing Refresh Token in Database
	ac.Logger.Info("storing refresh token in database")
	session := model.Session{
		Model: model.Model{
			UUID: sessionId,
		},
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(refreshTokenExpiration),
		IsRevoked:    false,
		Users:        []*model.User{&newUser},
	}

	if err = ac.AuthUsecase.CreateSession(session); err != nil {
		ac.Logger.Errorf("error storing refresh token in database %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Storing Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Setting HttpOnly Cookie to Client
	ac.Logger.Info("setting http only cookie to client")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenExpiration),
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   config.APP_ENV == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// Return users
	ac.Logger.Info("successfully registered user")
	response.WriteJSON(w, http.StatusCreated, types.Response{
		Message: "Successfully Register User",
		Data: dto.RegisteredUserResponse{
			Email:    newUser.Email,
			Username: newUser.Username,
			Token: dto.TokenResponse{
				TokenType:   "Bearer",
				AccessToken: accessToken,
				ExpiresIn:   accessTokenExpiration.String(),
			},
		},
		Status: http.StatusCreated,
	})
}

// @Tags Auth
// @Summary login as user
// @Description login as user
// @Router /auth/login [post]
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Request body"
// @Failure 400 {object} types.ErrorResponse "request body is invalid"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Success 200 {object} types.Response{data=dto.TokenResponse} "Successfully Logged In"
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	ac.Logger.Info("parsing request body")
	var loginPayload dto.LoginRequest
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.Header.Get("X-Forwarded-For")

	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// get the request body
	if err := utils.ParseJSON(r, &loginPayload); err != nil {
		ac.Logger.Errorf("error parsing request body %s", err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error Parsing Request Body",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Check if user exists
	ac.Logger.Info("checking if user exists")
	user, ok := ac.AuthUsecase.IsUserExists(loginPayload.Email)
	if !ok {
		ac.Logger.Error("user does not exist")
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "email or password is incorrect",
			Error:   "email or password is incorrect",
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Comparing request password with stored user password
	ac.Logger.Info("comparing request password with stored user password")
	if ok := auth.ComparePassword(loginPayload.Password, user.Password); !ok {
		ac.Logger.Errorf("email or password is incorrect")
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "email or password is incorrect",
			Error:   "email or password is incorrect",
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Generating Access Token and Refresh Token
	ac.Logger.Info("generating access token")
	sessionId := uuid.New()
	accessTokenExpiration := time.Minute * time.Duration(config.JWT_EXPIRATION)
	accessToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    user.UUID,
		SessionId: sessionId,
		Username:  user.Username,
		Email:     user.Email,
		Duration:  accessTokenExpiration,
	})
	if err != nil {
		ac.Logger.Errorf("error creating access token %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ac.Logger.Info("generating refresh token")
	refreshTokenExpiration := time.Hour * time.Duration(config.JWT_REFRESH_EXPIRATION)
	refreshToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    user.UUID,
		SessionId: sessionId,
		Username:  user.Username,
		Email:     user.Email,
		Duration:  refreshTokenExpiration,
	})
	if err != nil {
		ac.Logger.Errorf("error creating refresh token %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Storing Refresh Token in Database
	ac.Logger.Info("storing refresh token in database")
	session := model.Session{
		Model: model.Model{
			UUID: sessionId,
		},
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(refreshTokenExpiration),
		IsRevoked:    false,
		Users:        []*model.User{&user},
	}

	if err = ac.AuthUsecase.CreateSession(session); err != nil {
		ac.Logger.Errorf("error storing refresh token in database %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Storing Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Setting HttpOnly Cookie to Client
	ac.Logger.Info("setting http only cookie to client")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenExpiration),
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   config.APP_ENV == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	// Return Successfully Login With Its Data
	ac.Logger.Info("successfully logged in")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Login",
		Data: dto.TokenResponse{
			TokenType:   "Bearer",
			AccessToken: accessToken,
			ExpiresIn:   accessTokenExpiration.String(),
		},
		Status: http.StatusOK,
	})
}

// @Tags Auth
// @Summary Renew Access Token
/* @Description This endpoint allows the user to renew their access token using a refresh token stored in a cookie.
 * but sadly swagger ui can't do http only cookie authentication
 */
// @Accept json
// @Produce json
// @Router /auth/renew [post]
// @Success 200 {object} types.Response{data=dto.TokenResponse} "Successfully Renewed Access Token"
// @Failure 401 {object} types.ErrorResponse "Unauthorized: Invalid or expired refresh token"
// @Failure 303 {object} types.ErrorResponse "Redirect: User must log in again"
func (ac *AuthController) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token from cookie
	ac.Logger.Info("Renewing Access Token")
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		ac.Logger.Errorf("Error Parsing Token %s", err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// verify old access token
	ac.Logger.Info("Verifying Refresh Token")
	claims, err := auth.ParseToken(refreshToken.Value)
	if err != nil {
		ac.Logger.Errorf("Error Parsing Token %s", err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// check if user have refresh token in database based on user informations
	ac.Logger.Info("Checking User Session")
	session, err := ac.AuthUsecase.GetUserSession(claims.SessionId)
	if err != nil {
		ac.Logger.Errorf("Error Getting User Session %s", err)
		// redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if session.IsRevoked {
		// redirect to login page
		ac.Logger.Error("User Session is Revoked")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// renew access token
	ac.Logger.Info("creating new token")
	newAccessTokenExpiration := time.Minute * time.Duration(config.JWT_EXPIRATION)
	newAccessToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    claims.UserId,
		SessionId: claims.SessionId,
		Username:  claims.Username,
		Email:     claims.Email,
		Duration:  newAccessTokenExpiration,
	})

	// return new access token
	ac.Logger.Info("successfully renewed token")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Renew Access Token",
		Data: dto.TokenResponse{
			TokenType:   "Bearer",
			AccessToken: newAccessToken,
			ExpiresIn:   newAccessTokenExpiration.String(),
		},
		Status: http.StatusOK,
	})
}

// @Tags Auth
// @Summary Logout User
// @Description this endpoint used for logout user and remove user refresh token and accesstoken
// @Accept json
// @Produce json
// @Router /auth/logout [post]
// @Success 200 {object} types.Response{data=dto.TokenResponse} "Successfully Renewed Access Token"
// @Failure 401 {object} types.ErrorResponse "Unauthorized: Invalid or expired refresh token"
// @Failure 303 {object} types.ErrorResponse "Redirect: User must log in again"
func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Getting Claims From Auth Middleware Context
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// verify old access token
	claims, err := auth.ParseToken(refreshToken.Value)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Find User Session
	session, err := ac.AuthUsecase.GetUserSession(claims.SessionId)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error getting session",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Delete Existed User Session
	if err := ac.AuthUsecase.DeleteUserSession(session.UUID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error deleting session",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	// Remove HttpOnly Cookie from Client
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Return Successfully Logged Out
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Logged Out",
		Data:    "Successfully Logged Out",
		Status:  http.StatusOK,
	})
}

func (ac *AuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	url := config.GoogleOauthConfig.AuthCodeURL(state)

	http.SetCookie(w, &http.Cookie{
		Name:   "oauth2_state",
		Value:  state,
		MaxAge: 60,
	})

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (ac *AuthController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	cookieState, err := r.Cookie("oauth2_state")

	if cookieState.Value != state && err != nil {
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "Error",
			Error:   "code is required",
			Status:  http.StatusForbidden,
		})
		return
	}

	// Get Google Token
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "Error Exchange Token",
			Error:   err,
			Status:  http.StatusForbidden,
		})
		return
	}

	userInfo, err := oauth.GetUserInfo(token.AccessToken)

	// create user if not exist else update user verified email
	existingUser, exist := ac.AuthUsecase.IsUserExists(userInfo.Email)
	newUser := model.User{
		Email:           userInfo.Email,
		Username:        userInfo.Name,
		IsEmailVerified: true,
		AuthProvider:    "google",
	}

	if exist && !existingUser.IsEmailVerified {
		if err := ac.AuthUsecase.UpdateUser(&newUser); err != nil {
			response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
				Message: "Error Updating User",
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			})
			return
		}
	} else if !exist {
		if err := ac.AuthUsecase.CreateUser(&newUser); err != nil {
			response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
				Message: "Error Creating User",
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			})
			return
		}
	}

	accessTokenExpiration := time.Minute * time.Duration(config.JWT_EXPIRATION)
	accessToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    newUser.UUID,
		SessionId: uuid.MustParse(state),
		Username:  newUser.Username,
		Email:     newUser.Email,
		Duration:  accessTokenExpiration,
	})
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	refreshTokenExpiration := time.Hour * time.Duration(config.JWT_REFRESH_EXPIRATION)
	refreshToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    newUser.UUID,
		SessionId: uuid.MustParse(state),
		Username:  newUser.Username,
		Email:     newUser.Email,
		Duration:  refreshTokenExpiration,
	})
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Setting HttpOnly Cookie to Client
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTokenExpiration),
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   config.APP_ENV == "production",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Logged In",
		Data: dto.TokenResponse{
			TokenType:   "Bearer",
			AccessToken: accessToken,
			ExpiresIn:   accessTokenExpiration.String(),
		},
		Status: http.StatusOK,
	})
}
