FactoryBot.define do
  factory :receipt do
    account
    description { "Din Tai Fung" }

    factory :receipt_with_items do
      description { "Din Tai Fung" }

      after(:create) do |receipt|
        items_params = [
          {
            name: 'Dan Dan Noodle',
            price_cents: 880
          },
          {
            name: 'Xiao Long Bao',
            price_cents: 930
          },
          {
            name: 'Chili wanton',
            price_cents: 950
          },
          {
            name: 'Pork Rice',
            price_cents: 1380
          },
          {
            name: 'Iverson Noodles',
            price_cents: 1050
          },
          {
            name: 'Sharona Noodles',
            price_cents: 1080
          },
        ]

        tax_params = [
          {
            name: 'service charge',
            rate: 10,
          },
          {
            name: 'GST',
            rate: 7,
          }
        ]

        receipt.items = items_params.map { |params| Item.new(params) }
        receipt.taxes = tax_params.map { |params| Tax.new(params) }

        receipt.save
      end
    end
  end
end
