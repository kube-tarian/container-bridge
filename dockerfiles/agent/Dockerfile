FROM golang:1.19 AS builder
WORKDIR /
COPY ./ ./

RUN go mod download
RUN go build -o ./build/agent agent/main.go

FROM scratch
COPY --from=builder ./build/agent agent

USER 65532:65532
ENTRYPOINT ["./agent"]
