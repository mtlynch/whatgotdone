package entry

import (
	"strings"
	"testing"
)

func TestReadProject(t *testing.T) {
	var tests = []struct {
		explanation         string
		markdown            string
		project             string
		projectBodyExpected string
		errExpected         error
	}{
		{
			"finds project when it starts an entry",
			`# Donuts

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`,
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"finds project when it is in the middle of an entry",
			`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Donuts

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`,
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"finds project when it is at the end of an entry",
			`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Soup

* Soup is reportedly not as delicious as donuts

# Donuts

* Donuts are delicious
* Multiple studies confirm this`,
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"finds project when it's a hyperlink",
			`# [Donuts](https://donutpalace.com)

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`,
			"donuts",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"canonicalizes multi-word project",
			`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Donut Updates

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`,
			"donut-updates",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"canonicalizes header dots into dashes",
			`# Kittens

* Adopted 17 kittens
* Named all of them mittens

# Donuts.com

* Donuts are delicious
* Multiple studies confirm this

# Soup

* Soup is reportedly not as delicious as donuts`,
			"donuts-com",
			`* Donuts are delicious
* Multiple studies confirm this`,
			nil,
		},
		{
			"ignores headers within code blocks",
			strings.ReplaceAll(`# Kittens

Wrote this code sample:

'''
# This is not a header because it is in a code block
print('Hello, world!')
'''

* Adopted 17 kittens
* Named all of them mittens

# Soup

* Soup is reportedly not as delicious as donuts`, "'''", "```"),
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
