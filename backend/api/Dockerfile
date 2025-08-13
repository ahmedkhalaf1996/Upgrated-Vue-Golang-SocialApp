FROM golang:1.23.0 AS builder

WORKDIR /usr/src/app

# Copy go.mod and go.sum first (for caching dependencies)
COPY backend/api/go.mod ./
COPY backend/api/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY backend/api .

# Build the binary from the cmd folder
# This assumes all files in cmd/ are part of package main
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/.env .
EXPOSE 5000
CMD ["./main"]


# FROM golang:1.23.0 AS builder 
# WORKDIR /usr/src/app 
# COPY backend/api/go.mod ./
# COPY backend/api/go.sum ./
# RUN go mod download 
# COPY backend/api . 
# RUN CGO_ENABLED=0 GDOS=linux go build -o main . 

# FROM alpine:latest 
# RUN apk --no-cache add ca-certificates 
# WORKDIR /root/
# COPY --from=builder /usr/src/app/main . 
# COPY --from=builder /usr/src/app/.env . 
# EXPOSE 5000
# CMD [ "./main" ]






