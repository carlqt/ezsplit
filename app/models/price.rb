class Price
  attr_reader :value_in_cents

  def initialize(value_in_cents)
    @value_in_cents = value_in_cents
    @value = value_in_cents.to_d / 100
  end

  def display_value
    "%.2f" % @value.round(2)
  end
end
