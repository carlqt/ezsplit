table "accounts" {
  schema = schema.public
  column "id" {
    type = bigint
    identity {
      generated = ALWAYS
    }
  }
  column "username" {
    null    = false
    type    = varchar(25)
  }
  column "password" {
    null    = false
    type    = varchar(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
}

table "users" {
  schema = schema.public
  column "id" {
    type = bigint
    identity {
      generated = ALWAYS
    }
  }
  column "name" {
    null    = false
    type    = varchar(25)
    default = ""
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("now()")
  }
  column "account_id" {
    null = true
    type = bigint  
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "account_id" {
    columns     = [column.account_id]
    ref_columns = [table.accounts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_users_on_name" {
    columns = [column.name]
    unique = true
    where = "account_id IS NOT NULL"
  }
  index "idx_users_on_account_id" {
    columns = [column.account_id]
  }
}

table "receipts" {
  schema = schema.public
  column "id" {
    type = bigint
    identity {
      generated = ALWAYS
    }
  }
  column "user_id" {
    type = bigint
  }
  column "description" {
    null = false
    default = ""
    type = text
  }
  column "url_slug" {
    default = ""
    type = varchar(255)
  }
  column "total" {
    null = true
    type = int
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  check "non_empty_description" {
    expr = "description != ''"
  }
}

table "items" {
  schema = schema.public
  column "id" {
    type = bigint
    identity {
      generated = ALWAYS
    }
  }
  column "receipt_id" {
    type = bigint
  }
  column "price" {
    null = false
    type = int
    default = 0
  }
  column "name" {
    null = true
    type = character_varying(128)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "receipt_id" {
    columns     = [column.receipt_id]
    ref_columns = [table.receipts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "FK_items_receipt_id" {
    columns = [column.receipt_id]
  }
}

table "user_orders" {
  schema = schema.public
  column "id" {
    type = bigint
    identity {
      generated = ALWAYS
    }
  }
  column "user_id" {
    null = true
    type = bigint
  }
  column "item_id" {
    null = true
    type = bigint
  }
  column "status" {
    null = true
    type = character_varying(128)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "item_id" {
    columns     = [column.item_id]
    ref_columns = [table.items.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}

schema "public" {
}
