FROM golang:1.16-alpine AS build
ENV GO111MODULE=on GOOS=linux GOARCH=amd64
WORKDIR /comunigo
COPY . .
RUN go mod download && mkdir ./bin && go build -o ./bin/peer ./peer/cmd/main.go

FROM alpine:latest  
WORKDIR /comunigo
COPY --from=build /comunigo/bin/peer .
ENTRYPOINT [ "./peer" ]
