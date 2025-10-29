package calculator

import (
	"errors"

	"github.com/shopspring/decimal"
)

type PaymentTiming int

const (
	PaymentEndOfPeriod PaymentTiming = iota
	PaymentBeginningOfPeriod
)

// PMT calculates the payment for a loan based on constant payments and a constant interest rate.
// rate: Annual interest rate per period. Divide by 12 for monthly rate.
// nper: Total number of payment periods
// pv: Present value (loan amount)
// fv: Future value (default is 0)
// paymentTiming: indicates when payments are due (0 = end of period, 1 = beginning of period)
func PMT(rate decimal.Decimal, nper int, pv decimal.Decimal, fv decimal.Decimal, pt PaymentTiming) (decimal.Decimal, error) {
	// Validate inputs
	if nper <= 0 {
		return decimal.Zero, errors.New("number of periods must be greater than zero")
	}

	// If fv is provided (non-zero), use future value formula
	if pv.IsZero() && !fv.IsZero() {
		return futureValuePMT(rate, nper, fv, pt)
	}

	return presentValuePMT(rate, nper, pv, pt)
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
