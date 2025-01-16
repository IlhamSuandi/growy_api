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
	ac.Logger.Infof("[%s /auth/register] parsing request body", r.Method)
	var registerPayload dto.RegisterRequest
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.Header.Get("X-Forwarded-For")

	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// get the request body
	if err := utils.ParseJSON(r, &registerPayload); err != nil {
		ac.Logger.Errorf("[%s /auth/register] error parsing request body %s", r.Method, err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error Parsing Request Body",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	ac.Logger.Infof("[%s /auth/register] checking if user exists", r.Method)
	if _, exists := ac.AuthUsecase.IsUserExists(registerPayload.Email); exists {
		ac.Logger.Errorf("[%s /auth/register] user already exists", r.Method)
		response.WriteError(w, http.StatusConflict, types.ErrorResponse{
			Message: "Error Registering User",
			Error:   "User Already Exists",
			Status:  http.StatusConflict,
		})
		return
	}

	ac.Logger.Infof("[%s /auth/register] creating user", r.Method)
	newUser := model.User{
		Email:        registerPayload.Email,
		Username:     registerPayload.Username,
		Password:     registerPayload.Password,
		AuthProvider: "custom",
	}

	if err := ac.AuthUsecase.CreateUser(&newUser); err != nil {
		ac.Logger.Errorf("[%s /auth/register] error creating user %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Registering User",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Generating Access Token and Refresh Token
	ac.Logger.Infof("[%s /auth/register] generating access token", r.Method)
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
		ac.Logger.Errorf("[%s /auth/register] error creating access token %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ac.Logger.Infof("[%s /auth/register] generating refresh token", r.Method)
	refreshTokenExpiration := time.Hour * time.Duration(config.JWT_REFRESH_EXPIRATION)
	refreshToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    newUser.UUID,
		SessionId: sessionId,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Duration:  refreshTokenExpiration,
	})
	if err != nil {
		ac.Logger.Errorf("[%s /auth/register] error creating refresh token %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Storing Refresh Token in Database
	ac.Logger.Infof("[%s /auth/register] storing refresh token in database", r.Method)
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
		ac.Logger.Errorf("[%s /auth/register] error storing refresh token in database %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Storing Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Setting HttpOnly Cookie to Client
	ac.Logger.Infof("[%s /auth/register] setting http only cookie to client", r.Method)
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
	ac.Logger.Infof("[%s /auth/register] successfully registered user", r.Method)
	response.WriteJSON(w, http.StatusCreated, types.Response{
		Message: "Successfully Register User",
		Data: dto.RegisteredUserResponse{
			Email:    newUser.Email,
			Username: newUser.Username,
			Token: dto.TokenResponse{
				TokenType:   "Bearer",
				AccessToken: accessToken,
				ExpiresIn:   int64(accessTokenExpiration.Milliseconds()),
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
	ac.Logger.Infof("[%s /auth/login] parsing request body", r.Method)
	var loginPayload dto.LoginRequest
	userAgent := r.Header.Get("[/auth/login] User-Agent")
	ipAddress := r.Header.Get("X-Forwarded-For")

	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// get the request body
	if err := utils.ParseJSON(r, &loginPayload); err != nil {
		ac.Logger.Errorf("[%s /auth/login] error parsing request body %s", r.Method, err)
		response.WriteError(w, http.StatusBadRequest, types.ErrorResponse{
			Message: "Error Parsing Request Body",
			Error:   err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Check if user exists
	ac.Logger.Infof("[%s /auth/login] checking if user exists", r.Method)
	user, ok := ac.AuthUsecase.IsUserExists(loginPayload.Email)
	if !ok {
		ac.Logger.Errorf("[%s /auth/login] user does not exist", r.Method)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "email or password is incorrect",
			Error:   "email or password is incorrect",
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Comparing request password with stored user password
	ac.Logger.Infof("[%s /auth/login] comparing request password with stored user password", r.Method)
	if ok := auth.ComparePassword(loginPayload.Password, user.Password); !ok {
		ac.Logger.Errorf("[%s /auth/login] email or password is incorrect", r.Method)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "email or password is incorrect",
			Error:   "email or password is incorrect",
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Generating Access Token and Refresh Token
	ac.Logger.Infof("[%s /auth/login] generating access token", r.Method)
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
		ac.Logger.Errorf("[%s /auth/login] error creating access token %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ac.Logger.Infof("[%s /auth/login] generating refresh token", r.Method)
	refreshTokenExpiration := time.Hour * time.Duration(config.JWT_REFRESH_EXPIRATION)
	refreshToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    user.UUID,
		SessionId: sessionId,
		Username:  user.Username,
		Email:     user.Email,
		Duration:  refreshTokenExpiration,
	})
	if err != nil {
		ac.Logger.Errorf("[%s /auth/login] error creating refresh token %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Storing Refresh Token in Database
	ac.Logger.Infof("[%s /auth/login] storing refresh token in database", r.Method)
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
		ac.Logger.Errorf("[%s /auth/login] error storing refresh token in database %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Storing Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Setting HttpOnly Cookie to Client
	ac.Logger.Infof("[%s /auth/login] setting http only cookie to client", r.Method)
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
	ac.Logger.Infof("[%s /auth/login] successfully logged in", r.Method)
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Login",
		Data: dto.TokenResponse{
			TokenType:   "Bearer",
			AccessToken: accessToken,
			ExpiresIn:   int64(accessTokenExpiration.Milliseconds()),
		},
		Status: http.StatusOK,
	})
}

// @Tags Auth
// @Summary Renew Access Token
/*  @Description This endpoint allows the user to renew their access token using a refresh token stored in a cookie.
 *  but sadly swagger ui can't do http only cookie authentication
 */
// @Produce json
// @Router /auth/token/renew [get]
// @Success 200 {object} types.Response{data=dto.TokenResponse} "Successfully Renewed Access Token"
// @Failure 401 {object} types.ErrorResponse "Unauthorized: Invalid or expired refresh token"
// @Failure 303 {object} types.ErrorResponse "Redirect: User must log in again"
func (ac *AuthController) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token from cookie
	ac.Logger.Infof("[%s /auth/token/renew] Renewing Access Token", r.Method)
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		ac.Logger.Errorf("[%s /auth/token/renew] Error Parsing Token %s", r.Method, err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// verify old access token
	ac.Logger.Infof("[%s /auth/token/renew] Verifying Refresh Token", r.Method)
	claims, err := auth.ParseToken(refreshToken.Value)
	if err != nil {
		ac.Logger.Errorf("[%s /auth/token/renew] Error Parsing Token %s", r.Method, err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// check if user have refresh token in database based on user informations
	ac.Logger.Infof("[%s /auth/token/renew] Checking User Session", r.Method)
	session, err := ac.AuthUsecase.GetUserSession(claims.SessionId)
	if err != nil {
		ac.Logger.Errorf("[%s /auth/token/renew] Error Getting User Session %s", r.Method, err)
		// redirect to login page
		http.Redirect(w, r, "/api/v1/auth/login", http.StatusSeeOther)
		return
	}

	if session.IsRevoked {
		// redirect to login page
		ac.Logger.Error("User Session is Revoked")
		http.Redirect(w, r, "/api/v1/auth/login", http.StatusSeeOther)
		return
	}

	// renew access token
	ac.Logger.Infof("[%s /auth/token/renew] creating new token", r.Method)
	newAccessTokenExpiration := time.Minute * time.Duration(config.JWT_EXPIRATION)
	newAccessToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    claims.UserId,
		SessionId: claims.SessionId,
		Username:  claims.Username,
		Email:     claims.Email,
		Duration:  newAccessTokenExpiration,
	})
	if err != nil {
		ac.Logger.Errorf("[%s /auth/token/renew] error creating access token %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// return new access token
	ac.Logger.Infof("[%s /auth/token/renew] successfully renewed token", r.Method)
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Renew Access Token",
		Data: dto.TokenResponse{
			TokenType:   "Bearer",
			AccessToken: newAccessToken,
			ExpiresIn:   int64(newAccessTokenExpiration.Milliseconds()),
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
	ac.Logger.Infof("[%s /auth/logout] getting refresh token from cookie", r.Method)
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// verify old refresh token
	ac.Logger.Infof("[%s /auth/logout] parsing refresh token", r.Method)
	claims, err := auth.ParseToken(refreshToken.Value)
	if err != nil {
		ac.Logger.Errorf("[%s /auth/logout] error parsing refresh token %s", r.Method, err)
		response.WriteError(w, http.StatusUnauthorized, types.ErrorResponse{
			Message: "Error Parsing Token",
			Error:   err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	// Find User Session
	ac.Logger.Infof("[%s /auth/logout] getting user session", r.Method)
	session, err := ac.AuthUsecase.GetUserSession(claims.SessionId)
	if err != nil {
		ac.Logger.Errorf("[%s /auth/logout] error getting session %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error getting session",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Delete Existed User Session
	ac.Logger.Infof("[%s /auth/logout] deleting user session", r.Method)
	if err := ac.AuthUsecase.DeleteUserSession(session.UUID); err != nil {
		ac.Logger.Errorf("[%s /auth/logout] error deleting session %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "error deleting session",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	// Remove HttpOnly Cookie from Client
	ac.Logger.Infof("[%s /auth/logout] removing refresh token from cookie", r.Method)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Return Successfully Logged Out
	ac.Logger.Infof("[%s /auth/logout] successfully logged out", r.Method)
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Logged Out",
		Data:    "Successfully Logged Out",
		Status:  http.StatusOK,
	})
}

// @Tags Auth
// @Summary Redirects user to Google's OAuth 2.0 authentication page
// @Description Initiates Google OAuth 2.0 login by generating a state and redirecting to Google's authorization URL
// @Accept json
// @Produce json
// @Success 303 {string} string "Redirect to Google OAuth 2.0 login page"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/google/login [get]
func (ac *AuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	url := config.GoogleOauthConfig.AuthCodeURL(state)

	utils.Log.Info("[/auth/google/login] setting oauth cookie")
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth2_state",
		Value:  state,
		MaxAge: 60,
	})

	http.Redirect(w, r, url, http.StatusSeeOther)
}

// @Tags Auth
// @Summary Google OAuth callback handler
// @Description Handles the Google OAuth callback, exchanges the code for a token, and creates/updates a user.
// @Accept  json
// @Produce  json
// @Param state query string true "State for OAuth validation"
// @Param code query string true "Authorization code received from Google"
// @Success 201 {object} types.Response{data=dto.TokenResponse} "Successfully Logged In"
// @Failure 400 {object} types.ErrorResponse "Bad Request, Invalid Code or State"
// @Failure 403 {object} types.ErrorResponse "Forbidden, Invalid token exchange or user info"
// @Failure 500 {object} types.ErrorResponse "Internal Server Error"
// @Router /auth/google/callback [get]
func (ac *AuthController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("[/auth/google/callback] getting state and code")

	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.Header.Get("X-Forwarded-For")

	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	cookieState, err := r.Cookie("oauth2_state")

	utils.Log.Info("[/auth/google/callback] checking if state is equal")
	if cookieState.Value != state && err != nil {
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "Error",
			Error:   "code is required",
			Status:  http.StatusForbidden,
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "oauth2_state",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	// Get Google Token
	utils.Log.Info("[/auth/google/callback] exchanging code")
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.Log.Errorf("[/auth/google/callback] error exchanging code %s", err)
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "Error Exchange Token",
			Error:   err,
			Status:  http.StatusForbidden,
		})
		return
	}

	utils.Log.Info("[/auth/google/callback] getting user info")
	userInfo, err := oauth.GetUserInfo(token.AccessToken)
	if err != nil {
		utils.Log.Errorf("[/auth/google/callback] error getting user info %s", err)
		response.WriteError(w, http.StatusForbidden, types.ErrorResponse{
			Message: "Error Getting User Info",
			Error:   err.Error(),
			Status:  http.StatusForbidden,
		})
		return
	}

	// create user if not exist else update user verified email
	utils.Log.Info("[/auth/google/callback] checking if user exists")
	existingUser, exist := ac.AuthUsecase.IsUserExists(userInfo.Email)
	newUser := model.User{
		Email:           userInfo.Email,
		Username:        userInfo.Name,
		IsEmailVerified: true,
		AuthProvider:    "google",
	}

	var user model.User

	utils.Log.Info("[/auth/google/callback] checking if user is verified")
	if exist && !existingUser.IsEmailVerified {
		utils.Log.Info("[/auth/google/callback] user exist updating user")
		updatedUser, err := ac.AuthUsecase.UpdateUser(existingUser.Id, newUser)
		if err != nil {
			utils.Log.Errorf("[/auth/google/callback] error updating user %s", err)
			response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
				Message: "Error Updating User",
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			})
			return
		}
		user = *updatedUser
	} else if !exist {
		utils.Log.Info("[/auth/google/callback] creating user")
		if err := ac.AuthUsecase.CreateUser(&newUser); err != nil {
			utils.Log.Errorf("[/auth/google/callback] error creating user %s", err)
			response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
				Message: "Error Creating User",
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			})
			return
		}
		user = newUser
	} else {
		utils.Log.Info("[/auth/google/callback] user is exist")
		user = existingUser
	}

	utils.Log.Info("[/auth/google/callback] creating access token")
	accessTokenExpiration := time.Minute * time.Duration(config.JWT_EXPIRATION)
	accessToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    user.UUID,
		SessionId: uuid.MustParse(state),
		Username:  user.Username,
		Email:     user.Email,
		Duration:  accessTokenExpiration,
	})
	if err != nil {
		utils.Log.Errorf("[/auth/google/callback] error creating access token %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	utils.Log.Info("[/auth/google/callback] creating refresh token")
	refreshTokenExpiration := time.Hour * time.Duration(config.JWT_REFRESH_EXPIRATION)
	refreshToken, _, err := ac.AuthUsecase.CreateToken(types.CreateToken{
		UserId:    user.UUID,
		SessionId: uuid.MustParse(state),
		Username:  user.Username,
		Email:     user.Email,
		Duration:  refreshTokenExpiration,
	})
	if err != nil {
		utils.Log.Errorf("[/auth/google/callback] error creating refresh token %s", err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Creating Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	ac.Logger.Infof("[%s /auth/google/callback] storing refresh token in database", r.Method)
	session := model.Session{
		Model: model.Model{
			UUID: uuid.MustParse(state),
		},
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(refreshTokenExpiration),
		IsRevoked:    false,
		Users:        []*model.User{&user},
	}

	// Storing Refresh Token in Database
	ac.Logger.Infof("[%s /auth/google/callback] storing refresh token in database", r.Method)
	if err = ac.AuthUsecase.CreateSession(session); err != nil {
		ac.Logger.Errorf("[%s /auth/google/callback] error storing refresh token in database %s", r.Method, err)
		response.WriteError(w, http.StatusInternalServerError, types.ErrorResponse{
			Message: "Error Storing Token",
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}

	utils.Log.Info("[/auth/google/callback] setting http only cookie to client")
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

	utils.Log.Info("[/auth/google/callback] successfully logged in")
	response.WriteJSON(w, http.StatusOK, types.Response{
		Message: "Successfully Logged In",
		Data: dto.TokenResponse{
			TokenType:   "Bearer",
			AccessToken: accessToken,
			ExpiresIn:   int64(accessTokenExpiration.Milliseconds()),
		},
		Status: http.StatusOK,
	})
}
