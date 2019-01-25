class AccountsController < ApplicationController
  before_action :get_account

  def authenticate
    if @account.authenticate(params[:password])
      payload = @account.to_json(only: [:id, :email])
      @profiles = @account.profiles
      @groups = @account.groups
      @jwt_token = JWT.encode(payload, Rails.application.credentials.secret_key_base, "HS256")
    else
      # error
      render json: { message: "email or password incorrect" }, status: :unauthorized
    end
  end

  private

  def get_account
    @account = Account.find_by(email: params[:email])

    if @account.nil?
      render json: { message: "incorrect email or password" }, status: :unauthorized
    end
  end
end
