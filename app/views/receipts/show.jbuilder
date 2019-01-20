json.data do
  json.(
    @receipt,
    :id,
    :description,
    :price,
    :profile_id,
  )

  json.items(@items) do |item|
    json.(
      item,
      :id,
      :name,
      :quantity,
      :price,
    )
    if item.receipt_share
      json.claimed_by do
        json.id item.receipt_share.profile_id
        json.email item.receipt_share.profile.email
      end
    end
    json.checked item.receipt_share && (item.receipt_share.profile_id == @current_profile.id)
  end
end
