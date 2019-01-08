class ReceiptShare < ApplicationRecord
  belongs_to :receipt

  validates :receipt, presence: true
end
