FROM golang:1.16-alpine AS build
LABEL stage=builder
ENV GO111MODULE=on GOOS=linux GOARCH=amd64
WORKDIR /sequencer
COPY . .
RUN go mod download && mkdir ./bin/ && go build -o ./bin/seq-server ./sequencer/cmd/main.go

FROM alpine:latest  
WORKDIR /sequencer
COPY --from=build /sequencer/bin/seq-server .
ENTRYPOINT [ "./seq-server" ]