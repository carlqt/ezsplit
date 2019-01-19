json.data do
  json.(
    @group,
    :id,
    :name,
  )

  json.members(@members) do |member|
    json.profile_id member.id
    json.email member.email
    json.balance member.balance
  end
end
