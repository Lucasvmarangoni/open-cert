package handlers

import (
	"certifications/internal/modules/users/http/dto"
	"certifications/pkg/validate"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errs "github.com/Lucasvmarangoni/logella/err"
	"github.com/asaskevich/govalidator"

	// "github.com/go-chi/chi"
	// "github.com/go-chi/jwtauth"
	"github.com/rs/zerolog/log"
)

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.UserInput  true  "user data"
// @Success      200
// @Failure      400
// @Router       /authn/create [post]
func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.UserInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "BadRequest",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	err = validate.ValidatePassword(user.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "BadRequest",
			"message": fmt.Sprintf("Valid password is required. %v", err),
		})
		return
	}

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "BadRequest",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	err = u.userService.Create(user.Name, user.LastName, user.CPF, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "BadRequest",
			"message": fmt.Sprintf("%v", err),
		})
		log.Error().Stack().Err(errs.Stack()).Msg("Error create user")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Authentication godoc
// @Summary      Generate a user JWT
// @Description  Generate a user JWT. Requires either a CPF or an Email and Password.
// @Tags         Authn
// @Accept       json
// @Produce      json
// @Param        request     body      dto.AuthenticationInput  true  "Authentication input. Requires either a CPF or an Email and Password."
// @Success      200  {object}  dto.GetJWTOutput
// @Failure      400  {object}  string  "Both email and CPF are required for authentication."
// @Failure      401  {object}  string  "Unauthorized."
// @Router       /authn/jwt [post]
// func (u *UserHandler) Authentication(w http.ResponseWriter, r *http.Request) {
// 	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
// 	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)
// 	var user dto.AuthenticationInput

// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		log.Error().Err(errs.Stack()).Msg("Error decode request")
// 		return
// 	}

// 	err = validate.ValidatePassword(user.Password)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"status":  "BadRequest",
// 			"message": fmt.Sprintf("Valid password is required. %v", err),
// 		})
// 		return
// 	}

// 	err = u.validateUserUpdateInputForCPFAndEmail(&user)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"status":  "BadRequest",
// 			"message": fmt.Sprintf("%v", err),
// 		})
// 		return
// 	}

// 	unique := user.Email + user.CPF
// 	tokenString, err := u.userService.Authn(unique, user.Password, jwt, jwtExpiresIn)
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		log.Error().Stack().Err(errs.Stack()).Msg("Error authenticate user")
// 		return
// 	}

// 	accessToken := dto.GetJWTOutput{AccessToken: tokenString}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(accessToken)
// }

/*
// Authentication godoc
// @Summary      Generate 2FA Secret
// @Description  Generate 2FA Secret.
// @Tags         Authn
// @Produce      json
// @Success      200  {object}  dto.OTPOutput
// @Failure 	 500 {object} map[string]string "Error response"
// @Router       /totp/generate [get]
func (u *UserHandler) TwoFactorAuthn(w http.ResponseWriter, r *http.Request) {
	var response dto.OTPOutput

	id, err := u.GetSub(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "StatusInternalServerError",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	otpResponse, err := u.userService.GenerateTOTP(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "StatusInternalServerError",
			"message": fmt.Sprintf("%v", err),
		})
		log.Error().Err(errs.Stack()).Msg("Error generate totp")
		return
	}

	response = dto.OTPOutput{
		Base32:     otpResponse.Base32,
		OtpauthUrl: otpResponse.Otpauth_url,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Authentication godoc
// @Summary      Verify 2FA
// @Description Verify 2FA. The isValidate parameter should be "1" for the first validation attempt. For subsequent attempts, any value or an empty string is accepted.
// @Tags         Authn
// @Produce      json
// @Param        request     body      dto.OTPInput  true  "Authentication input. Requires a token and isValidate."
// @Success      200  {object}  dto.OTPOutput
// @Failure 	 500 {object} map[string]string "Error response"
// @Failure 	 400
// @Router       /totp/verify/{is_validate} [post]
func (u *UserHandler) TwoFactorVerify(w http.ResponseWriter, r *http.Request) {
	var totpToken dto.OTPInput

	id, err := u.GetSub(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "StatusInternalServerError",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&totpToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Err(errs.Stack()).Msg("Error decode request")
		return
	}

	isValidate := chi.URLParam(r, "is_validate")
	if isValidate == "" {
		isValidate = "false"
	}
	err = u.userService.VerifyTOTP(id, totpToken.Token, isValidate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "StatusInternalServerError",
			"message": fmt.Sprintf("%v", err),
		})
		log.Error().Err(errs.Stack()).Msg("Error validate totp")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"otp_verified": true})
}

// Authentication godoc
// @Summary      Disable 2FA
// @Description Disable 2FA.
// @Tags         Authn
// @Produce      json
// @Success      200  {object}  map[string]bool "otp_disabled"
// @Failure 	 500 {object} map[string]string "Error response"
// @Router       /totp/disable [patch]
func (u *UserHandler) TwoFactorDisable(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetSub(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "StatusInternalServerError",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	err = u.userService.DisableOTP(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "StatusInternalServerError",
			"message": fmt.Sprintf("%v", err),
		})
		log.Error().Err(errs.Stack()).Msg("Error disable two factor authentication")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"otp_disabled": true})
}

func (u *UserHandler) GetSub(w http.ResponseWriter, r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", errors.ErrCtx(err, "jwtauth.FromContext")
	}
	id, ok := claims["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return "", errors.ErrCtx(err, `claims["sub"].(string)`)
	}
	return id, nil
}
*/
func (u *UserHandler) validateUserUpdateInputForCPFAndEmail(user *dto.AuthenticationInput) error {

	if user.Email == "" && user.CPF == "" {
		return errors.New("an Email or a CPF is necessary")
	}
	if user.Email != "" && user.CPF != "" {
		user.CPF = ""
	}
	if err := u.validateEmail(&user.Email); err != nil {
		return err
	}
	if err := u.validateCPF(&user.CPF); err != nil {
		return err
	}
	return nil
}

