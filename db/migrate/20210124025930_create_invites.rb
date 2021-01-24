class CreateInvites < ActiveRecord::Migration[5.2]
  def change
    create_table :invites do |t|
      t.references :profile
      t.references :group
      t.string :token
      t.datetime :expired_at

      t.timestamps
    end
  end
end
