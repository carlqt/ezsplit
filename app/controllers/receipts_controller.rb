class ReceiptsController < ApplicationController
  before_action :authenticate_account!

  def new
    @profiles = accounts_in_group.select(:id, :email)
  end

  def create
  end

  private

  def accounts_in_group
   Account 
     .joins("LEFT JOIN profiles ON profiles.group_id = #{params[:home_id]} AND profiles.account_id = accounts.id")
  end

  def receipt_params
    params.require(:receipt)
      .permit(
        :description,
        :total,
        members: [:id],
        items: [:name, :quantity, :price],
      )
  end
end
