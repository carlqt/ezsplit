class AccountController < ApplicationController
  before_action :authenticate_account!

  def show
    @account = current_account
    @groups = current_account.groups
    @profiles = current_account.profiles
  end
end
