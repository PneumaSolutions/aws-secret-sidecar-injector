FROM golang:1.16.6-alpine3.14@sha256:bc2db47c5f4a682f1315e0d484811d65bf094d3bcd824459b170714c91656190 AS build
ENV CGO_ENABLED=0
WORKDIR /src/aws-secrets-manager
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /app -v ./cmd/aws-secrets-manager

FROM scratch
ENV MOUNT_POINT=/tmp
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /.
ENTRYPOINT ["/app"]
