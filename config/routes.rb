Rails.application.routes.draw do
  resources :invites, only: %i[show create] do
    post :confirm, on: :member
  end

  resources :home, only: [:show, :index], controller: :home do
    resources :receipts, only: [:show, :new, :create, :index] do
      post :claim_items, on: :member
    end

    get :members, on: :member
  end
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
  resources :accounts, only: [:show] do
    post :authenticate, on: :collection
  end

  resources :invite, only: [:show, :create], param: :token
  resource :account, only: [:show], controller: :account
end
