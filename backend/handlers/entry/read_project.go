package entry

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// ProjectNotFoundError occurs when an entry does not contain the given project.
type ProjectNotFoundError struct {
	Project string
}

func (f ProjectNotFoundError) Error() string {
	return fmt.Sprintf("Entry does not contain project %s", f.Project)
}

// ReadProject reads the body of a project, starting from a project header and
// ending at the following project header or the end of the entry.
func ReadProject(markdown types.EntryContent, project string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(markdown)))
	for scanner.Scan() {
		// TODO(mtlynch): Strip formatting from line.
		if readHeading(scanner.Text()) == project {
			return strings.TrimSpace(readUntilNextHeading(scanner)), nil
		}
	}
	return "", ProjectNotFoundError{
		Project: project,
	}
}

// TODO(mtlynch): Support other level headers.
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

func canonicalizeHeading(project string) string {
	re := regexp.MustCompile(`[^A-Za-z]+`)
	return re.ReplaceAllString(project, "-")
}

func readUntilNextHeading(scanner *bufio.Scanner) string {
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if lineHasCodeBlockDelimiter(line) {
			lines = append(lines, line)
			lines = append(lines, readUntilCodeBlockEnd(scanner))
			continue
		}
		if readHeading(line) != "" {
			return strings.Join(lines, "\n")
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func lineHasCodeBlockDelimiter(line string) bool {
	const codeBlockDelimiter = "```"
	return strings.HasPrefix(line, codeBlockDelimiter)
}

func readUntilCodeBlockEnd(scanner *bufio.Scanner) string {
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if lineHasCodeBlockDelimiter(line) {
			break
		}
	}
	return strings.Join(lines, "\n")
}
