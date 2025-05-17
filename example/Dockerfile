FROM archlinux:latest AS build

WORKDIR /app

RUN pacman -Syu  --noconfirm go gcc

COPY go.mod go.sum ./

RUN go mod download

COPY *.go .

COPY cmd/ca.template cmd/

COPY cmd/pubkey_to_label.go cmd/
RUN CGO_ENABLED=1 GOOS=linux go build -o ./pubkey_to_label cmd/pubkey_to_label.go 

COPY cmd/sshca_example.go cmd/
RUN CGO_ENABLED=1 GOOS=linux go build -o ./sshca_example cmd/sshca_example.go 


FROM archlinux:latest

# Set the working directory
WORKDIR /app

RUN pacman -Syu --noconfirm softhsm

RUN softhsm2-util --init-token --slot 0 --label "sshca_token" --pin 1234 --so-pin 123456

ADD /example/ca_ed25519.key.pkcs8 ca_ed25519.key.pkcs8

ADD /example/ca_ed25519.pub.ssh ca_ed25519.pub.ssh

COPY --from=build /app/pubkey_to_label .

RUN softhsm2-util --import ca_ed25519.key.pkcs8 --pin 1234 --token "sshca_token" --id 9c --label `./pubkey_to_label ca_ed25519.pub.ssh`


COPY --from=build /app/sshca_example .

# Set the entrypoint command
ENTRYPOINT ["/app/sshca_example", "./ca_ed25519.pub.ssh"]
