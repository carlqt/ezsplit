# EZSplit

A project app that makes it easy to split, track and manage bills in a group.

Have you ever found yourself shouldering the bill of your friends at a restaurant and forced to calculate your group's invidividual expenses? Adding insult to injury is when the bill has tax that's excluded from the items.

The goal of this app is to move the responsibility of tracking payments to your friends since they have the knowledge of the items they bought.

# Roadmap
[Coming Soon]

# Getting Started
## Requirements

- [Docker](https://docs.docker.com/engine/install/)
- [Task](https://taskfile.dev/installation/)

## Running the app

If Docker and Task are installed, you can start the server using the command

```
task up
```

Once the server is running, you need to apply the database schema. Applying the schema will create all the necessary tables needed for development:
```
task apply:db:schema
```

This project uses atlasgo tool to implement a declarative style of db management as opposed to the common versioned style. Read https://atlasgo.io/concepts/declarative-vs-versioned to learn more.
