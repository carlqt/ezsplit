json.data do
  json.members(@members) do |member|
    json.profile_id member.id
    json.email member.email
  end
end
