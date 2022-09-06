FROM golang:1.19 AS builder
WORKDIR /
COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/agent agent/main.go

FROM scratch
COPY --from=builder ./build/agent agent

USER 65532:65532
ENTRYPOINT ["./agent"]
