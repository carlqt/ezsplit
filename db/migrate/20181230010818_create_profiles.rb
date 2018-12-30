class CreateProfiles < ActiveRecord::Migration[5.2]
  def change
    create_table :profiles do |t|
      t.references :account
      t.references :group

      t.timestamps
    end
  end
end
