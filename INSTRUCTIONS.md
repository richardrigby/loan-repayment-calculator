# Technical Test

## Task (3-4 hours)

Create a service which exposes a REST API with a single route. Use golang as the programming language and grpc-gateway.
The purpose of the route is to calculate the monthly repayment amount for a loan given the inputs of loan amount, interest rate and number of payments. Use the PMT (excel) formula to do this calculation.
Commit your work to an open source repo. The repo should have instructions on how to run the program locally and also have instructions to deploy to k8s. We will be deploying the result to a local k8s cluster (mini kube, docker desktop, etc.).

## Presentation (1 hour)

Presentation & discussion with key stakeholders, engineering, product, and sales in the room. You have 30 mins max to present, and can choose any approach, method, tools, etc as you wish. After your presentation you can lead a 30 min discussion related to the problem and your solution with stakeholders.

Description: The PMT function calculates the payment for a loan based on constant payments and a constant interest rate.

Where:

- P = loan amount (principal)
- r = monthly interest rate (annual interest rate / 12)
- n = number of payments (loan term in months)

Loan Payment Formula (PMT):

```
PMT = P * r * (1 + r)^n / ((1 + r)^n - 1)
```


PMT = (rate * PV) / (1 - (1 + rate)^(-nper))

Where:
rate is the interest rate per period.
PV is the present value (the loan amount).
nper is the total number of payment periods.


From [Excel documentation](https://support.microsoft.com/en-us/office/pmt-function-0214da64-9a63-4996-bc20-214433fa6441):

PMT, one of the financial functions, calculates the payment for a loan based on constant payments and a constant interest rate.

## Syntax

PMT(rate, nper, pv, [fv], [type])

> Note: For a more complete description of the arguments in PMT, see the PV function.

The PMT function syntax has the following arguments:

- Rate    Required. The interest rate for the loan.

- Nper    Required. The total number of payments for the loan.

- Pv    Required. The present value, or the total amount that a series of future payments is worth now; also known as the principal.

- Fv    Optional. The future value, or a cash balance you want to attain after the last payment is made. If fv is omitted, it is assumed to be 0 (zero), that is, the future value of a loan is 0.

- Type    Optional. The number 0 (zero) or 1 and indicates when payments are due.

| **Set type equal to** |**If payments are due** |
| --- | --- |
| 0 or omitted | At the end of the period |
| 1 | At the beginning of the period |

## Remarks

- The payment returned by PMT includes principal and interest but no taxes, reserve payments, or fees sometimes associated with loans.

- Make sure that you are consistent about the units you use for specifying rate and nper. If you make monthly payments on a four-year loan at an annual interest rate of 12 percent, use 12%/12 for rate and 4*12 for nper. If you make annual payments on the same loan, use 12 percent for rate and 4 for nper.

**Tip**    To find the total amount paid over the duration of the loan, multiply the returned PMT value by nper.

## Example

Copy the example data in the following table, and paste it in cell A1 of a new Excel worksheet. For formulas to show results, select them, press F2, and then press Enter. If you need to, you can adjust the column widths to see all the data.

| **Data** |**Description** | |
| --- | --- | --- |
| 8% | Annual interest rate | |
| 10 | Number of months of payments | |
| $10,000 | Amount of loan | |
| --- | --- | --- |
| **Formula**| **Description** | **Result** |
| `=PMT(A2/12,A3,A4)` | Monthly payment for a loan with terms specified as arguments in A2:A4. | ($1,037.03) |
| `=PMT(A2/12,A3,A4,,1)` | Monthly payment for a loan with with terms specified as arguments in A2:A4, except payments are due at the beginning of the period. | ($1,030.61) |
| **Data** |**Description** | |
| --- | --- | --- |
| 6% | Annual interest rate | |
| 18 | Number of months of payments | |
| $50,000 | Amount of loan | |
| --- | --- | --- |
| **Formula**| **Description** | **Result** |
| `PMT(A9/12,A10*12, 0,A11)` | Amount to save each month to have $50,000 at the end of 18 years. | ($129.08) |

