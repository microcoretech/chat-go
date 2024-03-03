FROM golang:1.21-alpine as builder

RUN apk --no-cache add ca-certificates git

WORKDIR /go/src/gitlab.com/chat605743/backend

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./app ./cmd/backend \
    && chmod +x ./app

FROM alpine
WORKDIR /
COPY --from=builder /go/src/gitlab.com/chat605743/backend/app .
COPY /VERSION .
ENTRYPOINT ["./app"]