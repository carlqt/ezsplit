class Item < ApplicationRecord
  belongs_to :receipt
  has_many :claimed_items

  validates :receipt, presence: true

  scope :total_price, -> { sum("price_cents") }

  def price
    @price ||= Price.new(price_cents)
  end
end
