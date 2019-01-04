class HomeController < ApplicationController
  before_action :authenticate_account!

  def index
    @groups = profile_groups
  end

  def show
    @profiles = group.profiles
  end

  def members
    @members = group.account_members.select(:email, :id)
  end

  private

  def group
    @group ||= Group.find(params[:id])
  end

  def profile_groups
    Group
      .joins("LEFT JOIN profiles ON profiles.group_id = groups.id AND profiles.account_id = #{account_payload["id"]}")
  end
end
