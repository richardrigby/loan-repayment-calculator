package calculator

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestPMT_Basic(t *testing.T) {
	cases := []struct {
		name     string
		rate     decimal.Decimal
		nper     int
		pv       decimal.Decimal
		fv       decimal.Decimal
		pt       PaymentTiming
		expected decimal.Decimal
	}{
		{
			name: "8% per year, 10 months, pv=10000, end",
			rate: decimal.NewFromFloat(0.08).Div(decimal.NewFromInt(12)),
			nper: 10, // 10 months
			pv:   decimal.NewFromFloat(10_000),
			fv:   decimal.Zero,
			pt:   PaymentEndOfPeriod,
			expected: func() decimal.Decimal {
				num, _ := decimal.NewFromString("1037.0320893591529985")
				return num
			}(),
		},
		{
			name: "0.1% per month, 150 months, pv=10000, end",
			rate: decimal.NewFromFloat(0.001),
			nper: 150,
			pv:   decimal.NewFromFloat(10_000),
			fv:   decimal.Zero,
			pt:   PaymentEndOfPeriod,
			expected: func() decimal.Decimal {
				num, _ := decimal.NewFromString("71.8248852091213088")
				return num
			}(),
		},
		{
			name: "8% per year, 10 months, pv=10000, beginning",
			rate: decimal.NewFromFloat(0.08).Div(decimal.NewFromInt(12)),
			nper: 10, // 10 months
			pv:   decimal.NewFromFloat(10_000),
			fv:   decimal.Zero,
			pt:   PaymentBeginningOfPeriod,
			expected: func() decimal.Decimal {
				num, _ := decimal.NewFromString("1030.1643271779665207")
				return num
			}(),
		},
		{
			name:     "6% per period, 18 periods, fv=50000, end",
			rate:     decimal.NewFromFloat(0.06).Div(decimal.NewFromInt(12)),
			nper:     18 * 12, // 18 years monthly
			pv:       decimal.Zero,
			fv:       decimal.NewFromFloat(50_000),
			pt:       PaymentEndOfPeriod,
			expected: decimal.NewFromFloat(129.08116086799),
		},
		{
			name:     "6% per year, 18 years, fv=50000, end",
			rate:     decimal.NewFromFloat(0.06).Div(decimal.NewFromInt(12)),
			nper:     18 * 12, // 18 years monthly
			pv:       decimal.Zero,
			fv:       decimal.NewFromFloat(50_000),
			pt:       PaymentEndOfPeriod,
			expected: decimal.NewFromFloat(129.08116086799),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PMT(tc.rate, tc.nper, tc.pv, tc.fv, tc.pt)
			if err != nil {
				t.Fatalf("PMT returned unexpected error: %v", err)
			}

			if got.Cmp(tc.expected) != 0 {
				t.Fatalf("PMT result: got %s, expected %s", got.String(), tc.expected.String())
			}
		})
	}
}
