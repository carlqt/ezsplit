class Group < ApplicationRecord
  has_many :profiles

  # show only group you are invited
end
