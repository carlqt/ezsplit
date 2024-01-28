FROM golang:1.21.6-bullseye as base

LABEL maintainer="carltablante"

WORKDIR /app

RUN curl -sSf https://atlasgo.sh | sh

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM base as builder

RUN go build -o build/graphql_playground ./cmd/graphql_playground/main.go

COPY --from=stage /app/build/graphql_playground /app

# Add a script to be executed every time the container starts.
# COPY entrypoint.sh /usr/bin/
# RUN chmod +x /usr/bin/entrypoint.sh
# ENTRYPOINT ["entrypoint.sh"]
EXPOSE 8080

# Start the main process.
CMD ["./graphql_playground"]
