package entry

import (
	"bufio"
	"regexp"
	"strings"
)

func ReadTopic(markdown string, topic string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(markdown))
	for scanner.Scan() {

		// TODO(mtlynch): Strip links from line
		// TODO(mtlynch): Strip formatting from line.
		if readHeading(scanner.Text()) == topic {
			return strings.TrimSpace(readUntilNextHeading(scanner)), nil
		}
	}
	return "", nil
}

const headerPrefix = "# "

func readHeading(line string) string {
	if !strings.HasPrefix(line, headerPrefix) {
		return ""
	}
	heading := line[len(headerPrefix):]
	heading = strings.ToLower(heading)
	heading = stripMarkdownLink(heading)
	heading = canonicalizeHeading(heading)
	return heading
}

func stripMarkdownLink(line string) string {
	re := regexp.MustCompile(`\[(.+)\]\(.+\)`)
	match := re.FindStringSubmatch(line)
	if len(match) == 0 {
		return line

	}
	return match[1]
}

func canonicalizeHeading(topic string) string {
	re := regexp.MustCompile(`[^A-Za-z]+`)
	return re.ReplaceAllString(topic, "-")
}

func readUntilNextHeading(scanner *bufio.Scanner) string {
	// TODO(mtlynch): Handle code blocks.
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if readHeading(line) != "" {
			return strings.Join(lines, "\n")
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
