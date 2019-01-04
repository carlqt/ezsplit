class ApplicationController < ActionController::API
  rescue_from JWT::DecodeError, with: :unauthorized_access

  def authenticate_account!
    @decoded_token = JWT.decode(bearer_token, Rails.application.credentials.secret_key_base, true)
  end

  def account_payload
    JSON.parse(@decoded_token.first)
  end

  def current_account
    @current_account ||= Account.find(account_payload["id"])
  end

  def unauthorized_access(exception)
    render json: { message: exception }, status: :unauthorized
  end

  private

  def bearer_token
    pattern = /^Bearer /
    header  = request.headers['Authorization']
    header.gsub(pattern, '') if header && header.match(pattern)
  end
end
