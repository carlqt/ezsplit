# EZSplit

A project app that makes it easy to split, track and manage bills in a group.

Have you ever found yourself shouldering the bill of your friends at a restaurant and forced to calculate your group's invidividual expenses? Adding insult to injury is when the bill has tax that's excluded from the items.

The goal of this app is to move the responsibility of tracking payments to your friends since they have the knowledge of the items they bought.

# Roadmap
[Coming Soon]

# Getting Started
## Introduction

The backend uses `Go` version 1.21.6. It's not required to download the Go runtime since the project uses docker but you're free to do so if it helps with productivity.

The project uses Graphql for the API. To learn more about the concept, visit https://graphql.org/.

The project uses the gqlgen to help with organizing the queries/mutations and generating the graphql schema. Please check out the library to learn more about the tool's style and convention at https://gqlgen.com/getting-started/


## Requirements

- [Docker](https://docs.docker.com/engine/install/)
- [Task](https://taskfile.dev/installation/)

## Running the app

If Docker and Task are installed, you can start the server using the command

```
task up
```
Verify if the server running:
```
task ping
```

After verifying that the server is running, you need to apply the database schema. Applying the schema will create all the necessary tables needed for development:
```
task apply:db:schema
```

This project uses atlasgo tool to implement a declarative style of db management as opposed to the common versioned style. Read https://atlasgo.io/concepts/declarative-vs-versioned to learn more.

To view the database, the project ships with `pgweb` client which is found on `localhost:8082`

To see more available commands, you can run:
```
task list
```
Or checkout `Taskfile.yml`
