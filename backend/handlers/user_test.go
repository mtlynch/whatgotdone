package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (ds mockDatastore) GetUserProfile(username string) (types.UserProfile, error) {
	return ds.userProfile, nil
}

func (ds *mockDatastore) SetUserProfile(username string, p types.UserProfile) error {
	ds.userProfile = p
	return nil
}

func TestUserPostStoresValidProfile(t *testing.T) {
	ds := mockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ "aboutMarkdown": "I'm a little teapot", "twitterHandle": "someTweeter", "emailAddress": "hi@example.com" }`)
	req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	userProfileExpected := types.UserProfile{
		AboutMarkdown: "I'm a little teapot",
		TwitterHandle: "someTweeter",
		EmailAddress:  "hi@example.com",
	}
	if !reflect.DeepEqual(ds.userProfile, userProfileExpected) {
		t.Fatalf("unexpected user profile: got %v want %v",
			ds.userProfile, userProfileExpected)
	}
}

func TestUserPostStoresEmptyProfileWhenRequestIsPartiallyPopulatedDict(t *testing.T) {
	ds := mockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ "twitterHandle": "someTweeter" }`)
	req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	userProfileExpected := types.UserProfile{
		TwitterHandle: "someTweeter",
	}
	if !reflect.DeepEqual(ds.userProfile, userProfileExpected) {
		t.Fatalf("unexpected user profile: got %v want %v",
			ds.userProfile, userProfileExpected)
	}
}

func TestUserPostStoresEmptyProfileWhenRequestIsEmptyDict(t *testing.T) {
	ds := mockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ }`)
	req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !reflect.DeepEqual(ds.userProfile, types.UserProfile{}) {
		t.Fatalf("unexpected user profile: got %v want %v",
			ds.userProfile, types.UserProfile{})
	}
}

func TestUserPostReturnsStatusBadRequestWhenRequestIsEmptyString(t *testing.T) {
	ds := mockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte("")
	req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
