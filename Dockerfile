# syntax=docker/dockerfile:1

FROM golang:1.23.1-alpine3.20 as base

LABEL maintainer="carltablante"

WORKDIR /app

RUN apk add curl

RUN curl -sSf https://atlasgo.sh | sh

# This is only for tests
RUN go install gotest.tools/gotestsum@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Add a script to be executed every time the container starts.
# COPY entrypoint.sh /usr/bin/
# RUN chmod +x /usr/bin/entrypoint.sh
# ENTRYPOINT ["entrypoint.sh"]
EXPOSE 8080

# Start the main process.
CMD ["./graphql_playground"]
