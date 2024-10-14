package parse

import (
	"strings"
	"testing"
)

func TestUserBio(t *testing.T) {
	for _, tt := range []struct {
		explanation   string
		bio           string
		validExpected bool
	}{
		{
			"simple text bio is valid",
			"Hi, I'm dummyuser. Thanks for visiting.",
			true,
		},
		{
			"bio with basic formatting is valid",
			"Hi, I'm **dummyuser**. Thanks for visiting.",
			true,
		},
		{
			"bio with a link is valid",
			"Hi, I'm dummyuser. Check out my [blog](https://blog.example.com).",
			true,
		},
		{
			"bio with newlines is valid",
			"Hi, I'm dummyuser.\n\nThanks for visiting.",
			true,
		},
		{
			"bio with hashtags is valid",
			"Hi, I'm dummyuser. I love #programming and #coffee.",
			true,
		},
		{
			"empty bio is valid",
			"",
			true,
		},
		{
			"bio that's exactly 300 characters is valid",
			strings.Repeat("A", UserBioMaxLength),
			true,
		},
		{
			"bio longer than 300 characters is invalid",
			strings.Repeat("A", UserBioMaxLength+1),
			false,
		},
		{
			"bio with a leading heading is invalid",
			"# My Life Story\n\n It all started 10 years ago...",
			false,
		},
		{
			"bio with a heading in the middle is invalid",
			"Welcome: \n\n# My Life Story\n\n It all started 10 years ago...",
			false,
		},
		{
			"bio with fenced code is invalid",
			"My life is like this app:\n\n```\nprint 'hello, world!'\n```\n",
			false,
		},
		{
			"bio with an inline-style image is invalid",
			"Here's me: ![image alt text](http://example.com/photo \"My avatar\")",
			false,
		},
		{
			"bio with an inline-style image is invalid",
			"Here's me: ![image alt text][me]\n\n[me]: http://example.com/photo \"My avatar\"",
			false,
		},
	} {
		t.Run(tt.explanation, func(t *testing.T) {
			_, err := UserBio(tt.bio)
			if got, want := err == nil, tt.validExpected; got != want {
				t.Fatalf("valid=%v, want %v", got, want)
			}
		})
	}
}
