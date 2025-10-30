package calculator

import (
	"errors"
	"log/slog"

	"github.com/shopspring/decimal"
)

type PaymentTiming int32

const (
	PaymentEndOfPeriod PaymentTiming = iota
	PaymentBeginningOfPeriod
)

type PMTInput struct {
	Rate            decimal.Decimal // interest rate per period
	NumberOfPeriods int             // total number of payment periods
	PresentValue    decimal.Decimal // present value (loan amount)
	FutureValue     decimal.Decimal // future value (default is 0)
	PaymentTiming   PaymentTiming   // indicates when payments are due
}

type Calculator struct {
	logger *slog.Logger
}

func NewCalculator(logger *slog.Logger) *Calculator {
	return &Calculator{
		logger: logger,
	}
}

// PMT calculates the payment for a loan based on constant payments and a constant interest rate.
// rate: Annual interest rate per period. Divide by 12 for monthly rate.
// nper: Total number of payment periods
// pv: Present value (loan amount)
// fv: Future value (default is 0)
// paymentTiming: indicates when payments are due (0 = end of period, 1 = beginning of period)
func (c *Calculator) PMT(input PMTInput) (decimal.Decimal, error) {

	c.logger.
		With(slog.Any("rate", input.Rate.String())).
		With(slog.Any("present_value", input.PresentValue.String())).
		With(slog.Int("number_of_periods", input.NumberOfPeriods)).
		With(slog.Any("future_value", input.FutureValue.String())).
		With(slog.Any("payment_timing", input.PaymentTiming)).
		Info("Calculating PMT")

	// If fv is provided (non-zero), use future value formula
	if input.PresentValue.IsZero() && !input.FutureValue.IsZero() {
		return futureValuePMT(input.Rate, input.NumberOfPeriods, input.FutureValue, input.PaymentTiming)
	}

	return presentValuePMT(input.Rate, input.NumberOfPeriods, input.PresentValue, input.PaymentTiming)
}

// presentValuePMT calculates the payment based on present value
func presentValuePMT(rate decimal.Decimal, nper int, pv decimal.Decimal, pt PaymentTiming) (decimal.Decimal, error) {
	var pmt decimal.Decimal

	// if pv and is zero, payment is zero
	if pv.IsZero() {
		return decimal.Zero, nil
	}

	// Calculate (1 + rate)^-nper
	power := decimal.NewFromInt(1).Add(rate).Pow(decimal.NewFromInt(int64(-nper)))

	// Calculate denominator; 1 - (1 + rate)^-nper
	denominator := decimal.NewFromInt(1).Sub(power)

	// Calculate payment; PMT = pv * rate / (1 - (1 + rate)^-nper)
	pmt = pv.Mul(rate).Div(denominator)

	// Adjust for payment timing; if payments are at beginning of period, divide by (1 + rate)
	if pt == PaymentBeginningOfPeriod {
		pmt = pmt.Div(decimal.NewFromInt(1).Add(rate))
		return pmt, nil
	}

	return pmt, nil
}

// futureValuePMT calculates the payment based on future value
func futureValuePMT(rate decimal.Decimal, nper int, fv decimal.Decimal, pt PaymentTiming) (decimal.Decimal, error) {
	var pmt decimal.Decimal

	// Validate inputs
	if nper <= 0 {
		return decimal.Zero, errors.New("number of periods must be greater than zero")
	}

	// if fv is zero, payment is zero
	if fv.IsZero() {
		return decimal.Zero, nil
	}

	// Calculate (1 + rate)^nper
	power := decimal.NewFromInt(1).Add(rate).Pow(decimal.NewFromInt(int64(nper)))
	// Calculate denominator; (1 + rate)^nper - 1
	denominator := power.Sub(decimal.NewFromInt(1))
	// Calculate payment; PMT = fv * rate / ((1 + rate)^nper - 1)
	pmt = fv.Mul(rate.Div(denominator))

	// Adjust for payment timing; if payments are at beginning of period, divide by (1 + rate)
	if pt == PaymentBeginningOfPeriod {
		pmt = pmt.Div(decimal.NewFromInt(1).Add(rate))
	}

	return pmt, nil
}
