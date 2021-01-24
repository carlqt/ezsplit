class InvitesController < ApplicationController
  before_action :set_invite, only: [:show, :update, :destroy]

  # GET /invites/1
  # GET /invites/1.json
  def show
  end

  # POST /invites
  # POST /invites.json
  def create
    @invite = Invite.new(invite_params)

    if @invite.save
      render :show, status: :created, location: @invite
    else
      render json: @invite.errors, status: :unprocessable_entity
    end
  end

  private

  # Use callbacks to share common setup or constraints between actions.
  def set_invite
    @invite = Invite.find(params[:id])
  end

  # Never trust parameters from the scary internet, only allow the white list through.
  def invite_params
    params.require(:invite).permit(:profile, :group, :token)
  end
end
