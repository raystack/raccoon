FROM golang:1.14

WORKDIR /app
RUN apt-get update && apt-get install unzip  --no-install-recommends --assume-yes
COPY . .
RUN make update-deps && make compile

FROM debian:buster-slim
WORKDIR /app
COPY --from=0 /app/raccoon ./raccoon
COPY . .
CMD ["./raccoon"]
