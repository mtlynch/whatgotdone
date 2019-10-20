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

func TestUserPost(t *testing.T) {
	var userPostTests = []struct {
		requestBody         string
		httpStatusExpected  int
		userProfileExpected types.UserProfile
	}{
		// Valid case.
		{
			`{ "aboutMarkdown": "I'm a little teapot", "twitterHandle": "someTweeter", "emailAddress": "hi@example.com" }`,
			http.StatusOK,
			types.UserProfile{
				AboutMarkdown: "I'm a little teapot",
				TwitterHandle: "someTweeter",
				EmailAddress:  "hi@example.com",
			},
		},
		// Partially populated request body.
		{
			`{ "twitterHandle": "someTweeter" }`,
			http.StatusOK,
			types.UserProfile{
				TwitterHandle: "someTweeter",
			},
		},
		// When request body is empty dict, store empty profile.
		{
			`{ }`,
			http.StatusOK,
			types.UserProfile{},
		},
		// When request body is empty string, return bad request error.
		{
			"",
			http.StatusBadRequest,
			types.UserProfile{},
		},
		// When request body is malformed dict, return bad request error.
		{
			"{",
			http.StatusBadRequest,
			types.UserProfile{},
		},
		// If the request contains a malformed email, reject it.
		{
			`{ "aboutMarkdown": "I'm a little teapot", "twitterHandle": "someTweeter", "emailAddress": "hi[at]example.com" }`,
			http.StatusBadRequest,
			types.UserProfile{},
		},
		// If the request contains an invalid Twitter handle, reject it.
		{
			`{ "aboutMarkdown": "I'm a little teapot", "twitterHandle": "inv@liD==h&le", "emailAddress": "hi@example.com" }`,
			http.StatusBadRequest,
			types.UserProfile{},
		},
	}

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

	for _, tt := range userPostTests {
		req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer([]byte(tt.requestBody)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		if status := w.Code; status != tt.httpStatusExpected {
			t.Fatalf("for input [%s], handler returned wrong status code: got %v want %v",
				tt.requestBody, status, tt.httpStatusExpected)
		}
		if tt.httpStatusExpected == http.StatusOK && !reflect.DeepEqual(ds.userProfile, tt.userProfileExpected) {
			t.Fatalf("for input [%s], unexpected user profile: got %v want %v",
				tt.requestBody, ds.userProfile, tt.userProfileExpected)
		}
	}
}
