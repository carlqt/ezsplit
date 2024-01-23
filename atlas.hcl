// Define an environment named "local"
env "local" {
  // Declare where the schema definition resides.
  // Also supported:
  //   src = "./dir/with/schema"
  //   src = ["multi.hcl", "file.hcl"]
  src = "./schema.hcl"

  // Define the URL of the database which is managed in
  // this environment.
  url = "postgres://postgres:postgres@localhost:5432/ezsplit_dev?sslmode=disable"

  // The schemas in the database that are managed by Atlas.
  schemas = ["public"]
}