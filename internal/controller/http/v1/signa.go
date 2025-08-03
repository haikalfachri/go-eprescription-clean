package v1

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go-eprescription-clean/internal/controller/http/v1/request"
	"go-eprescription-clean/internal/entity"
	"github.com/go-playground/validator/v10"
)

// @Summary     Create signa
// @Description Create a new signa
// @ID          create-signa
// @Tags        signa
// @Accept      json
// @Produce     json
// @Param       request body request.Signa true "Signa input"
// @Success 	201 {object} response.SuccessSigna
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /signas [post]
func (r *V1) createSigna(ctx *fiber.Ctx) error {
	var req request.Signa

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - createSigna")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - createSigna - validation")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, http.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	signa, err := r.u.Signa.Create(ctx.UserContext(), entity.Signa{
		Signa:       req.Signa,
		Description: req.Description,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - createSigna")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to create signa")
	}

	return successResponse(ctx, http.StatusCreated, "signa created successfully", signa)
}

// @Summary     Get all signas
// @Description Retrieve all signa records
// @ID          get-all-signa
// @Tags        signa
// @Accept      json
// @Produce     json
// @Success 	200 {object} response.SuccessSignaList
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /signas [get]
func (r *V1) getAllSigna(ctx *fiber.Ctx) error {
	signas, err := r.u.Signa.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllSigna")
		if strings.Contains(err.Error(), "no signas found") {
			return errorResponse(ctx, http.StatusNotFound, "no signas found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch signas")
	}

	return successResponse(ctx, http.StatusOK, "signas retrieved successfully", signas)
}

// @Summary     Get signa by ID
// @Description Retrieve a single signa by ID
// @ID          get-signa-by-id
// @Tags        signa
// @Accept      json
// @Produce     json
// @Param       id path string true "Signa ID"
// @Success 	200 {object} response.SuccessSigna
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /signas/{id} [get]
func (r *V1) getByIDSigna(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		r.l.Error("http - v1 - getByIDSigna")
		return errorResponse(ctx, http.StatusBadRequest, "signa ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - getByIDSigna")
		return errorResponse(ctx, http.StatusBadRequest, "invalid signa ID")
	}

	signa, err := r.u.Signa.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getSignaByID")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to retrieve signa")
	}

	return successResponse(ctx, http.StatusOK, "signa retrieved successfully", signa)
}

// @Summary     Update signa
// @Description Update a signa by ID
// @ID          update-signa
// @Tags        signa
// @Accept      json
// @Produce     json
// @Param       id path string true "Signa ID"
// @Param       request body request.Signa true "Signa input"
// @Success 	200 {object} response.SuccessSigna
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /signas/{id} [patch]
func (r *V1) updateSigna(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "signa ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - updateSigna - uuid validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid signa ID")
	}

	var req request.Signa

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - updateSigna")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - updateSigna - validation")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, http.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	signa, err := r.u.Signa.Update(ctx.UserContext(), id, entity.Signa{
		Signa:       req.Signa,
		Description: req.Description,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - updateSigna")
		if strings.Contains(err.Error(), "signa not found") {
			return errorResponse(ctx, http.StatusNotFound, "signa not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update signa")
	}

	return successResponse(ctx, http.StatusOK, "signa updated successfully", signa)
}


// @Summary     Delete signa
// @Description Delete a signa by ID
// @ID          delete-signa
// @Tags        signa
// @Accept      json
// @Produce     json
// @Param       id path string true "Signa ID"
// @Success 	200 {object} response.SuccessString
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /signas/{id} [delete]
func (r *V1) deleteSigna(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := r.u.Signa.Delete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - deleteSigna")
		if strings.Contains(err.Error(), "signa not found") {
			return errorResponse(ctx, http.StatusNotFound, "signa not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete signa")
	}

	return successResponse(ctx, http.StatusOK, "signa deleted successfully", "deleted")
}
