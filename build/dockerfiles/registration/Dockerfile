FROM golang:1.16-alpine AS build
LABEL stage=builder
ENV GO111MODULE=on GOOS=linux GOARCH=amd64
WORKDIR /registration
COPY . .
RUN go mod download && mkdir ./bin/ && go build -o ./bin/reg-server ./registration/cmd/main.go

FROM alpine:latest  
WORKDIR /registration
COPY --from=build /registration/bin/reg-server .
ENTRYPOINT [ "./reg-server" ]