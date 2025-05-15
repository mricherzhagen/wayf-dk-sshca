#!/bin/bash

cd "$(dirname "$0")"

#ssh-keygen -t ed25519 -f ./ca_ed25519 -N ""

F=ca_ed25519
PRIVATE=${F}.key.pem
PUBLIC=${F}.pub.pem
PUBLIC_SSH=${F}.pub.ssh
PRIVATE_PKCS8=${F}.key.pkcs8
openssl genpkey -algorithm ed25519 -out $PRIVATE
openssl pkey -in $PRIVATE -pubout -out $PUBLIC


SSHPK_CONTAINER=$(podman build -q . -f Dockerfile_sshpk) 
cat $PUBLIC | podman run -i --rm $SSHPK_CONTAINER -T pem -t ssh > $PUBLIC_SSH

cat $PRIVATE | podman run -i --rm $SSHPK_CONTAINER -p -T pem -t pkcs8 > $PRIVATE_PKCS8

echo "cert-authority `cat $PUBLIC_SSH`" > ./authorized_keys
