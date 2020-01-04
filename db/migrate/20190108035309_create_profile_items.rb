class CreateProfileItems < ActiveRecord::Migration[5.2]
  def change
    create_table :profile_items do |t|
      t.references :item
      t.references :profile
      t.integer :status, default: 0

      t.timestamps
    end
  end
end
