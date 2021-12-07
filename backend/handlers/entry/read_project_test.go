package entry

import (
	"strings"
	"testing"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestReadProject(t *testing.T) {
	var tests = []struct {
		explanation         string
		markdown            types.EntryContent
		project             string
		projectBodyExpected string
		errExpected         error
	}{
		{
			"finds project when it starts an entry",
			types.EntryContent(`# Donuts

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`),
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"finds project when it is in the middle of an entry",
			types.EntryContent(`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Donuts

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`),
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"finds project when it is at the end of an entry",
			types.EntryContent(`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Soup

* Soup is reportedly not as delicious as donuts

# Donuts

* Donuts are delicious
* Multiple studies confirm this`),
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"finds project when it's a hyperlink",
			types.EntryContent(`# [Donuts](https://donutpalace.com)

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`),
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"canonicalizes multi-word project",
			types.EntryContent(`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Donut Updates

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`),
			"donut-updates",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"canonicalizes header dots into dashes",
			types.EntryContent(`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Donuts.com

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`),
			"donuts-com",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"ignores headers within code blocks",
			types.EntryContent(strings.ReplaceAll(`# Kittens

Wrote this code sample:

'''
# This is not a header because it is in a code block
print('Hello, world!')
'''

* Adopted 17 kittens
* Named all of them mittens

# Soup

* Soup is reportedly not as delicious as donuts`, "'''", "```")),
			"kittens",
			strings.ReplaceAll(`Wrote this code sample:

'''
# This is not a header because it is in a code block
print('Hello, world!')
'''

* Adopted 17 kittens
* Named all of them mittens`, "'''", "```"),
			nil,
		},
		{
			"returns ProjectNotFoundError when no project matches",
			types.EntryContent(`# Donuts

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`),
			"pineapples",
			"",
			ProjectNotFoundError{
				Project: "pineapples",
			},
		},
	}

	for _, tt := range tests {
		projectBodyActual, errActual := ReadProject(tt.markdown, tt.project)
		if errActual != tt.errExpected {
			t.Errorf("%s: input (%s, %s), got %v, want %v", tt.explanation, tt.markdown, tt.project, errActual, tt.errExpected)
		} else if tt.projectBodyExpected != projectBodyActual {
			t.Errorf("%s: input (%s, %s), got [%v], want [%v]", tt.explanation, tt.markdown, tt.project, projectBodyActual, tt.projectBodyExpected)
		}
	}
}
