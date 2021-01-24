class Invite < ApplicationRecord
  belongs_to :group
  belongs_to :profile

  has_secure_token :token

  before_create :set_expired_at

  private

  def set_expired_at
    self.expired_at = Time.current + 7.days
  end
end
