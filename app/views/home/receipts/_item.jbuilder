json.(
  item,
  :id,
  :name,
  :quantity,
  :price,
)

json.claimed_by(item.claimed_items) do |claimed_item|
  json.id claimed_item.profile_id
  json.email claimed_item.profile.email
end

json.checked item.claimed_items.any?{|i| i.profile == @current_profile}
