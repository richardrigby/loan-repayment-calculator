package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/richardrigby/loan-repayment-calculator/api/loan_payment_calculator/v1/go"
	"github.com/richardrigby/loan-repayment-calculator/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcAddr = ":50051"
	httpAddr = ":50052"
)

func main() {
	logger := slog.Default()
	logger.Info("Starting Loan Repayment Calculator gRPC server...")
	if err := run(logger); err != nil {
		logger.Error("Failed to run application", "error", err)
	}
}

func run(logger *slog.Logger) error {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start gRPC server
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.With("address", grpcAddr).Error("Failed to listen")
		return err
	}
	defer lis.Close()
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	// Create and register service
	service := service.NewServer(logger)
	pb.RegisterLoanPaymentCalculatorServiceServer(grpcServer, service)
	reflection.Register(grpcServer)

	// Start HTTP server (gRPC-Gateway)
	mux := runtime.NewServeMux()
	err = pb.RegisterLoanPaymentCalculatorServiceHandlerServer(ctx, mux, service)
	if err != nil {
		logger.With("address", httpAddr).Error("Failed to register HTTP server")
		return err
	}

	httpErrChan := make(chan error)
	grpcErrChan := make(chan error)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	go func() {
		logger.With("address", httpAddr).Info("HTTP server listening")

		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			logger.With("address", httpAddr).Error("Failed to start HTTP server")
			httpErrChan <- err
		}
		logger.Info("HTTP server stopped")
	}()

	go func() {
		logger.With("address", lis.Addr()).Info("gRPC server listening")
		if err := grpcServer.Serve(lis); err != nil {
			logger.With("address", lis.Addr()).Error("Failed to serve gRPC server")
			grpcErrChan <- err
		}
		logger.Info("gRPC server stopped")
	}()

	select {
	case err := <-httpErrChan:
		logger.With("error", err).Error("HTTP server error")
		return err
	case err := <-grpcErrChan:
		logger.With("error", err).Error("gRPC server error")
		return err
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		logger.Info("Shutting down server...")
		stop()
	}

	logger.Info("server shut down successfully")

	return nil
}
