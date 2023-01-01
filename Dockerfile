FROM golang:1.19 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bench .

FROM alpine
COPY --from=builder /app/input ./input
COPY --from=builder /app/bench .
CMD ["./bench"]