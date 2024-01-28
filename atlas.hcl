variable "url" {
  type = string
  default = getenv("DATABASE_URL")
}

variable "port" {
  type = string
  default = getenv("DB_PORT")
}

variable "host" {
  type = string
  default = getenv("DB_HOST")
}

variable "user" {
  type = string
  default = getenv("DB_USER")
}

variable "name" {
  type = string
  default = getenv("DB_NAME")
}

variable "password" {
  type = string
  default = getenv("DB_PASSWORD")
}

// Define an environment named "local"
env "local" {
  // Declare where the schema definition resides.
  // Also supported:
  //   src = "./dir/with/schema"
  //   src = ["multi.hcl", "file.hcl"]
  src = "./schema.hcl"

  // Define the URL of the database which is managed in
  // this environment.
  url = "postgres://${var.user}:${var.password}@${var.host}:${var.port}/${var.name}?sslmode=disable"

  // The schemas in the database that are managed by Atlas.
  schemas = ["public"]
}