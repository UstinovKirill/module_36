package api

import (
	storage "module_31/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	db storage.Interface
}

// Функция New выполняет роль конструктора
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Маршрутизатор запросов
func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {

	api.r.HandleFunc("/news/{n}", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Функция postsHandler возвращает все публикации
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)

	news, err := api.db.Posts(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)

}
