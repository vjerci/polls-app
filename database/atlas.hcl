schema "public" {}

table "users" {
  schema = schema.public
  column "id" {
    type = varchar(255)
  }
  column "group_id" {
    type = varchar(255)
  }
  column "name" {
    type = varchar(255)
  }
  primary_key {
    columns = [
      column.id,
      column.group_id
    ]
  }
  index "idx_unique_constraint" {
    columns = [
      column.id,
      column.group_id,
    ]
    unique = true
  }
}