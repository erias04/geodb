FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /
COPY . .
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./web-app ./main.go

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=build . .
EXPOSE 80
ENTRYPOINT ./web-app --port 80