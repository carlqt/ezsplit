json.extract! invite, :id, :profile, :group, :token, :expired_at, :created_at, :updated_at
json.url invite_url(invite, format: :json)
