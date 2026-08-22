FROM golang:1.27.0-alpine3.24 AS builder
WORKDIR /debter
COPY . .
RUN go build -o main cmd/main.go

FROM golang:1.27.0-alpine3.24 AS server
WORKDIR /debter
COPY --from=builder /debter/main .
COPY --from=builder /debter/.env .
COPY --from=builder /debter/internal/migration ./migration
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY --from=builder /debter/start.sh .
EXPOSE 8080
CMD [ "./main" ]
ENTRYPOINT [ "./start.sh" ]