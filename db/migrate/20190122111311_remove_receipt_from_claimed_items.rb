class RemoveReceiptFromClaimedItems < ActiveRecord::Migration[5.2]
  def change
    remove_column :claimed_items, :receipt_id, :bigint
  end
end
