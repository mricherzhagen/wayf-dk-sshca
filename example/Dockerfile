FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go .

COPY cmd/*.go cmd/ca.template cmd/


RUN GOOS=linux go build -o ./sshca_example cmd/sshca_example.go 


FROM alpine:edge

# Set the working directory
WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata softhsm


RUN softhsm2-util --init-token --slot 0 --label "sshca_token" --pin 1234 --so-pin 123456

ADD /example/ca_ed25519.key.pkcs8 ca_ed25519.key.pkcs8

RUN softhsm2-util --import ca_ed25519.key.pkcs8 --pin 1234 --token "sshca_token" --id 9c --label "sshca_keypair"

# Copy the binary from the build stage
COPY --from=build /app/sshca_example .


# Set the entrypoint command
ENTRYPOINT ["/app/sshca_example"]
