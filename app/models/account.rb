class Account < ApplicationRecord
  has_secure_password
  has_many :profiles, dependent: :destroy
  has_many :groups, through: :profiles
  has_many :receipts

  validates :email, uniqueness: true
end
