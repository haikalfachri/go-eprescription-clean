package v1

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go-eprescription-clean/internal/controller/http/v1/request"
	"go-eprescription-clean/internal/entity"
	"github.com/go-playground/validator/v10"
)

// @Summary     Create medicinedetail
// @Description Create a new medicine detail
// @ID          create-medicine-detail
// @Tags        medicine-detail
// @Accept      json
// @Produce     json
// @Param       request body request.MedicineDetail true "MedicineDetail input"
// @Success 	201 {object} response.SuccessMedicineDetail
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicinedetails [post]
func (r *V1) createMedicineDetail(ctx *fiber.Ctx) error {
	var req request.MedicineDetail

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - createMedicineDetail")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - createMedicineDetail - validation")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, http.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	medicinedetail, err := r.u.MedicineDetail.Create(ctx.UserContext(), entity.MedicineDetail{
		MedicineID: req.MedicineID,
		SignaID:    req.SignaID,
		Description: req.Description,
		Quantity:    req.Quantity,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - createMedicineDetail")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to create medicinedetail")
	}

	return successResponse(ctx, http.StatusCreated, "medicinedetails created successfully", medicinedetail)
}

// @Summary     Get all medicine details
// @Description Retrieve all medicine detail records
// @ID          get-all-medicine-details
// @Tags        medicine-detail
// @Accept      json
// @Produce     json
// @Success 	200 {object} response.SuccessMedicineDetailList
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicinedetails [get]
func (r *V1) getAllMedicineDetail(ctx *fiber.Ctx) error {
	medicineDetails, err := r.u.MedicineDetail.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllMedicineDetails")
		if strings.Contains(err.Error(), "no medicine details found") {
			return errorResponse(ctx, http.StatusNotFound, "no medicine details found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch medicine details")
	}

	return successResponse(ctx, http.StatusOK, "medicine details retrieved successfully", medicineDetails)
}

// @Summary     Get medicine details by ID
// @Description Retrieve a single medicine detail by ID
// @ID          get-medicine-detail-by-id
// @Tags        medicine-detail
// @Accept      json
// @Produce     json
// @Param       id path string true "MedicineDetail ID"
// @Success 	200 {object} response.SuccessMedicineDetail
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicinedetails/{id} [get]
func (r *V1) getByIDMedicineDetail(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		r.l.Error("http - v1 - getByIDMedicineDetails")
		return errorResponse(ctx, http.StatusBadRequest, "medicinedetail ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - getByIDMedicineDetails")
		return errorResponse(ctx, http.StatusBadRequest, "invalid medicinedetail ID")
	}

	medicinedetail, err := r.u.MedicineDetail.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getMedicineDetailsByID")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to retrieve medicinedetail")
	}

	return successResponse(ctx, http.StatusOK, "medicinedetails retrieved successfully", medicinedetail)
}

// @Summary     Update medicinedetail
// @Description Update a medicine detail by ID
// @ID          update-medicine-detail
// @Tags        medicine-detail
// @Accept      json
// @Produce     json
// @Param       id path string true "MedicineDetail ID"
// @Param       request body request.MedicineDetail true "MedicineDetail input"
// @Success 	200 {object} response.SuccessMedicineDetail
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicinedetails/{id} [patch]
func (r *V1) updateMedicineDetail(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "medicine detail ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - updateMedicineDetail - uuid validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid medicine detail ID")
	}

	var req request.MedicineDetail

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - updateMedicineDetail")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - updateMedicineDetail - validation")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, http.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	medicinedetail, err := r.u.MedicineDetail.Update(ctx.UserContext(), id, entity.MedicineDetail{
		Description: req.Description,
		MedicineID:  req.MedicineID,
		SignaID:     req.SignaID,
		Quantity:    req.Quantity,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - updateMedicineDetail")
		if strings.Contains(err.Error(), "medicinedetail not found") {
			return errorResponse(ctx, http.StatusNotFound, "medicine detail not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update medicine detail")
	}

	return successResponse(ctx, http.StatusOK, "medicine detail updated successfully", medicinedetail)
}


// @Summary     Delete medicine detail
// @Description Delete a medicine detail by ID
// @ID          delete-medicine-detail
// @Tags        medicine-detail
// @Accept      json
// @Produce     json
// @Param       id path string true "MedicineDetail ID"
// @Success 	200 {object} response.SuccessString
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /medicinedetails/{id} [delete]
func (r *V1) deleteMedicineDetail(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := r.u.MedicineDetail.Delete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - deleteMedicineDetail")
		if strings.Contains(err.Error(), "medicinedetail not found") {
			return errorResponse(ctx, http.StatusNotFound, "medicinedetail not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete medicinedetail")
	}

	return successResponse(ctx, http.StatusOK, "medicinedetail deleted successfully", "deleted")
}
