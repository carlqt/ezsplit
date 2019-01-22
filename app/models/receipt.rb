class Receipt < ApplicationRecord
  has_many :items
  has_many :claimed_items
end
