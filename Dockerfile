FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /hris ./cmd/api


FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /hris /hris
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/hris"]
