class Group < ApplicationRecord
  has_many :profiles
  has_many :account_members, through: :profiles, source: :account
  has_secure_token :invite_token

  def receipts
    Receipt.where(profile_id: profiles)
  end
end
