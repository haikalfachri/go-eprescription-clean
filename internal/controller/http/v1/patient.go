package v1

import (
	"go-eprescription-clean/internal/controller/http/v1/request"
	"go-eprescription-clean/internal/entity"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Summary     Create Patient
// @Description Create a new patient
// @ID          create-patient
// @Tags        patient
// @Accept      json
// @Produce     json
// @Param       request body request.Patient true "Patient input"
// @Success     201 {object} response.SuccessPatient
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /patients [post]
func (r *V1) createPatient(ctx *fiber.Ctx) error {
	var req request.Patient

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - createPatient")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - createPatient")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, fiber.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	patient, err := r.u.Patient.Create(ctx.UserContext(), entity.Patient{
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - createPatient")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to create patient")
	}

	return successResponse(ctx, http.StatusCreated, "patient created successfully", patient)
}

// @Summary     Get all patients
// @Description Retrieve all patient records
// @ID          get-all-patients
// @Tags        patient
// @Accept      json
// @Produce     json
// @Success 	200 {object} response.SuccessPatientList
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /patients [get]
func (r *V1) getAllPatients(ctx *fiber.Ctx) error {
	patients, err := r.u.Patient.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllPatients")
		if strings.Contains(err.Error(), "no patients found") {
			return errorResponse(ctx, http.StatusNotFound, "no patients found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to retrieve patients")
	}

	return successResponse(ctx, http.StatusOK, "patients retrieved successfully", patients)
}

// @Summary     Get patient by ID
// @Description Retrieve a single patient by ID	
// @ID          get-patient-by-id
// @Tags        patient
// @Accept      json
// @Produce     json
// @Param       id path string true "Patient ID"
// @Success 	200 {object} response.SuccessPatient
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /patients/{id} [get]
func (r *V1) getByIDPatient(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		r.l.Error("http - v1 - getByIDPatient")
		return errorResponse(ctx, http.StatusBadRequest, "patient ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - getByIDPatient")
		return errorResponse(ctx, http.StatusBadRequest, "invalid patient ID")
	}

	patient, err := r.u.Patient.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getByIDPatient")
		if strings.Contains(err.Error(), "patient not found") {
			return errorResponse(ctx, http.StatusNotFound, "patient not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to retrieve patient")
	}

	return successResponse(ctx, http.StatusOK, "patient retrieved successfully", patient)
}


// @Summary     Update Patient
// @Description Update an existing patient
// @ID          update-patient
// @Tags        patient
// @Accept      json
// @Produce     json
// @Param       id path string true "Patient ID"
// @Param       request body request.Patient true "Patient input"
// @Success 	200 {object} response.SuccessPatient
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /patients/{id} [patch]
func (r *V1) updatePatient(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		r.l.Error("http - v1 - updatePatient")
		return errorResponse(ctx, http.StatusBadRequest, "patient ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - updatePatient")
		return errorResponse(ctx, http.StatusBadRequest, "invalid patient ID")
	}

	var req request.Patient

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - updatePatient")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - updatePatient")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, fiber.StatusBadRequest, formatValidationErrors(ve))
		} 
	}

	patient, err := r.u.Patient.Update(ctx.UserContext(), id, entity.Patient{
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - updatePatient")
		if strings.Contains(err.Error(), "patient not found") {
			return errorResponse(ctx, http.StatusNotFound, "patient not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update patient")
	}

	return successResponse(ctx, http.StatusOK, "patient updated successfully", patient)
}

// @Summary     Delete Patient
// @Description Delete a patient by ID
// @ID          delete-patient
// @Tags        patient
// @Accept      json
// @Produce     json
// @Param       id path string true "Patient ID"
// @Success 	200 {object} response.SuccessString
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /patients/{id} [delete]
func (r *V1) deletePatient(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "patient ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid patient ID")
	}

	err := r.u.Patient.Delete(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - deletePatient")
		if strings.Contains(err.Error(), "patient not found") {
			return errorResponse(ctx, http.StatusNotFound, "patient not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete patient")
	}

	return successResponse(ctx, http.StatusOK, "patient deleted successfully", "deleted")
}


