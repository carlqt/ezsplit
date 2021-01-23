json.data do
  json.(
    @receipt,
    :id,
    :description,
    :price,
    :profile_id,
  )

  json.items(@items) do |item|
    json.partial! 'item', item: item
  end
end
