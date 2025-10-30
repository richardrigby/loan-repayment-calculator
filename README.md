# Loan Repayment Calculator

Calculate the monthly repayment amount for a loan using the PMT formula.

## Overview

This application provides a simple API to calculate the monthly repayment amount for a loan based on the loan amount, interest rate, and number of payments using the PMT formula.

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Docker
- Kubernetes cluster (e.g., Minikube, Docker Desktop)
- kubectl

### Running Locally

There is a Makefile provided to simplify the setup and running of the application locally.

#### To generate the Go gRPC code from the protobuf definitions:

```bash
make gen
```

#### To run the application locally

```bash
make run
```

The gRPC server will be available at `http://localhost:50051` and the gRPC-Gateway at `http://localhost:50052`.

#### To run the tests

```bash
make test
```

#### To run the application with Docker

```bash
make run-docker
```

This will build the Docker image and run the container.
The gRPC server will be available at `http://localhost:50051` and the gRPC-Gateway at `http://localhost:50052`.

### Deploying to Kubernetes

First build the Docker image:

```bash
make build
```

Then deploy to your Kubernetes cluster:

```bash
make run-k8s
```

This will create a deployment and service in your Kubernetes cluster.

You can access the server using port forwarding:

```bash
make port-forward
```

The gRPC server will be available at `http://localhost:50051` and the gRPC-Gateway at `http://localhost:50052`.

## API Documentation

The API exposes a single endpoint to calculate the monthly repayment amount.

### Endpoint

```bash
POST /v1/calculate-pmt
```

Request Body:

```json
{
    "rate": {
        "value": "string" // Interest rate per period (e.g., monthly rate)
    },
    "nper": "int", // Total number of payments
    "present_value": {
        "value": "string" // Present value or principal amount
    },
    "future_value": {
        "value": "string" // Future value, not available in conjunction with present_value (optional, default is 0)
    },
    "payment_timing": "int" // 0 for end of period, 1 for beginning of period (optional, default is 0)
}
```

⚠️ Note: `rate`, `present_value`, and `future_value` are represented as strings to accommodate high precision decimal values.

#### Example cURL Request

```json
curl --location 'http://localhost:50052/v1/calculate-pmt' \
--header 'Content-Type: application/json' \
--data '{
    "nper": 10,
    "payment_timing": 1,
    "present_value": {
        "value": "10000"
    },
    "rate": {
        "value": "0.006666666667"
    }
}'
```

#### Response

```json
{
    "payment": {
        "value": "1030.16"
    }
}
```

### Testing with Postman

There is a public Postman collection provided [here](https://www.postman.com/rrig004/workspace/loan-repayment-calculator) for manual testing of the API.
