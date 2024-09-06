FROM debian:bookworm-slim
WORKDIR /app
COPY raccoon .
ENTRYPOINT [ "/app/raccoon" ]
