class HomeController < ApplicationController
  before_action :authenticate_account!
  before_action :set_group, only: [:show, :members]

  def index
    @groups = profile_groups
  end

  def show
    @members = @group.profiles.joins(:account).select("accounts.email, profiles.id")
  end

  def members
    # @members = group.account_members.select(:email, :id)
    @members = @group.profiles.joins(:account).select("accounts.email, profiles.id")
  end

  private

  def set_group
    @group = Group.find(params[:id])
  end

  def profile_groups
    Group
      .joins("LEFT JOIN profiles ON profiles.group_id = groups.id AND profiles.account_id = #{account_payload["id"]}")
  end
end
