FROM golang:1.24.2 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o api .

FROM gcr.io/distroless/base-debian11
COPY --from=build /app/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
