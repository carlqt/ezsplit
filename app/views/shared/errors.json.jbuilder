json.errors(errors) do |attribute, detail|
  json.partial! 'shared/error', attribute: attribute, detail: detail
end
