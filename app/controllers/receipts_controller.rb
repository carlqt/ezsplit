class ReceiptsController < ApplicationController
  before_action :authenticate_account!, only: %i[new create]

  def new
    @receipt = current_account.receipts.new
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
    @receipt = current_account.receipts.new do |r|
      r.description = receipt_params[:description]
      r.items = receipt_items
      r.taxes = receipt_taxes
    end

    if @receipt.save
      render status: :created
    else
      render status: :unprocessable_entity
    end
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
        items: [:name, :price],
        taxes: [:name, :rate],
      )
  end

  def receipt_taxes
    # optional

    receipt_params[:taxes].map { |tax| Tax.new(tax) }
  end

  def receipt_items
    receipt_params[:items].map do |item|
      Item.new(
        price_cents: item[:price].to_i * 100,
        name: item[:name],
      )
    end
  end
end
