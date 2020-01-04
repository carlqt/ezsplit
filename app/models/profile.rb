class Profile < ApplicationRecord
  has_many :receipts
  has_many :claimed_items
  has_many :items, through: :claimed_items
  belongs_to :account
  belongs_to :group

  validates :account, uniqueness: { scope: :group_id }

  delegate :email, to: :account

  def balance
    receipts.sum(:price_cents) - claimed_items.includes(:item).sum(&:price_cents)
  end
end
