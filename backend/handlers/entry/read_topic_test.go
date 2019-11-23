package entry

import (
	"testing"
)

func TestReadTopic(t *testing.T) {
	var tests = []struct {
		explanation       string
		markdown          string
		topic             string
		topicBodyExpected string
		errExpected       error
	}{
		{
			"finds topic when it starts an entry",
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
			"finds topic when it is in the middle of an entry",
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
			"finds topic when it is at the end of an entry",
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
			"finds topic when it's a hyperlink",
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
			"canonicalizes multi-word topic",
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
			"canonicalizes dots into dashes",
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
	}

	for _, tt := range tests {
		topicBodyActual, errActual := ReadTopic(tt.markdown, tt.topic)
		if errActual != tt.errExpected {
			t.Errorf("%s: input (%s, %s), got %v, want %v", tt.explanation, tt.markdown, tt.topic, errActual, tt.errExpected)
		} else if tt.topicBodyExpected != topicBodyActual {
			t.Errorf("%s: input (%s, %s), got [%v], want [%v]", tt.explanation, tt.markdown, tt.topic, topicBodyActual, tt.topicBodyExpected)
		}
	}
}
