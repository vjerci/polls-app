schema "public" {}

// varchars for ids are used to simplify app login it should probabbly be uuid
table "users" {
  schema = schema.public
  column "id" {
    type = varchar(255)
  }
  column "name" {
    type = varchar(255)
  }
  primary_key {
    columns = [
      column.id,
    ]
  }
  index "idx_unique_constraint" {
    columns = [
      column.id,
    ]
    unique = true
  }
}