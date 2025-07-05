package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
	"github.com/mtlynch/whatgotdone/backend/types/export"
)

func TestExportPopulatedUserAccount(t *testing.T) {
	ds := mock.MockDatastore{
		JournalDrafts: []types.JournalEntry{
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-12"),
				LastModified: mustParseTime("2021-11-12T00:00:00Z"),
				Markdown:     "thought about fishing",
			},
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-19"),
				LastModified: mustParseTime("2021-11-19T00:00:00Z"),
				Markdown:     "went to the store today",
			},
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-26"),
				LastModified: mustParseTime("2021-11-20T00:00:00Z"),
				Markdown:     "bought a new car",
			},
		},
		JournalEntries: []types.JournalEntry{
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-12"),
				LastModified: mustParseTime("2021-11-12T00:00:00Z"),
				Markdown:     "thought about fishing",
			},
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-19"),
				LastModified: mustParseTime("2021-11-19T00:00:00Z"),
				Markdown:     "went to the store today",
			},
		},
		Reactions: map[types.Username]map[types.EntryDate][]types.Reaction{
			"dummyUserA": {
				"2021-11-19": []types.Reaction{
					{
						Username:  types.Username("dummyUserB"),
						Symbol:    "👍",
						Timestamp: mustParseTime("2021-11-20T11:57:02Z"),
					},
				},
			},
		},
		UserFollows: map[types.Username][]types.Username{
			types.Username("dummyUserA"): {types.Username("dummyUserC")},
		},
		UserPreferences: map[types.Username]types.Preferences{
			types.Username("dummyUserA"): {
				EntryTemplate: "# My weekly template",
			},
		},
		UserProfile: types.UserProfile{
			AboutMarkdown: "I'm just a dummy user",
			EmailAddress:  "dummy@example.com",
		},
		Usernames: []types.Username{"dummyUserA", "dummyUserB", "dummyUserC"},
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUserA",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/export", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response export.UserData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	exportExpected := export.UserData{
		Drafts: []export.JournalEntry{
			{
				Date:         types.EntryDate("2021-11-12"),
				LastModified: "2021-11-12T00:00:00Z",
				Markdown:     "thought about fishing",
			},
			{
				Date:         types.EntryDate("2021-11-19"),
				LastModified: "2021-11-19T00:00:00Z",
				Markdown:     "went to the store today",
			},
			{
				Date:         types.EntryDate("2021-11-26"),
				LastModified: "2021-11-20T00:00:00Z",
				Markdown:     "bought a new car",
			},
		},
		Entries: []export.JournalEntry{
			{
				Date:         types.EntryDate("2021-11-12"),
				LastModified: "2021-11-12T00:00:00Z",
				Markdown:     "thought about fishing",
			},
			{
				Date:         types.EntryDate("2021-11-19"),
				LastModified: "2021-11-19T00:00:00Z",
				Markdown:     "went to the store today",
			},
		},
		Reactions: map[types.EntryDate][]export.Reaction{
			"2021-11-19": {
				{
					Username:  types.Username("dummyUserB"),
					Symbol:    "👍",
					Timestamp: "2021-11-20T11:57:02Z",
				},
			},
		},
		Following: []types.Username{types.Username("dummyUserC")},
		Preferences: export.Preferences{
			EntryTemplate: "# My weekly template",
		},
		Profile: export.UserProfile{
			AboutMarkdown: "I'm just a dummy user",
			EmailAddress:  "dummy@example.com",
		},
	}
	if !reflect.DeepEqual(response, exportExpected) {
		t.Fatalf("Unexpected response: got %+v want %+v", response, exportExpected)
	}
}

func TestExportEmptyUserAccount(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUserA",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/export", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response export.UserData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	exportExpected := export.UserData{
		Entries:   []export.JournalEntry{},
		Reactions: map[types.EntryDate][]export.Reaction{},
		Drafts:    []export.JournalEntry{},
	}
	if !reflect.DeepEqual(response, exportExpected) {
		t.Fatalf("Unexpected response: got %#v want %#v", response, exportExpected)
	}
}

func TestExportUnauthenticatedAccount(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUserA",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/export", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock-invalid-token", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusUnauthorized {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestEntryToMarkdown(t *testing.T) {
	for _, tt := range []struct {
		explanation string
		input       types.JournalEntry
		expected    string
	}{
		{
			"simple export when lastmod is the same as the journal date",
			types.JournalEntry{
				Date:         types.EntryDate("2025-07-04"),
				LastModified: mustParseTime("2025-07-04T11:00:00Z"),
				Markdown:     types.EntryContent("Good week!"),
			},
			`
---
date: 2025-07-04
---
Good week!
			`,
		},
		{
			"simple export when lastmod is different from the journal date",
			types.JournalEntry{
				Date:         types.EntryDate("2025-07-04"),
				LastModified: mustParseTime("2025-07-05T09:30:00Z"),
				Markdown:     types.EntryContent("Wrote this one a little late..."),
			},
			`
---
date: 2025-07-04
lastmod: 2025-07-05
---
Wrote this one a little late...
			`,
		},
	} {
		t.Run(tt.explanation, func(t *testing.T) {
			actual := entryToMarkdown(tt.input)
			if got, want := actual, strings.TrimSpace(tt.expected); got != want {
				t.Errorf("markdown=%v, want=%v", got, want)
			}
		})
	}
}

func TestPackageEntriesAsMarkdown(t *testing.T) {
	for _, tt := range []struct {
		description string
		entries     []types.JournalEntry
		wantFiles   map[string]string // filepath -> expected content
	}{
		{
			"single entry creates one markdown file",
			[]types.JournalEntry{
				{
					Date:         types.EntryDate("2025-07-04"),
					LastModified: mustParseTime("2025-07-04T11:00:00Z"),
					Markdown:     types.EntryContent("Good week!"),
				},
			},
			map[string]string{
				"2025-07-04/index.md": strings.TrimSpace(`
---
date: 2025-07-04
---
Good week!
				`),
			},
		},
		{
			"multiple entries create multiple markdown files",
			[]types.JournalEntry{
				{
					Date:         types.EntryDate("2025-07-04"),
					LastModified: mustParseTime("2025-07-04T11:00:00Z"),
					Markdown:     types.EntryContent("Good week!"),
				},
				{
					Date:         types.EntryDate("2025-06-27"),
					LastModified: mustParseTime("2025-06-28T09:30:00Z"),
					Markdown:     types.EntryContent("Busy week with lots of coding."),
				},
				{
					Date:         types.EntryDate("2025-06-20"),
					LastModified: mustParseTime("2025-06-20T15:45:00Z"),
					Markdown:     types.EntryContent("Started a new project!"),
				},
			},
			map[string]string{
				"2025-07-04/index.md": strings.TrimSpace(`
---
date: 2025-07-04
---
Good week!
				`),
				"2025-06-27/index.md": strings.TrimSpace(`
---
date: 2025-06-27
lastmod: 2025-06-28
---
Busy week with lots of coding.
				`),
				"2025-06-20/index.md": strings.TrimSpace(`
---
date: 2025-06-20
---
Started a new project!
				`),
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			reader, err := packageEntriesAsMarkdown(tt.entries)
			if err != nil {
				t.Fatalf("packageEntriesAsMarkdown() failed: %v", err)
			}

			// Extract zip contents
			zipContents := extractZipContents(t, reader)

			// Verify we have the expected number of files
			if len(zipContents) != len(tt.wantFiles) {
				t.Errorf("expected %d files in zip, got %d", len(tt.wantFiles), len(zipContents))
			}

			// Verify each expected file exists with correct content
			for expectedPath, expectedContent := range tt.wantFiles {
				actualContent, exists := zipContents[expectedPath]
				if !exists {
					t.Errorf("expected file %s not found in zip", expectedPath)
					continue
				}

				if actualContent != expectedContent {
					t.Errorf("file %s content mismatch:\ngot:\n%s\nwant:\n%s",
						expectedPath, actualContent, expectedContent)
				}
			}
		})
	}
}

func TestFindMediaURLs(t *testing.T) {
	for _, tt := range []struct {
		description string
		markdown    string
		expected    []string
	}{
		{
			"no media URLs",
			"Just some regular text with no images.",
			[]string{},
		},
		{
			"single media.whatgotdone.com URL",
			"Check out this image: https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png",
			[]string{"https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png"},
		},
		{
			"single storage.googleapis.com URL",
			"Check out this image: https://storage.googleapis.com/media.whatgotdone.com/uploads/michael/20200501/HK5a.png",
			[]string{"https://storage.googleapis.com/media.whatgotdone.com/uploads/michael/20200501/HK5a.png"},
		},
		{
			"multiple URLs of different types",
			"Images: https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png and https://storage.googleapis.com/media.whatgotdone.com/uploads/michael/20200501/HK5b.png",
			[]string{
				"https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png",
				"https://storage.googleapis.com/media.whatgotdone.com/uploads/michael/20200501/HK5b.png",
			},
		},
		{
			"URL in markdown link syntax",
			"* [Before](https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png)",
			[]string{"https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png"},
		},
		{
			"URL in markdown image syntax",
			"![Screenshot](https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png)",
			[]string{"https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png"},
		},
		{
			"ignores other domains",
			"External image: https://example.com/image.png and our image: https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png",
			[]string{"https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png"},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			actual := findMediaURLs(tt.markdown)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("findMediaURLs() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestExtractFilename(t *testing.T) {
	for _, tt := range []struct {
		url      string
		expected string
	}{
		{
			"https://media.whatgotdone.com/uploads/michael/20200501/HK5a.png",
			"HK5a.png",
		},
		{
			"https://storage.googleapis.com/media.whatgotdone.com/uploads/michael/20200501/HK5b.jpg",
			"HK5b.jpg",
		},
		{
			"https://media.whatgotdone.com/simple.gif",
			"simple.gif",
		},
	} {
		t.Run(tt.url, func(t *testing.T) {
			actual := extractFilename(tt.url)
			if actual != tt.expected {
				t.Errorf("extractFilename(%s) = %s, want %s", tt.url, actual, tt.expected)
			}
		})
	}
}

func TestProcessMarkdownWithImages(t *testing.T) {
	// Note: This test doesn't actually download images since that would require
	// setting up a mock HTTP server. In a real implementation, you might want
	// to add integration tests that use a mock server.

	for _, tt := range []struct {
		description string
		markdown    string
		expected    string
		expectError bool
	}{
		{
			"no media URLs - no changes",
			"Just some regular text with no images.",
			"Just some regular text with no images.",
			false,
		},
		{
			"markdown with non-media URLs - no changes",
			"External image: https://example.com/image.png",
			"External image: https://example.com/image.png",
			false,
		},
		{
			"markdown with WhatGotDone media URL - should fail to download",
			"Check out this image: https://media.whatgotdone.com/uploads/michael/20200501/nonexistent.png",
			"",
			true,
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			zipWriter := zip.NewWriter(&buf)

			actual, err := processMarkdownWithImages(tt.markdown, zipWriter, "2020-04-10/")

			if tt.expectError {
				if err == nil {
					t.Errorf("processMarkdownWithImages() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("processMarkdownWithImages() failed: %v", err)
			}

			zipWriter.Close()

			if actual != tt.expected {
				t.Errorf("processMarkdownWithImages() = %s, want %s", actual, tt.expected)
			}
		})
	}
}

// extractZipContents reads a zip file from an io.Reader and returns a map of file paths to their contents
func extractZipContents(t *testing.T, r io.Reader) map[string]string {
	// Read all data from the reader
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read zip data: %v", err)
	}

	// Create a zip reader from the data
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("failed to create zip reader: %v", err)
	}

	// Extract all files
	contents := make(map[string]string)
	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			t.Fatalf("failed to open file %s in zip: %v", file.Name, err)
		}

		fileData, err := io.ReadAll(fileReader)
		if err != nil {
			t.Fatalf("failed to read file %s content: %v", file.Name, err)
		}
		fileReader.Close()

		contents[file.Name] = string(fileData)
	}

	return contents
}
