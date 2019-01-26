require 'test_helper'

class InviteControllerTest < ActionDispatch::IntegrationTest
  test "should get show" do
    get invite_show_url
    assert_response :success
  end

  test "should get create" do
    get invite_create_url
    assert_response :success
  end

end
