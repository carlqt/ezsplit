class CreateTaxes < ActiveRecord::Migration[5.2]
  def change
    create_table :taxes do |t|
      t.string :name
      t.decimal :rate, scale: 2, precision: 5
      t.references :receipt
    end
  end
end
