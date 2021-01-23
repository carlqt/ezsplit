require 'rails_helper'

describe Receipt do
  describe '#total' do
    subject { receipt.total }

    let(:receipt) { build(:receipt) }
  end
end
