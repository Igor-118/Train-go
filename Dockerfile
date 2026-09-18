FROM golang:1.25 AS bobrkurwa

WORKDIR /app

COPY . .

RUN go mod download && go build -o app .

FROM golang:1.25

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

COPY --from=bobrkurwa /app/app .
COPY --from=bobrkurwa /app/.env .
COPY --from=bobrkurwa /app/entrypoint.sh .
COPY --from=bobrkurwa /app/migrations /app/migrations

# RUN apk --no-cache add ca-certificates
ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=postgres://postgres:postgres@pg-docker:5432/postgres
CMD ["./entrypoint.sh"]