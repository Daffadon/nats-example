#!/bin/bash

if [[ $(basename "$PWD") == "bin" ]]; then
  echo "Run this from project root!!"
  exit 1
fi
# check nsc installed or not
echo "Checking nsc package....."
if ! command -v nsc &>/dev/null; then
  # install nsc
  echo "nsc package is not exist. Installing nsc..."
  curl -L https://github.com/nats-io/nsc/releases/latest/download/nsc-linux-amd64.zip -o nsc.zip
  unzip nsc.zip -d /usr/local/bin
  chmod +x /usr/local/bin/nsc
  rm nsc.zip
fi
echo "nsc has been installed...."

echo "Check if the NATS_OPERATOR exist...."
if ! nsc describe operator NATS_OPERATOR &>/dev/null; then
  echo "Operator NATS_OPERATOR does not exist. Creating..."
  echo "Creating nsc operator with system account and regular user"
  nsc add operator --generate-signing-key --sys --name NATS_OPERATOR
  nsc edit operator --require-signing-keys --account-jwt-server-url "nats://0.0.0.0:4221"
  nsc add account ADMIN_ROLE
  nsc edit account ADMIN_ROLE --sk generate
  nsc add user --account ADMIN_ROLE --name admin
  # change regular account to have js access
  nsc edit account -n ADMIN_ROLE --js-mem-storage $((1024 * 1024 * 1024)) --js-disk-storage $((50 * 1024 * 1024 * 1024)) --js-streams 10 --js-consumer 10
fi
else
  echo "Operator NATS_OPERATOR already exists. Skipping creation."
  nsc select operator NATS_OPERATOR
fi
# generate sys and user regular account creds
echo "generate sys creds to config/creds/sys.creds"
nsc select account ADMIN_ROLE
nsc generate creds -n admin >config/creds/admin.creds

echo "generate user regular creds to config/creds/sys.creds"
nsc select account SYS
nsc generate creds -n sys >config/creds/sys.creds
# generate resolver.conf

echo "generate resolver with nats-resolver base"
nsc generate config --nats-resolver >config/creds/resolver.conf

# install docker & docker compose
echo "Checking Docker installation..."
if ! command -v docker &>/dev/null; then
  echo "Docker is not installed. Installing Docker..."
  curl -fsSL https://get.docker.com -o get-docker.sh
  sh get-docker.sh
  rm get-docker.sh
else
  echo "Docker is already installed."
fi

echo "Checking Docker Compose installation..."
if ! command -v docker-compose &>/dev/null; then
  echo "Docker Compose is not installed. Installing Docker Compose..."
  curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
  chmod +x /usr/local/bin/docker-compose
else
  echo "Docker Compose is already installed."
fi
