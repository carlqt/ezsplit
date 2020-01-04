class AddPriceToReceipts < ActiveRecord::Migration[5.2]
  def change
    add_column :receipts, :price_cents, :bigint, default: 0
  end
end
