FROM golang:latest

WORKDIR /app
COPY . .
RUN make copy-config && make compile

EXPOSE 8080

# Command to run the executable
CMD ["./out/raccoon"]
