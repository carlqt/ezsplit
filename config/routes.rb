Rails.application.routes.draw do
  resources :home, only: [:show, :index], controller: :home do
    resources :receipts, only: [:show, :new, :create, :index, :edit]
    get :members, on: :member
  end
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
  resources :accounts, only: [:show] do
    post :authenticate, on: :collection
  end

  resource :account, only: :show, controller: :account
end
