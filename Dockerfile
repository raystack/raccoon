FROM golang:latest

WORKDIR /app
COPY . .
RUN go mod download && \
    make build

EXPOSE 8080

# Command to run the executable
CMD ["./out/raccoon"]
