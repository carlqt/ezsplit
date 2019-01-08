class Account < ApplicationRecord
  has_secure_password
  has_many :profiles
  has_many :groups, through: :profiles
end
