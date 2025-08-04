// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-eprescription-clean/config"
	"go-eprescription-clean/internal/controller/http"
	v1 "go-eprescription-clean/internal/controller/http/v1"
	"go-eprescription-clean/internal/repo/payment_gateway"
	"go-eprescription-clean/internal/repo/persistent"
	"go-eprescription-clean/internal/usecase/medicine"
	"go-eprescription-clean/internal/usecase/medicine_detail"
	"go-eprescription-clean/internal/usecase/patient"
	"go-eprescription-clean/internal/usecase/signa"
	"go-eprescription-clean/internal/usecase/transaction"

	"go-eprescription-clean/pkg/httpserver"
	"go-eprescription-clean/pkg/logger"
	"go-eprescription-clean/pkg/midtrans"
	"go-eprescription-clean/pkg/xendit"
	"go-eprescription-clean/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Midtrans
	mtClient := midtrans.New()

	// Xendit
	xenditClient := xendit.New()
	
	// Use-Case
	signaUseCase := signa.New(persistent.NewSignaRepo(pg))
	patientUseCase := patient.New(persistent.NewPatientRepo(pg))
	medicineUseCase := medicine.New(persistent.NewMedicineRepo(pg))
	transactionUseCase := transaction.New(
		persistent.NewTransactionRepo(pg), 
		persistent.NewMedicineDetailsRepo(pg), 
		persistent.NewMedicineRepo(pg),
		persistent.NewPatientRepo(pg),
		persistent.NewSignaRepo(pg),
		payment_gateway.NewMidtransRepo(mtClient),
		payment_gateway.NewXenditRepo(xenditClient),
	)
	medicineDetailUseCase := medicine_detail.New(persistent.NewMedicineDetailsRepo(pg))

	usecases := v1.Usecases{
		Signa: signaUseCase,
		Patient: patientUseCase,
		Medicine: medicineUseCase,
		Transaction: transactionUseCase,
		MedicineDetail: medicineDetailUseCase,
	}

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, usecases, l)

	// Start servers
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
