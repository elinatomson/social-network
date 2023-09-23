FROM golang:latest

#setting the working directory inside the container
WORKDIR /app

#coping the Go module files from the root folder and download dependencies
COPY go.mod go.sum ./
RUN go mod download

#coping the entire backend source code from the host into the container
COPY ./back-end/ ./back-end/

#building the Go application
RUN go build -o main ./back-end/cmd/api

#exposing the port the application runs on
EXPOSE 8080

#running the Go application
CMD ["./main"]

