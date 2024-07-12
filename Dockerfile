# syntax=docker/dockerfile:1

FROM golang:1.22.4-alpine3.19 as base

LABEL maintainer="carltablante"

WORKDIR /app

RUN apk add curl

RUN curl -sSf https://atlasgo.sh | sh

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
