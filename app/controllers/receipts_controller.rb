class ReceiptsController < ApplicationController
  before_action :authenticate_account!

  def index
    @receipts = current_group.receipts
  end

  def show
    @current_profile = Profile.find_by(account_id: account_payload[:id], group_id: params[:home_id])
    @receipt = current_profile.receipts.find(params[:id])
    @items = @receipt.items.includes(:claimed_items)
  end

  def new
    @profiles = accounts_in_group.select(:id, :email)
  end

  def claim
    @receipt = current_group.receipts.find(params[:id])
  end

  def create
    @receipt = current_profile.receipts.new do |r|
      r.description = receipt_params[:description]
      r.price = receipt_params[:total]
      r.items = receipt_items
      r.receipt_shares = receipt_shares
    end

    @receipt.save
  end

  def edit
    @receipt = current_account.receipts.find(params[:id])
  end

  private

  def current_group
    current_profile.group
  end

  def current_profile
    @current_profile ||= Profile.find_by(account_id: account_payload[:id], group_id: params[:home_id])
  end

  def accounts_in_group
   Account 
     .joins("LEFT JOIN profiles ON profiles.group_id = #{params[:home_id]} AND profiles.account_id = accounts.id")
  end

  def receipt_params
    params.require(:receipt)
      .permit(
        :description,
        :total,
        members: [:profile_id],
        items: [:name, :quantity, :price],
      )
  end

  def receipt_shares
    receipt_params[:members].map{|m| ReceiptShare.new(m)}
  end

  def receipt_items
    receipt_params[:items].map{ |item| Item.new(item) }
  end
end
