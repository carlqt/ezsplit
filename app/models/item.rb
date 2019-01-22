class Item < ApplicationRecord
  belongs_to :receipt
  has_many :claimed_items
end
