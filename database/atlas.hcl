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
  column "date_created" {
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


table "answers" {
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
  foreign_key "answers_foreign_key_poll" {
    columns = [column.poll_id]
    ref_columns = [table.polls.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "answers_unique_constraint_id" {
    columns = [
      column.id,
    ]
    unique = true
  }
}


table "user_answers" {
   schema = schema.public
   column "user_id" {
    type    = varchar(255)
  }
  column "answer_id" {
    type = uuid
  }
  foreign_key "user_answers_fk_user_id" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  foreign_key "user_answers_fk_answer_id" {
    columns = [column.answer_id]
    ref_columns = [table.answers.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "user_answers_answer_id" {
    columns = [
      column.answer_id,
    ]
  }
}
