package api

import (
	storage "module_31/pkg/storage"
	db "module_31/pkg/storage/db"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPI_postsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dbase, _ := db.New(ctx, "postgres://postgres:rootroot@localhost:5432/aggregator")
	dbase.AddPost(storage.Post{})

	api := New(dbase)

	req := httptest.NewRequest(http.MethodGet, "/news/10", nil)

	rr := httptest.NewRecorder()

	api.r.ServeHTTP(rr, req)

	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	var data []storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	const wantLen = 10
	if len(data) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(data), wantLen)
	}
}
