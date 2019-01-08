class CreateReceiptShares < ActiveRecord::Migration[5.2]
  def change
    create_table :receipt_shares do |t|
      t.references :item
      t.references :receipt
      t.references :profile

      t.timestamps
    end
  end
end
