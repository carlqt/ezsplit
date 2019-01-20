class Item < ApplicationRecord
  belongs_to :receipt
  has_one :receipt_share
end
