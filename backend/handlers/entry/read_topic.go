package entry

import (
	"bufio"
	"fmt"
	"strings"
)

func ReadTopic(markdown string, topic string) (string, error) {

	topicStart := fmt.Sprintf("# %s", strings.ToLower(topic))
	scanner := bufio.NewScanner(strings.NewReader(markdown))
	for scanner.Scan() {
		// TODO(mtlynch): Strip links from line
		if strings.ToLower(scanner.Text()) == topicStart {
			return strings.TrimSpace(readUntilNextHeading(scanner)), nil
		}
	}
	return "", nil
}

func readUntilNextHeading(scanner *bufio.Scanner) string {
	// TODO(mtlynch): Handle code blocks.
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			return strings.Join(lines, "\n")
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
