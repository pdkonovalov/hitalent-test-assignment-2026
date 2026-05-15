FROM golang:1.26 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM gcr.io/distroless/static-debian13 AS release-stage

WORKDIR /

COPY --from=build-stage /server /server

USER nonroot:nonroot

ENTRYPOINT ["/server"]
