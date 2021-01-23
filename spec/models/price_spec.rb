require 'rails_helper'

describe Price do
  describe '#display_value' do
    subject { price.display_value }

    let(:price) { Price.new(4878) }

    it { is_expected.to eq '48.78' }
  end
end

