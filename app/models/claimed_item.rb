class ClaimedItem < ApplicationRecord
  belongs_to :item
  belongs_to :profile

  validates :receipt, presence: true
end
