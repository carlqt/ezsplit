class Item < ApplicationRecord
  belongs_to :receipt
  has_many :claimed_items

  validates :receipt, presence: true

  scope :total_price, -> { sum("price_cents * quantity") }

  def total_price
    quantity.nil? ? 0 : (price_cents / quantity)
  end

  def display_price
    price_cents.to_d / 100
  end
end
