schema "public" {}

// varchars are used for id to simplify ux of app login, in real world scenario it should probabbly be uuid
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
  index "users_unique_constraint" {
    columns = [
      column.id,
    ]
    unique = true
  }
}


table "polls" {
  schema = schema.public
  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "name" {
    type = varchar(255)
  }
  column "data_created" {
    type = int
  }
  primary_key {
    columns = [
      column.id,
    ]
  }
  index "polls_unique_constraint" {
    columns = [
      column.id,
    ]
    unique = true
  }
}


table "polls_possible_answers" {
   schema = schema.public
   column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "poll_id" {
    type    = uuid
  }
  column "name" {
    type = varchar(255)
  }
  primary_key {
    columns = [
      column.id,
    ]
  }
  foreign_key "poll_foreign_key" {
    columns = [column.poll_id]
    ref_columns = [table.polls.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "polls_possible_answers_unique_constraint" {
    columns = [
      column.id,
    ]
    unique = true
  }
}