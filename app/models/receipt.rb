class Receipt < ApplicationRecord
  has_many :items
  has_many :claimed_items

  belongs_to :account
end
