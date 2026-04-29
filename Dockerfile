FROM golang:1.26.2-bookworm AS base

LABEL maintainer="carltablante"

WORKDIR /app

RUN apt-get update && apt-get install -y curl

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
