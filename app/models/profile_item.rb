class ProfileItem < ApplicationRecord
  belongs_to :item
  belongs_to :profile

  scope :on_receipt, -> (receipt) { joins(item: :receipt).where(items: { receipt: receipt })  }

  delegate :price, to: :item

  enum status: %i(awaiting_payment paid)
end
