FROM golang:1.22.0

WORKDIR /usr/src/app

# Install swag for swagger documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Remove old docs if exists and generate fresh swagger docs
RUN rm -rf docs && swag init --parseDependency --parseInternal

# Run go mod tidy after docs are generated
RUN go mod tidy

EXPOSE 8081

CMD ["go", "run", "./main.go", "api"]

