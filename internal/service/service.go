package service

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/richardrigby/loan-repayment-calculator/api/loan_payment_calculator/v1/go"
	"github.com/richardrigby/loan-repayment-calculator/internal/calculator"
	"github.com/shopspring/decimal"
)

// Server implements the LoanPaymentCalculatorServiceServer interface
type Server struct {
	pb.UnimplementedLoanPaymentCalculatorServiceServer
	logger     *slog.Logger
	calculator *calculator.Calculator
}

func NewServer(logger *slog.Logger) *Server {
	return &Server{
		logger:     logger,
		calculator: calculator.NewCalculator(logger),
	}
}

// SayHello implements greet.GreeterServer
func (s *Server) CalculatePMT(ctx context.Context, in *pb.CalculatePMTRequest) (*pb.CalculatePMTResponse, error) {
	s.logger.With(slog.Any("in", in)).Info("Received CalculatePMT request")

	// parse input values
	input, err := grpcToDomainCalculatePMTRequest(in)
	if err != nil {
		s.logger.With(slog.Any("error", err)).Error("Failed to parse request")
		return nil, err
	}
	// Perform calculation here
	pmt, err := s.calculator.PMT(input)
	if err != nil {
		s.logger.With(slog.Any("error", err)).Error("Failed to calculate PMT")
		return nil, err
	}
	return &pb.CalculatePMTResponse{Payment: &pb.Decimal{Value: pmt.StringFixedBank(2)}}, nil
}

func grpcToDomainCalculatePMTRequest(in *pb.CalculatePMTRequest) (calculator.PMTInput, error) {
	if in == nil {
		return calculator.PMTInput{}, errors.New("input request is nil")
	}
	if in.Rate == nil {
		return calculator.PMTInput{}, errors.New("rate must be provided")
	}
	if in.PresentValue == nil && in.FutureValue == nil {
		return calculator.PMTInput{}, errors.New("present value or future value must be provided")
	}
	if in.PresentValue != nil && in.FutureValue != nil {
		return calculator.PMTInput{}, errors.New("only one of present value or future value should be provided")
	}
	if in.Nper <= 0 {
		return calculator.PMTInput{}, errors.New("number of periods must be greater than or equal to zero")
	}

	rate, err := decimal.NewFromString(in.Rate.Value)
	if err != nil {
		return calculator.PMTInput{}, err
	}
	presentValue, err := decimal.NewFromString(in.PresentValue.Value)
	if err != nil {
		return calculator.PMTInput{}, err
	}
	futureValue := decimal.Zero
	if in.FutureValue != nil {
		futureValue, err = decimal.NewFromString(in.FutureValue.Value)
		if err != nil {
			return calculator.PMTInput{}, err
		}
	}

	var paymentTiming calculator.PaymentTiming
	if in.PaymentTiming != nil {
		// Convert the protobuf enum pointer to the domain enum type
		if enumPtr := in.PaymentTiming.Enum(); enumPtr != nil {
			paymentTiming = calculator.PaymentTiming(*enumPtr)
		}
	}

	// map grpc request to domain model
	return calculator.PMTInput{
		Rate:            rate,
		NumberOfPeriods: int(in.Nper),
		PresentValue:    presentValue,
		FutureValue:     futureValue,
		PaymentTiming:   paymentTiming,
	}, nil
}
