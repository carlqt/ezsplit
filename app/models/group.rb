class Group < ApplicationRecord
  has_many :profiles
  has_many :account_members, through: :profiles, source: :account
end
