# syntax=docker/dockerfile:1.4
FROM --platform=$BUILDPLATFORM golang:1.25.1-alpine as base

WORKDIR /app

COPY . .
RUN go mod download

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /usr/local/bin/kustomize-job-hasher ./

FROM scratch

COPY --from=base /usr/local/bin/kustomize-job-hasher .

ENTRYPOINT ["./kustomize-job-hasher"]