class Receipt < ApplicationRecord
  has_many :items
  has_many :claimed_items
  has_many :taxes

  belongs_to :account
end
