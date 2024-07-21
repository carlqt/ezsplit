package playground

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./playground/graphiql.html")
}
