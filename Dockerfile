FROM golang:1.22 as builder

WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

COPY . .

RUN make setup
RUN make swag
RUN make wire
RUN make build

# FROM alpine:latest

# RUN apk --no-cache add ca-certificates

# WORKDIR /root/

# COPY --from=builder /app/main .

CMD ["./build/apiserver"]

# Path: Dockerfile

