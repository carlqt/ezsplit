class RenameReceiptShareToClaimedItems < ActiveRecord::Migration[5.2]
  def change
    rename_table :receipt_shares, :claimed_items
  end
end
