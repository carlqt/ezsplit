class AddAccountToReceipts < ActiveRecord::Migration[5.2]
  def change
    add_reference :receipts, :account, foreign_key: true
  end
end
