FROM golang:1.22.4

WORKDIR /app
RUN apt-get update && apt-get install unzip  --no-install-recommends --assume-yes
RUN PROTOC_ZIP=protoc-3.17.3-linux-x86_64.zip && \
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/$PROTOC_ZIP && \
    unzip -o $PROTOC_ZIP -d /usr/local bin/protoc && \
    unzip -o $PROTOC_ZIP -d /usr/local 'include/*' && \
    rm -f $PROTOC_ZIP
COPY . .
RUN make build


FROM debian:bookworm-slim
WORKDIR /app
COPY --from=0 /app/raccoon ./raccoon
CMD ["./raccoon", "server"]
