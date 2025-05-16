FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go .

COPY cmd/*.go cmd/ca.template cmd/


RUN GOOS=linux go build -o ./sshca_example cmd/sshca_example.go 


FROM docker.io/andyshinn/alpine-abuild AS softhsm_build

RUN abuild-keygen -n -i && cp /home/builder/.abuild/*.rsa /home/builder/ssh.rsa && cp /home/builder/.abuild/*.rsa.pub /home/builder/ssh.rsa.pub

ADD https://gitlab.alpinelinux.org/alpine/aports/-/raw/76712e926abdeec2eb121a378437301c78f46f5b/main/softhsm/APKBUILD ./APKBUILD

RUN abuilder -r -v

FROM alpine:edge

# Set the working directory
WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=softhsm_build /packages/softhsm.apk .

RUN apk --no-cache add softhsm.apk

RUN softhsm2-util --init-token --slot 0 --label "sshca_token" --pin 1234 --so-pin 123456

ADD /example/ca_ed25519.key.pkcs8 ca_ed25519.key.pkcs8

RUN softhsm2-util --import ca_ed25519.key.pkcs8 --pin 1234 --token "sshca_token" --id 9c --label "sshca_keypair"

# Copy the binary from the build stage
COPY --from=build /app/sshca_example .


# Set the entrypoint command
ENTRYPOINT ["/app/sshca_example"]
