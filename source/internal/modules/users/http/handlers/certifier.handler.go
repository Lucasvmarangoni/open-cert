package handlers

import (
	"certifications/internal/middlewares"
	"certifications/internal/modules/users/http/dto"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/asaskevich/govalidator"
	"github.com/rs/zerolog/log"
)

func (u *UsersHandler) CreateCertifier(w http.ResponseWriter, r *http.Request) {
	var certifier dto.CertifierInput
	requestID := middlewares.GetRequestID(r.Context())

	err := json.NewDecoder(r.Body).Decode(&certifier)
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

	sub, err := u.GetSub(w, r)
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

	_, err = govalidator.ValidateStruct(certifier)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  errs.Status[errs.Assertion(err).Code],
			"message": fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msgf("error validate certifier. | (%s)", requestID)
		return
	}

	err = u.certifierService.CreateCertifier(sub, certifier.Description)
	if err != nil {
		w.WriteHeader(errs.Assertion(err).Code)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  errs.Status[errs.Assertion(err).Code],
			"message": fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msgf("error create certifier. | (%s)", requestID)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (u *UsersHandler) DocIdentity(w http.ResponseWriter, r *http.Request) {
	requestID := middlewares.GetRequestID(r.Context())

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  errs.Status[http.StatusBadRequest],
			"message": fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msgf("error processing form. | (%s)", requestID)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  errs.Status[http.StatusBadRequest],
			"message": fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
			"request_id": requestID,
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msgf("error to getting file. | (%s)", requestID)
		return
	}
	defer file.Close()

	/*
		Verificação de identidade.
			- veriff
	*/
}

func (u *UsersHandler) FaceIdentity(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  errs.Status[http.StatusBadRequest],
			"message": fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msg("error processing form")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "BadRequest",
			"message": fmt.Sprintf("%v", errs.Assertion(err).ToClient()),
		})
		log.Error().Stack().Err(errs.Assertion(err).Stack()).Msg("error to getting file")
		return
	}
	defer file.Close()

	/*
		Verificação de identidade.
			- veriff
	*/
}
