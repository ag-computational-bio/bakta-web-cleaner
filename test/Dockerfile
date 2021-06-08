FROM golang:latest as builder
WORKDIR /bin

RUN wget https://dl.min.io/client/mc/release/linux-amd64/mc
RUN chmod 755 /bin/mc

WORKDIR /app_build

COPY . .