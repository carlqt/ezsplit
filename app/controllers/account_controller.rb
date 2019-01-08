class AccountController < ApplicationController
  before_action :authenticate_account!

  def show
    @groups = current_account.groups
    @profiles = current_account.profiles
  end
end
