package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
)

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
		// Trims leading and trailing whitespace from about field.
		{
			`{ "aboutMarkdown": " I'm a little teapot\t\n", "twitterHandle": "someTweeter", "emailAddress": "hi@example.com" }`,
			http.StatusOK,
			types.UserProfile{
				AboutMarkdown: "I'm a little teapot",
				TwitterHandle: "someTweeter",
				EmailAddress:  "hi@example.com",
			},
		},
		// Accept an update with only about field populated.
		{
			`{ "aboutMarkdown": "I'm a little teapot" }`,
			http.StatusOK,
			types.UserProfile{
				AboutMarkdown: "I'm a little teapot",
			},
		},
		// Accept an update with only twitter field populated.
		{
			`{ "twitterHandle": "someTweeter" }`,
			http.StatusOK,
			types.UserProfile{
				TwitterHandle: "someTweeter",
			},
		},
		// Accept an update with only email address field populated.
		{
			`{ "emailAddress": "hi@example.com" }`,
			http.StatusOK,
			types.UserProfile{
				EmailAddress: "hi@example.com",
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
		// If the request contains an illegal bio, reject it.
		{
			`{ "aboutMarkdown": "# Headings are invalid", "twitterHandle": "someTweeter", "emailAddress": "hi@example.com" }`,
			http.StatusBadRequest,
			types.UserProfile{},
		},
	}

	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	for _, tt := range userPostTests {
		req, err := http.NewRequest("POST", "/api/user", strings.NewReader(tt.requestBody))
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
		if tt.httpStatusExpected == http.StatusOK && !reflect.DeepEqual(ds.UserProfile, tt.userProfileExpected) {
			t.Fatalf("for input [%s], unexpected user profile: got %+v want %+v",
				tt.requestBody, ds.UserProfile, tt.userProfileExpected)
		}
	}
}
