FROM golang:1.14

WORKDIR /app
COPY . .
RUN make compile

EXPOSE 8080

# Command to run the executable
CMD ["./out/raccoon"]
