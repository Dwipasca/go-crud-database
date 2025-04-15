# base image golang
FROM golang:1.23.4

# set working directory
WORKDIR /app

# copy all the files to container
COPY . .

# download dependencies
RUN go mod tidy

# build binary
RUN go build -o main ./cmd

# run application
CMD ["./main"]