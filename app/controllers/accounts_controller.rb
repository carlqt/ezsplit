class AccountsController < ApplicationController
  def authenticate
    account = Account.find_by(email: params[:email])

    if account.authenticate(params[:password])
      payload = account.to_json(only: [:id, :email])
      @jwt_token = JWT.encode(payload, Rails.application.credentials.secret_key_base, "HS256")
    else
      # error
    end
  end
end
