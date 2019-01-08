class Receipt < ApplicationRecord
  has_many :items
  has_many :receipt_shares
end
