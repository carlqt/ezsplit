class Profile < ApplicationRecord
  has_many :receipts
  has_many :claimed_items
  has_many :items, through: :claimed_items
  belongs_to :account, optional: true
  belongs_to :group

  validates :account, uniqueness: { scope: :group_id }, if: :account_exists?

  delegate :email, to: :account

  def balance
    receipts.sum(:price_cents) - claimed_items.includes(:item).sum(&:price_cents)
  end

  private

  def account_exists?
    account.present?
  end
end
