require 'rails_helper'

describe ReceiptsController, type: :request do
  describe '#create' do
    let(:account) { create(:account) }
    let(:payload) { account.to_json(only: %i[id email]) }
    let(:token) { JWT.encode(payload, Rails.application.credentials.secret_key_base, "HS256") }
    let(:headers) { { "Authorization" => "Bearer #{token}"} }

    let(:params) do
      {
        receipt: {
          description: "Aunty Anne's Pretzel",
          items: [
            {
              name: 'Dan Dan Noodle',
              price: 8.8
            },
            {
              name: 'Xiao Long Bao',
              price: 9.3
            },
            {
              name: 'Chili wanton',
              price: 9.5
            },
            {
              name: 'Pork Rice',
              price: 13.8
            },
            {
              name: 'Iverson Noodles',
              price: 10.5
            },
            {
              name: 'Sharona Noodles',
              price: 10.8
            },
          ],
          taxes: [
            {
              name: 'service charge',
              rate: 10,
            },
            {
              name: 'GST',
              rate: 7,
            }
          ]
        }
      }
    end

    before do
      post receipts_path, headers: headers, params: params
    end

    it 'responds with http status created' do
      expect(response).to have_http_status(:created)
    end
  end
end
