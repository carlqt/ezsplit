Rails.application.routes.draw do
  resources :home, only: [:show, :index], controller: :home
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
  resources :accounts, only: [] do
    post :authenticate, on: :collection
  end
end
