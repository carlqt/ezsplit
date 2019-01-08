class Profile < ApplicationRecord
  has_many :receipts
  has_many :receipt_shares
  belongs_to :account
  belongs_to :group

  validates :account, uniqueness: { scope: :group_id }

  def balance
    receipts.sum(:price)
  end
end
