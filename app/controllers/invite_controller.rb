class InviteController < ApplicationController
  before_action :find_group

  def show
  end

  def create
    @account = Account.new(account_params) do |account|
      account.profiles = [Profile.new(account: account, group: @group)]
    end

    if @account.save
      payload = @account.to_json(only: [:id, :email])
      @jwt_token = JWT.encode(payload, Rails.application.credentials.secret_key_base, "HS256")
      @profiles = @account.profiles
      @groups = @account.groups
    else
      render json: { message: @account.errors.full_messages[0] }, status: :unauthorized
    end
  end

  private

  def find_group
    @group = Group.find_by(invite_token: params[:token])

    if @group.nil?
      render json: { message: "page not found" }, status: 404
    end
  end

  def account_params
    params[:invite].require(:account)
      .permit(
        :email,
        :password,
        :password_confirmation
      )
  end
end
