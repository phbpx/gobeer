package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/phbpx/gobeer/internal/adding"
	"github.com/phbpx/gobeer/internal/beers"
	"github.com/phbpx/gobeer/internal/http/rest"
	"github.com/phbpx/gobeer/internal/reviewing"
	"github.com/phbpx/gobeer/internal/storage/postgres/dbtest"
	"github.com/phbpx/gobeer/pkg/docker"
	"go.opentelemetry.io/otel"
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

	test := dbtest.NewTest(t, c)
	defer test.Teardown()

	h := rest.NewHandler(rest.Config{
		Log:    test.Log,
		Tracer: otel.Tracer(""),
		DB:     test.DB,
	})

	testPostBeer201(t, h)
	testPostBeer400(t, h)
	testPostBeer409(t, h)
	testGetBeers200(t, h)
	testPostBeerReview201(t, h)
	testPostBeerReview400(t, h)
	testPostBeerReview404(t, h)
	testGetBeerReviews200(t, h)
	testGetBeerReviews204(t, h)
	testGetBeerReviews400(t, h)
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

func testGetBeers200(t *testing.T, h *rest.Handler) {
	r := httptest.NewRequest("GET", "/beers", nil)
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a list of beers can be retrieved.")
	{
		t.Log("\tWhen checking the response code.")
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t\t[ERROR] Should receive a 200 status code. Got %d", w.Code)
			}
			t.Log("\t\t[OK] Should receive a 200 status code.")
		}
	}
}

func testPostBeerReview201(t *testing.T, h *rest.Handler) {
	nr := reviewing.NewReview{
		UserID:  uuid.NewString(),
		Score:   5,
		Comment: "Test Comment",
	}

	body, err := json.Marshal(nr)
	if err != nil {
		t.Fatal(err)
	}

	beers := getBeers(t, h)
	if len(beers) == 0 {
		t.Fatal("No beers found")
	}

	beerID := beers[0].ID

	r := httptest.NewRequest("POST", fmt.Sprintf("/beers/%s/reviews", beerID), bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a new beer review can be added.")
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

func testPostBeerReview400(t *testing.T, h *rest.Handler) {
	r := httptest.NewRequest("POST", "/beers/123/reviews", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a new beer review can't be added with invalid payload.")
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

func testPostBeerReview404(t *testing.T, h *rest.Handler) {
	nr := reviewing.NewReview{
		UserID:  uuid.NewString(),
		Score:   3.0,
		Comment: "Test Comment",
	}

	body, err := json.Marshal(nr)
	if err != nil {
		t.Fatal(err)
	}

	beerID := uuid.NewString()

	r := httptest.NewRequest("POST", fmt.Sprintf("/beers/%s/reviews", beerID), bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a new beer review can't be added with a non existing beer.")
	{
		t.Log("\tWhen checking the response code.")
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t\t[ERROR] Should receive a 404 status code. Got %d", w.Code)
			}
			t.Log("\t\t[OK] Should receive a 404 status code.")
		}
	}
}

func testGetBeerReviews200(t *testing.T, h *rest.Handler) {
	beers := getBeers(t, h)
	if len(beers) == 0 {
		t.Fatal("No beers found")
	}

	beerID := beers[0].ID

	r := httptest.NewRequest("GET", fmt.Sprintf("/beers/%s/reviews", beerID), nil)
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a list of beer reviews can be retrieved.")
	{
		t.Log("\tWhen checking the response code.")
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t\t[ERROR] Should receive a 200 status code. Got %d", w.Code)
			}
			t.Log("\t\t[OK] Should receive a 200 status code.")
		}
	}
}

func testGetBeerReviews204(t *testing.T, h *rest.Handler) {
	beerID := uuid.NewString()

	r := httptest.NewRequest("GET", fmt.Sprintf("/beers/%s/reviews", beerID), nil)
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a list of beer reviews can't be retrieved with a non existing beer.")
	{
		t.Log("\tWhen checking the response code.")
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t\t[ERROR] Should receive a 204 status code. Got %d", w.Code)
			}
			t.Log("\t\t[OK] Should receive a 204 status code.")
		}
	}
}

func testGetBeerReviews400(t *testing.T, h *rest.Handler) {
	r := httptest.NewRequest("GET", "/beers/invalid/reviews", nil)
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	t.Log("Given the neeed to validate a list of beer reviews can't be retrieved with an invalid beer ID.")
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

func getBeers(t *testing.T, h *rest.Handler) []beers.Beer {
	r := httptest.NewRequest("GET", "/beers", nil)
	w := httptest.NewRecorder()

	h.Router().ServeHTTP(w, r)

	var beers []beers.Beer
	if err := json.NewDecoder(w.Body).Decode(&beers); err != nil {
		t.Fatal(err)
	}

	return beers
}
