package v1

import (
	"net/http"
	"strings"

	"go-eprescription-clean/internal/controller/http/v1/request"
	"go-eprescription-clean/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)



// @Summary     Create Medicine
// @Description Create a new medicine
// @ID          create-medicine
// @Tags        medicine
// @Accept      json
// @Produce     json
// @Param       request body request.Medicine true "Medicine input"
// @Success 	201 {object} response.SuccessMedicine
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicines [post]
func (r *V1) createMedicine(ctx *fiber.Ctx) error {
	var req request.Medicine

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - createMedicine")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - createMedicine")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, fiber.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	medicine, err := r.u.Medicine.Create(ctx.UserContext(), entity.Medicine{
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - createMedicine")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to create medicine")
	}

	return successResponse(ctx, http.StatusCreated, "medicine created successfully", medicine)
}

// @Summary     Get all medicines
// @Description Retrieve all medicine records
// @ID          get-all-medicines
// @Tags        medicine
// @Accept      json
// @Produce     json
// @Success 	200 {object} response.SuccessMedicineList
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicines [get]
func (r *V1) getAllMedicines(ctx *fiber.Ctx) error {
	medicines, err := r.u.Medicine.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllMedicines")
		if strings.Contains(err.Error(), "no medicines found") {
			return errorResponse(ctx, http.StatusNotFound, "no medicines found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to retrieve medicines")
	}

	return successResponse(ctx, http.StatusOK, "medicines retrieved successfully", medicines)
}

// @Summary     Get medicine by ID
// @Description Retrieve a single medicine by ID
// @ID          get-by-id-medicine
// @Tags        medicine
// @Accept      json
// @Produce     json
// @Param       id path string true "Medicine ID"
// @Success 	200 {object} response.SuccessMedicine
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicines/{id} [get]
func (r *V1) getByIDMedicine(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "medicine ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid medicine ID")
	}

	medicine, err := r.u.Medicine.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getByIDMedicine")
		if strings.Contains(err.Error(), "medicine not found") {
			return errorResponse(ctx, http.StatusNotFound, "medicine not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch medicine")
	}

	return successResponse(ctx, http.StatusOK, "medicine retrieved successfully", medicine)
}

// @Summary     Update medicine
// @Description Update an existing medicine
// @ID          update-medicine
// @Tags        medicine
// @Accept      json
// @Produce     json
// @Param       id path string true "Medicine ID"
// @Param       request body request.Medicine true "Medicine input"
// @Success 	200 {object} response.SuccessMedicine
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicines/{id} [patch]	
func (r *V1) updateMedicine(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "medicine ID is required")
	}

	var req request.Medicine
	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - updateMedicine")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - updateMedicine")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, fiber.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	updatedMedicine, err := r.u.Medicine.Update(ctx.UserContext(), id, entity.Medicine{
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - updateMedicine")
		if strings.Contains(err.Error(), "medicine not found") {
			return errorResponse(ctx, http.StatusNotFound, "medicine not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update medicine")
	}

	return successResponse(ctx, http.StatusOK, "medicine updated successfully", updatedMedicine)
}

// @Summary     Delete medicine
// @Description Delete a medicine by ID
// @ID          delete-medicine
// @Tags        medicine
// @Accept      json
// @Produce     json
// @Param       id path string true "Medicine ID"
// @Success 	200 {object} response.SuccessString
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicines/{id} [delete]
func (r *V1) deleteMedicine(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "medicine ID is required")
	}

	if err := r.u.Medicine.Delete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - deleteMedicine")
		if strings.Contains(err.Error(), "medicine not found") {
			return errorResponse(ctx, http.StatusNotFound, "medicine not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete medicine")
	}

	return successResponse(ctx, http.StatusOK, "medicine deleted successfully", "deleted")
}