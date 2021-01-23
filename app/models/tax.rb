class Tax < ApplicationRecord
  belongs_to :receipt

  def self.multiplier
    self.pluck(:rate).reduce(1) do |mem, rate|
      mem = (1 + (rate.to_d / 100)) * mem
    end
  end
end
