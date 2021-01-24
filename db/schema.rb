# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2021_01_24_025930) do

  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "accounts", force: :cascade do |t|
    t.string "email"
    t.string "password_digest"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["email"], name: "index_accounts_on_email", unique: true
  end

  create_table "groups", force: :cascade do |t|
    t.string "name"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "invite_token"
    t.index ["invite_token"], name: "index_groups_on_invite_token", unique: true
  end

  create_table "invites", force: :cascade do |t|
    t.bigint "profile_id"
    t.bigint "group_id"
    t.string "token"
    t.datetime "expired_at"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["group_id"], name: "index_invites_on_group_id"
    t.index ["profile_id"], name: "index_invites_on_profile_id"
  end

  create_table "items", force: :cascade do |t|
    t.bigint "receipt_id"
    t.string "name"
    t.integer "quantity", default: 0
    t.bigint "price_cents", default: 0
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["receipt_id"], name: "index_items_on_receipt_id"
  end

  create_table "profile_items", force: :cascade do |t|
    t.bigint "item_id"
    t.bigint "profile_id"
    t.integer "status", default: 0
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["item_id"], name: "index_profile_items_on_item_id"
    t.index ["profile_id"], name: "index_profile_items_on_profile_id"
  end

  create_table "profiles", force: :cascade do |t|
    t.bigint "account_id"
    t.bigint "group_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["account_id"], name: "index_profiles_on_account_id"
    t.index ["group_id"], name: "index_profiles_on_group_id"
  end

  create_table "receipts", force: :cascade do |t|
    t.bigint "profile_id"
    t.text "description"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.bigint "price_cents", default: 0
    t.bigint "account_id"
    t.index ["account_id"], name: "index_receipts_on_account_id"
    t.index ["profile_id"], name: "index_receipts_on_profile_id"
  end

  create_table "taxes", force: :cascade do |t|
    t.string "name"
    t.decimal "rate", precision: 5, scale: 2
    t.bigint "receipt_id"
    t.index ["receipt_id"], name: "index_taxes_on_receipt_id"
  end

  add_foreign_key "receipts", "accounts"
end
