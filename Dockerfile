FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:latest as builder

RUN mkdir /BAKTA-Web-Cleaner
WORKDIR /BAKTA-Web-Cleaner
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o BaktaCleaner .

FROM scratch
ARG GITHUB_SHA
ENV GITHUB_SHA=$GITHUB_SHA
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /BAKTA-Web-Cleaner/BaktaCleaner .

ENTRYPOINT [ "/BaktaBackend" ]