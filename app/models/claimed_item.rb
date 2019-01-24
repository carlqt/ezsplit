class ClaimedItem < ApplicationRecord
  belongs_to :item
  belongs_to :profile

  validates :profile, uniqueness: { scope: :item_id}

  scope :on_receipt, -> (receipt) { joins(item: :receipt).where(items: { receipt: receipt })  }

  def price
    item.price / item.claimed_items.count
  end
end
