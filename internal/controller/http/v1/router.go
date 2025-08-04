package v1

import (
	"go-eprescription-clean/internal/usecase"
	"go-eprescription-clean/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewSignaRoutes -.
func NewSignaRoutes(apiV1Group fiber.Router, s usecase.Signa, l logger.Interface) {
	r := &V1{
		u: Usecases{
			Signa: s,
		},
		l: l,
		v: validator.New(),
	}

	signaGroup := apiV1Group.Group("/signas")
	{
		signaGroup.Get("/", r.getAllSigna)
		signaGroup.Get("/:id", r.getByIDSigna)
		signaGroup.Post("/", r.createSigna)
		signaGroup.Patch("/:id", r.updateSigna)
		signaGroup.Delete("/:id", r.deleteSigna)
	}
}

// NewPatientRoutes -.
func NewPatientRoutes(apiV1Group fiber.Router, p usecase.Patient, l logger.Interface) {
	r := &V1{
		u: Usecases{
			Patient: p,
		},
		l: l,
		v: validator.New(),
	}

	patientGroup := apiV1Group.Group("/patients")
	{
		patientGroup.Post("/", r.createPatient)
		patientGroup.Get("/", r.getAllPatients)
		patientGroup.Get("/:id", r.getByIDPatient)
		patientGroup.Patch("/:id", r.updatePatient)
		patientGroup.Delete("/:id", r.deletePatient)
	}
}

// NewMedicineRoutes -.
func NewMedicineRoutes(apiV1Group fiber.Router, m usecase.Medicine, l logger.Interface) {
	r := &V1{
		u: Usecases{
			Medicine: m,
		},
		l: l,
		v: validator.New(),
	}

	medicineGroup := apiV1Group.Group("/medicines")
	{
		medicineGroup.Get("/", r.getAllMedicines)
		medicineGroup.Get("/:id", r.getByIDMedicine)
		medicineGroup.Post("/", r.createMedicine)
		medicineGroup.Patch("/:id", r.updateMedicine)
		medicineGroup.Delete("/:id", r.deleteMedicine)
	}
}

// NewTransactionRoutes -.
func NewTransactionRoutes(apiV1Group fiber.Router, t usecase.Transaction, md usecase.MedicineDetail, m usecase.Medicine, p usecase.Patient, l logger.Interface) {
	r := &V1{
		u: Usecases{
			Transaction:   t,
			MedicineDetail: md,
			Medicine: m,
			Patient: p,
		},
		l: l,
		v: validator.New(),
	}

	transactionGroup := apiV1Group.Group("/transactions")
	{
		transactionGroup.Get("/", r.getAllTransactions)
		transactionGroup.Get("/:id", r.getByIDTransaction)
		transactionGroup.Get("/patients/:patient_id", r.getAllTransactionsByPatientID)
		transactionGroup.Post("/", r.createTransactionWithMedDetail)
		transactionGroup.Patch("/:id", r.updateTransaction)
		transactionGroup.Delete("/:id", r.deleteTransaction)
		transactionGroup.Post("/midtrans/callbacks", r.handleMidtransCallback)
		transactionGroup.Post("/xendit/callbacks", r.handleXenditCallback)
	}
}


// NewMedicineDetailRoutes -.
func NewMedicineDetailRoutes(apiV1Group fiber.Router, md usecase.MedicineDetail, l logger.Interface) {
	r := &V1{
		u: Usecases{
			MedicineDetail: md,
		},
		l: l,
		v: validator.New(),
	}

	medicineDetailGroup := apiV1Group.Group("/medicine-details")
	{
		medicineDetailGroup.Get("/", r.getAllMedicineDetail)
		medicineDetailGroup.Get("/:id", r.getByIDMedicineDetail)
		medicineDetailGroup.Post("/", r.getByIDTransaction)
		medicineDetailGroup.Post("/", r.createMedicineDetail)
		medicineDetailGroup.Patch("/:id", r.updateMedicineDetail)
		medicineDetailGroup.Delete("/:id", r.deleteMedicineDetail)
	}
}