class Account < ApplicationRecord
  has_secure_password
  has_many :profiles, dependent: :destroy
  has_many :groups, through: :profiles

  validates :email, uniqueness: true
end
