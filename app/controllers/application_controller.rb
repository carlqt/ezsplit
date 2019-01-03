class ApplicationController < ActionController::API
  before_action :current_user_info

  def current_user_info
    JWT.decode(bearer_token, Rails.application.credentials.secret_key_base, true).first
  rescue JWT::DecodeError
    {}
  end

  private

  def bearer_token
    pattern = /^Bearer /
    header  = request.headers['Authorization']
    header.gsub(pattern, '') if header && header.match(pattern)
  end
end
