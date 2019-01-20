class ReceiptShare < ApplicationRecord
  belongs_to :receipt
  belongs_to :item
  belongs_to :profile

  validates :receipt, presence: true
end
