FROM golang:1.14

WORKDIR /app
RUN apt-get update && apt-get install unzip  --no-install-recommends --assume-yes
RUN PROTOC_ZIP=protoc-3.14.0-linux-x86_64.zip && \
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP && \
unzip -o $PROTOC_ZIP -d /usr/local bin/protoc && \
unzip -o $PROTOC_ZIP -d /usr/local 'include/*' && \
rm -f $PROTOC_ZIP
COPY . .
RUN make install-protoc
RUN make generate-proto
RUN make update-deps
RUN make compile

EXPOSE 8080

# Command to run the executable
CMD ["./out/raccoon"]
