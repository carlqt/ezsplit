class Group < ApplicationRecord
  has_many :profiles
  has_many :account_members, through: :profiles, source: :account

  def receipts
    Receipt.where(profile_id: profiles)
  end
end
