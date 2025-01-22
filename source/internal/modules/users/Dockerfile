FROM golang@sha256:51a6466e8dbf3e00e422eb0f7a97ac450b2d57b33617bbe8d2ee0bddcd9d0d37 AS build

WORKDIR /app

COPY config/ ./config
COPY pkg/ ./pkg
COPY internal/modules/users/ ./internal/modules/users
COPY internal/infra/ ./internal/infra
COPY go.mod .
COPY go.sum .

RUN go mod tidy 

WORKDIR /app/internal/modules/users/cmd

RUN go build -o /server && chmod +x /server

FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

COPY --from=build /server /server

VOLUME /app/logs

EXPOSE  8000

USER nonroot:nonroot

ENTRYPOINT ["/server"]