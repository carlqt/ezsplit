class ReceiptsController < ApplicationController
  before_action :authenticate_account!

  def index
    @receipts = current_group.receipts.order(:created_at)
  end

  def show
    @current_profile = Profile.find_by(account_id: account_payload[:id], group_id: params[:home_id])
    @receipt = current_group.receipts.find(params[:id])
    @items = @receipt.items.includes(:claimed_items)
  end

  def new
    @profiles = accounts_in_group.select(:id, :email)
  end

  def claim_items
    @receipt = current_group.receipts.find(params[:id])
    items_to_claim = @receipt.items.where(id: params[:items]).ids
    claimed_items = current_profile.claimed_items.on_receipt(@receipt).map(&:item_id)

    # Claiming
    (items_to_claim - claimed_items).each do |id|
      current_profile.claimed_items.create(item_id: id).item
    end

    # Unclaim
    (claimed_items - items_to_claim).each do |id|
      current_profile.claimed_items.find_by(item_id: id).destroy
    end
  end

  def create
    @receipt = current_profile.receipts.new do |r|
      r.description = receipt_params[:description]
      r.price = receipt_params[:total]
      r.items = receipt_items
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
        items: [:name, :quantity, :price],
      )
  end

  def receipt_items
    receipt_params[:items].map{ |item| Item.new(item) }
  end
end
