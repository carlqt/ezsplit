require 'rails_helper'

describe Receipt do
  describe '#total' do
    subject { receipt.total.display_value }

    let(:receipt) { create(:receipt_with_items) }

    it { is_expected.to eq "73.80" }
  end
end
