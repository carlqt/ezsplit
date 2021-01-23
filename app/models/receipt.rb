class Receipt < ApplicationRecord
  has_many :items
  has_many :claimed_items
  has_many :taxes

  belongs_to :account

  def total
    @total ||= Price.new(items_price * tax_multiplier)
  end

  def tax_multiplier
    taxes.multiplier
  end

  def items_price
    items.total_price
  end
end
