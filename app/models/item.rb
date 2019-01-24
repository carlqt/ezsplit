class Item < ApplicationRecord
  belongs_to :receipt
  has_many :claimed_items

  scope :total_price, -> { sum("price * quantity") }

  def raw_price
    quantity.nil? ? 0 : (price / quantity)
  end
end
