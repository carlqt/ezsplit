class CreateItems < ActiveRecord::Migration[5.2]
  def change
    create_table :items do |t|
      t.references :receipt
      t.string :name
      t.integer :quantity, default: 0
      t.bigint :price_cents, default: 0

      t.timestamps
    end
  end
end
