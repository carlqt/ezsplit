json.data do
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
