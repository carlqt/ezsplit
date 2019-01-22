json.(
  item,
  :id,
  :name,
  :quantity,
  :price,
)
unless item.claimed_items.empty?
  json.claimed_by item.claimed_items.map{ |i| i.profile.email }.to_sentence
end

json.checked item.claimed_items.any?{|i| i.profile == @current_profile}
