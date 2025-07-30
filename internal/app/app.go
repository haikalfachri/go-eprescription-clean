// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-eprescription-clean/config"
	// amqprpc "go-eprescription-clean/internal/controller/amqp_rpc"
	// "go-eprescription-clean/internal/controller/grpc"
	"go-eprescription-clean/internal/controller/http"
	v1 "go-eprescription-clean/internal/controller/http/v1"
	"go-eprescription-clean/internal/repo/persistent"
	"go-eprescription-clean/internal/usecase/medicine"
	"go-eprescription-clean/internal/usecase/patient"
	"go-eprescription-clean/internal/usecase/signa"
	"go-eprescription-clean/internal/usecase/transaction"
	"go-eprescription-clean/internal/usecase/medicine_detail"

	// "go-eprescription-clean/pkg/grpcserver"
	"go-eprescription-clean/pkg/httpserver"
	"go-eprescription-clean/pkg/logger"
	"go-eprescription-clean/pkg/postgres"
	// "go-eprescription-clean/pkg/rabbitmq/rmq_rpc/server"
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

	// Use-Case
	signaUseCase := signa.New(persistent.NewSignaRepo(pg))
	patientUseCase := patient.New(persistent.NewPatientRepo(pg))
	medicineUseCase := medicine.New(persistent.NewMedicineRepo(pg))
	transactionUseCase := transaction.New(persistent.NewTransactionRepo(pg), persistent.NewMedicineDetailsRepo(pg), persistent.NewMedicineRepo(pg))
	medicineDetailUseCase := medicine_detail.New(persistent.NewMedicineDetailsRepo(pg))

	usecases := v1.Usecases{
		Signa: signaUseCase,
		Patient: patientUseCase,
		Medicine: medicineUseCase,
		Transaction: transactionUseCase,
		MedicineDetail: medicineDetailUseCase,
	}

	// RabbitMQ RPC Server
	// rmqRouter := amqprpc.NewRouter(translationUseCase, l)

	// rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	// if err != nil {
	// 	l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	// }

	// gRPC Server
	// grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	// grpc.NewRouter(grpcServer.App, translationUseCase, l)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, usecases, l)

	// Start servers
	// rmqServer.Start()
	// grpcServer.Start()
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// case err = <-grpcServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
		// case err = <-rmqServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
		// }

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}

		// err = grpcServer.Shutdown()
		// if err != nil {
		// 	l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
		// }

		// err = rmqServer.Shutdown()
		// if err != nil {
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
		// }
	}
}
