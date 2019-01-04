class HomeController < ApplicationController
  before_action :authenticate_account!

  def index
    @groups = profile_groups
  end

  def show
    group = Group.find(params[:id])
    @profiles = group.profiles
  end

  private

  def profile_groups
    Group
      .joins("LEFT JOIN profiles ON profiles.group_id = groups.id AND profiles.account_id = #{account_payload["id"]}")
  end
end
