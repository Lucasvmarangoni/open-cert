package handlers

import (
	"certifications/internal/middlewares"
	"certifications/internal/modules/users/http/dto"
	"certifications/pkg/validate"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/jwtauth"

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
func (u *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.UserInput
	requestID := middlewares.GetRequestID(r.Context())

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[http.StatusBadRequest],
			"message":    fmt.Sprintf("%v", err),
			"request_id": requestID,
		})
		return
	}
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		return
	}

	err = u.userService.Create(user.Name, user.LastName, user.CPF, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msgf("error create user. | (%s)", requestID)
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
func (u *UsersHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	var user dto.AuthenticationInput
	requestID := middlewares.GetRequestID(r.Context())

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     "bad_request",
			"message":    "Invalid JSON format",
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Msgf("error decode request. | (%s)", requestID)
		return
	}

	err = validate.ValidatePassword(user.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		return
	}

	authdata, err := u.userService.Authn(user.Email, user.Password)
	if err != nil {
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msgf("error authenticate user. | (%s)", requestID)
		return
	}

	if authdata.IsTemporaryToken {
		response := dto.GetJWTOutput{
			AccessToken:      authdata.Token,
			TwoAuthnRequirid: true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.GetJWTOutput{
		AccessToken:      authdata.Token,
		TwoAuthnRequirid: authdata.IsTemporaryToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Authentication godoc
// @Summary      Generate 2FA Secret
// @Description  Generate 2FA Secret.
// @Tags         Authn
// @Produce      json
// @Success      200  {object}  dto.OTPOutput
// @Failure 	 500 {object} map[string]string "Error response"
// @Router       /totp/generate [get]
func (u *UsersHandler) TwoFactorAuthn(w http.ResponseWriter, r *http.Request) {
	var response dto.OTPOutput
	requestID := middlewares.GetRequestID(r.Context())

	id, err := u.GetSub(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		return
	}

	otpResponse, err := u.userService.GenerateTOTP(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Msgf("error generate totp. | (%s)", requestID)
		return
	}
	response = dto.OTPOutput{
		Base32:      otpResponse.Base32,
		OtpauthUrl:  otpResponse.Otpauth_url,
		AccessToken: otpResponse.AccessToken,
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
func (u *UsersHandler) TwoFactorVerify(w http.ResponseWriter, r *http.Request) {
	var totpToken dto.OTPInput
	requestID := middlewares.GetRequestID(r.Context())

	id, err := u.GetSub(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&totpToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[http.StatusBadRequest],
			"message":    fmt.Sprintf("%v", err),
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Msgf("error decode request. | (%s)", requestID)
		return
	}

	authdata, err := u.userService.VerifyTOTP(id, totpToken.Token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Msgf("error validate totp. | (%s)", requestID)
		return
	}
	response := dto.GetJWTOutput{
		AccessToken:      authdata.Token,
		TwoAuthnRequirid: authdata.IsTemporaryToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Authentication godoc
// @Summary      Disable 2FA
// @Description Disable 2FA.
// @Tags         Authn
// @Produce      json
// @Success      200  {object}  map[string]bool "otp_disabled"
// @Failure 	 500 {object} map[string]string "Error response"
// @Router       /totp/disable [patch]
func (u *UsersHandler) TwoFactorDisable(w http.ResponseWriter, r *http.Request) {
	requestID := middlewares.GetRequestID(r.Context())

	id, err := u.GetSub(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		return
	}

	err = u.userService.DisableOTP(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Msgf("error disable two factor authentication. | (%s)", requestID)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"otp_disabled": true})
}

func (u *UsersHandler) GetSub(w http.ResponseWriter, r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", errs.Wrap(err, "jwtauth.FromContext", http.StatusInternalServerError)
	}
	id, ok := claims["sub"].(string)
	if !ok {
		return "", errs.Wrap(err, `claims["sub"].(string)`, http.StatusInternalServerError)
	}
	return id, nil
}
