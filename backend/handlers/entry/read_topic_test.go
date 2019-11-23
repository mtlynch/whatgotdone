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
