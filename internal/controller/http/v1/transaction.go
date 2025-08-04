package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go-eprescription-clean/internal/controller/http/v1/request"
	"go-eprescription-clean/internal/entity"
)

// @Summary     Create transaction with medicine details
// @Description Create a new transaction with medicine details
// @ID          create-transaction-with-med-detail
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Param       request body request.CreateTransaction true "Transaction input"
// @Success 	201 {object} response.SuccessTransaction
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions [post]
func (r *V1) createTransactionWithMedDetail(ctx *fiber.Ctx) error {
	var req request.CreateTransaction

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - createTransactionWithMedDetail")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - createTransactionWithMedDetail - validation")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, http.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	transaction, err := r.u.Transaction.CreateWithMedicineDetail(ctx.UserContext(), entity.Transaction{
		PatientID:    req.PatientID,
		MedicineType: req.MedicineType,
	}, req.Medicines, req.Signas, req.Descriptions, req.Quantities)

	if err != nil {
		r.l.Error(err, "http - v1 - createTransactionWithMedDetail")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to create transaction")
	}

	return successResponse(ctx, http.StatusCreated, "transaction created successfully", transaction)
}

// @Summary     Get all transactions
// @Description Retrieve all transaction records
// @ID          get-all-transaction
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Success 	200 {object} response.SuccessTransactionList
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions [get]
func (r *V1) getAllTransactions(ctx *fiber.Ctx) error {
	transactions, err := r.u.Transaction.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllTransaction")
		if strings.Contains(err.Error(), "no transactions found") {
			return errorResponse(ctx, http.StatusNotFound, "no transactions found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch transactions")
	}

	return successResponse(ctx, http.StatusOK, "transactions retrieved successfully", transactions)
}

// @Summary     Get all transactions by patient ID
// @Description Retrieve all transactions for a specific patient
// @ID          get-all-transactions-by-patient-id
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Param       patient_id path string true "Patient ID"
// @Success 	200 {object} response.SuccessTransactionList
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions/patients/{patient_id} [get]
func (r *V1) getAllTransactionsByPatientID(ctx *fiber.Ctx) error {
	patientID := ctx.Params("patient_id")
	if patientID == "" {
		r.l.Error("http - v1 - getAllTransactionsByPatientID")
		return errorResponse(ctx, http.StatusBadRequest, "patient ID is required")
	}

	if err := r.v.Var(patientID, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - getByIDTransaction")
		return errorResponse(ctx, http.StatusBadRequest, "invalid transaction ID")
	}

	transactions, err := r.u.Transaction.GetAllByPatientID(ctx.UserContext(), patientID)
	if err != nil {
		r.l.Error(err, "http - v1 - getAllTransactionsByPatientID")
		if strings.Contains(err.Error(), "no transactions found") {
			return errorResponse(ctx, http.StatusNotFound, "no transactions found for this patient")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch transactions for this patient")
	}

	return successResponse(ctx, http.StatusOK, "transactions retrieved successfully", transactions)
}

// @Summary     Get transaction by ID
// @Description Retrieve a single transaction by ID
// @ID          get-transaction-by-id
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Param       id path string true "Transaction ID"
// @Success 	200 {object} response.SuccessTransaction
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions/{id} [get]
func (r *V1) getByIDTransaction(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		r.l.Error("http - v1 - getByIDTransaction")
		return errorResponse(ctx, http.StatusBadRequest, "transaction ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - getByIDTransaction")
		return errorResponse(ctx, http.StatusBadRequest, "invalid transaction ID")
	}

	transaction, err := r.u.Transaction.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getByIDTransaction")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to retrieve transaction")
	}

	return successResponse(ctx, http.StatusOK, "transaction retrieved successfully", transaction)
}

// @Summary     Update transaction
// @Description Update a transaction by ID
// @ID          update-transaction
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Param       id path string true "Transaction ID"
// @Param       request body request.UpdateTransaction true "Transaction input"
// @Success 	200 {object} response.SuccessTransaction
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions/{id} [patch]
func (r *V1) updateTransaction(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "transaction ID is required")
	}

	if err := r.v.Var(id, "uuid4"); err != nil {
		r.l.Error(err, "http - v1 - updateTransaction - uuid validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid transaction ID")
	}

	var req request.UpdateTransaction

	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error(err, "http - v1 - updateTransaction")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(req); err != nil {
		r.l.Error(err, "http - v1 - updateTransaction - validation")
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errorResponse(ctx, http.StatusBadRequest, formatValidationErrors(ve))
		}
	}

	transaction, err := r.u.Transaction.Update(ctx.UserContext(), id, entity.Transaction{
		MedicineType: req.MedicineType,
		TotalPrice:  req.TotalPrice,
		TotalMedicines: req.TotalMedicines,
	})

	if err != nil {
		r.l.Error(err, "http - v1 - updateTransaction")
		if strings.Contains(err.Error(), "transaction not found") {
			return errorResponse(ctx, http.StatusNotFound, "transaction not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update transaction")
	}

	return successResponse(ctx, http.StatusOK, "transaction updated successfully", transaction)
}

// @Summary     Delete transaction
// @Description Delete a transaction by ID
// @ID          delete-transaction
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Param       id path string true "Transaction ID"
// @Success 	200 {object} response.SuccessString
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions/{id} [delete]
func (r *V1) deleteTransaction(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := r.u.Transaction.Delete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - deleteTransaction")
		if strings.Contains(err.Error(), "transaction not found") {
			return errorResponse(ctx, http.StatusNotFound, "transaction not found")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete transaction")
	}

	return successResponse(ctx, http.StatusOK, "transaction deleted successfully", "deleted")
}

// @Summary     Handle Midtrans callback notification
// @Description Receives and processes callback notifications from Midtrans
// @ID          handle-midtrans-callback
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Param       payload body object true "Midtrans Notification Payload"
// @Success     200 {object} response.SuccessString
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions/midtrans/callbacks [post]
func (r *V1) handleMidtransCallback(ctx *fiber.Ctx) error {
	body := ctx.Body()
	r.l.Info(fmt.Sprintf("Midtrans callback payload: %s", string(body)))

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		r.l.Error(fmt.Sprintf("Failed to unmarshal Midtrans callback: %v", err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid midtrans payload")
	}

	transactionID, ok1 := payload["order_id"].(string)
	transactionStatus, ok2 := payload["transaction_status"].(string)
	fraudStatus, ok3 := payload["fraud_status"].(string)

	if !ok1 || !ok2 || !ok3 {
		r.l.Error("Missing required fields in midtrans callback")
		return errorResponse(ctx, http.StatusBadRequest, "missing fields in midtrans payload")
	}

	_, err := r.u.Transaction.HandleMidtransNotification(ctx.UserContext(), transactionID, transactionStatus, fraudStatus)
	if err != nil {
		r.l.Error(fmt.Sprintf("Failed to handle midtrans notification: %v", err))
		return errorResponse(ctx, http.StatusInternalServerError, "failed to process midtrans callback")
	}

	return ctx.SendStatus(http.StatusOK)
}

// @Summary     Handle Xendit callback notification
// @Description Receives and processes callback notifications from Xendit
// @ID          handle-xendit-callback
// @Tags        transaction
// @Accept      json
// @Produce     json
// @Param       payload body object true "Xendit Notification Payload"
// @Success     200 {object} response.SuccessString
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /transactions/xendit/callbacks [post]
func (r *V1) handleXenditCallback(ctx *fiber.Ctx) error {
	body := ctx.Body()
	r.l.Info(fmt.Sprintf("Xendit callback payload: %s", string(body)))

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		r.l.Error(fmt.Sprintf("Failed to unmarshal Xendit callback: %v", err))
		return errorResponse(ctx, http.StatusBadRequest, "invalid xendit payload")
	}

	//print transaction ID, status, and fraud status
	r.l.Info(fmt.Sprintf("Xendit callback transaction ID: %s", payload["external_id"]))
	r.l.Info(fmt.Sprintf("Xendit callback status: %s", payload["status"]))

	transactionID, ok1 := payload["external_id"].(string)
	transactionStatus, ok2 := payload["status"].(string)

	if !ok1 || !ok2 {
		r.l.Error("Missing required fields in xendit callback")
		return errorResponse(ctx, http.StatusBadRequest, "missing fields in xendit payload")
	}

	_, err := r.u.Transaction.HandleXenditNotification(ctx.UserContext(), transactionID, transactionStatus)
	if err != nil {
		r.l.Error(fmt.Sprintf("Failed to handle xendit notification: %v", err))
		return errorResponse(ctx, http.StatusInternalServerError, "failed to process xendit callback")
	}

	return ctx.SendStatus(http.StatusOK) 
}

