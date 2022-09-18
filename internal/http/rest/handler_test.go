package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/phbpx/gobeer/internal/adding"
	"github.com/phbpx/gobeer/internal/http/rest"
	"github.com/phbpx/gobeer/internal/storage/postgres"
	"github.com/phbpx/gobeer/internal/storage/postgres/dbtest"
	"github.com/phbpx/gobeer/kit/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error

	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func TestHandler(t *testing.T) {
	t.Parallel()

	db, teardown := dbtest.New(t, c)
	defer teardown()

	repository := postgres.NewStorage(db)

	h := rest.NewHandler(rest.Config{
		Adding: adding.NewService(repository),
	})

	testPostBeer201(t, h)
	testPostBeer400(t, h)
	testPostBeer409(t, h)
}

func testPostBeer201(t *testing.T, h *rest.Handler) {
	nb := adding.NewBeer{
		Name:      "Test Beer",
		Brewery:   "Test Brewery",
		ShortDesc: "Test Short Description",
		Style:     "Test Style",
		ABV:       5.5,
	}

	body, err := json.Marshal(nb)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/beers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a new beer can be added.")
	{
		t.Log("\tWhen checking the response code.")
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t\t[ERROR] Should receive a 201 status code. Got %d", w.Code)
			}
			t.Log("\t\t[OK] Should receive a 201 status code.")
		}
	}
}

func testPostBeer400(t *testing.T, h *rest.Handler) {
	r := httptest.NewRequest("POST", "/beers", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a new beer can't be added with invalid payload.")
	{
		t.Log("\tWhen checking the response code.")
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t\t[ERROR] Should receive a 400 status code. Got %d", w.Code)
			}
			t.Log("\t\t[OK] Should receive a 400 status code.")
		}
	}
}

func testPostBeer409(t *testing.T, h *rest.Handler) {
	nb := adding.NewBeer{
		Name:      "Test Beer",
		Brewery:   "Test Brewery",
		ShortDesc: "Test Short Description",
		Style:     "Test Style",
		ABV:       5.5,
	}

	body, err := json.Marshal(nb)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/beers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a new beer can't be added with a duplicate name.")
	{
		t.Log("\tWhen checking the response code.")
		{
			if w.Code != http.StatusConflict {
				t.Fatalf("\t\t[ERROR] Should receive a 409 status code. Got %d", w.Code)
			}
			t.Log("\t\t[OK] Should receive a 409 status code.")
		}
	}
}
