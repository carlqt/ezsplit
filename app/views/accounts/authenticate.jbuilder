json.data do
  json.token @jwt_token
  json.id @account.id
  json.email @account.email
  json.profiles(@profiles) do |profile|
    json.(
      profile,
      :id,
      :group_id,
    )
  end

  json.groups(@groups) do |group|
    json.(
      group,
      :id,
      :name,
    )
  end
end
