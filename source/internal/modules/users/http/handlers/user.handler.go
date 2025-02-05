package handlers

import (
	"certifications/internal/middlewares"
	"certifications/internal/modules/users/http/dto"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/go-chi/jwtauth"
	"github.com/rs/zerolog/log"
)

// Get user godoc
// @Summary      Get me user data
// @Description  Get me user data
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 		 {object}  dto.UserOutput
// @Failure      500
// @Router       /user/me [get]
// @Security ApiKeyAuth
func (u *UsersHandler) Me(w http.ResponseWriter, r *http.Request) {
	requestID := middlewares.GetRequestID(r.Context())

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[http.StatusInternalServerError],
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Msgf("failed to get JWT claims. | (%s)", requestID)
		return
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[http.StatusInternalServerError],
			"request_id": requestID,
		})
		log.Error().Msgf("sub claim is missing or not a string. | (%s)", requestID)
		return
	}
	finduser, err := u.userService.FindById(sub, nil)
	if err != nil {
		w.WriteHeader(errs.Assertion(err).Code)
		log.Error().Err(errs.Assertion(err).Stack()).Stack().Msgf("user not found. | (%s)", requestID)
		return
	}
	user := dto.UserOutput{
		ID:        finduser.ID,
		Name:      finduser.Name,
		LastName:  finduser.LastName,
		CPF:       finduser.CPF,
		Email:     finduser.Email,
		Certifier: finduser.Certifier,
	}
	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[http.StatusInternalServerError],
			"request_id": requestID,
		})
		log.Error().Msgf("error to marshal indent user data. | (%s)", requestID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

// Update user godoc
// @Summary      Update user
// @Description  Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.UserInput  true  "user data update, password is required"
// @Success      200
// @Failure      400
// @Router       /user/update [put]
// @Security ApiKeyAuth
func (u *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user dto.UserUpdateInput
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

	id, err := u.GetSub(w, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Stack().Msgf("faleiled get sub. | (%s)", requestID)
		return
	}

	err = u.userService.Update(id, user.Password, user.NewEmail, user.NewPassword)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Err(errs.Assertion(err).Stack()).Stack().Msgf("failed to Update user. | (%s)", requestID)
		return
	}
	log.Info().Str("context", "UsersHandler").Msgf("user updated successfully (%s)", id)
	w.WriteHeader(http.StatusOK)
}

// Delete user godoc
// @Summary      Delete user
// @Description  Delete user
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Router       /user/del [delete]
// @Security ApiKeyAuth
func (u *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		log.Error().Err(errs.Assertion(err).Stack()).Stack().Msgf("faleiled get sub. | (%s)", requestID)
		return
	}
	err = u.userService.Delete(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     errs.Status[errs.Assertion(err).Code],
			"message":    fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msg("error delete user ")
		return
	}
	log.Info().Str("context", "UsersHandler").Msgf("user deleted successfully (%s)", id)
	w.WriteHeader(http.StatusOK)
}
