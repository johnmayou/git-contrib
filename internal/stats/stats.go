package stats

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// Scans repo commits and prints a commits graph.
func Stats(repos []string, email string, now func() time.Time, since time.Time, stdout io.Writer) error {
	since = SnappedToSundayMorning(since)

	commitsByDay := make(map[string][]*object.Commit, (int(now().Sub(since).Hours()/24))+1)
	for _, repo := range repos {
		if err := scanRepoCommits(repo, email, since, commitsByDay); err != nil {
			return err
		}
	}

	return printCommitHeatmap(commitsByDay, now, since, stdout)
}

// Returns time set to 00:00 and snapped to the previous Sunday (if not already a Sunday).
func SnappedToSundayMorning(t time.Time) time.Time {
	return t.AddDate(0, 0, -int(t.Weekday())).Truncate(24 * time.Hour)
}

// Returns all commits authored by the given email, commited after the given time.
func scanRepoCommits(path string, email string, since time.Time, commitsByDay map[string][]*object.Commit) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	iter, err := repo.Log(&git.LogOptions{Since: &since})
	if err != nil {
		return err
	}
	defer iter.Close()
	return iter.ForEach(func(c *object.Commit) error {
		if c.Author.Email == email {
			day := timeToDateString(c.Author.When)
			commitsByDay[day] = append(commitsByDay[day], c)
		}
		return nil
	})
}

const dayBlock = "■"

type heatLevel int

const (
	heatLevel0 heatLevel = iota
	heatLevel1
	heatLevel2
	heatLevel3
	heatLevel4
)

var heatmapStyles = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("#161b22")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#0e4429")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#006d32")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#26a641")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#39d353")),
}

// Prints a daily commit heatmap to stdout.
func printCommitHeatmap(commitsByDay map[string][]*object.Commit, now func() time.Time, since time.Time, stdout io.Writer) error {
	if since.Weekday() != time.Sunday {
		return fmt.Errorf("since should be a Sunday, but got %s", since.Weekday().String())
	}

	weeks := int(now().Sub(since).Hours() / (24 * 7))
	for row := range 7 {
		for col := range weeks {
			count, day := 0, timeToDateString(since.AddDate(0, 0, (col*7)+row))
			if commits, ok := commitsByDay[day]; ok {
				count = len(commits)
			}
			level := commitCountHeatLevel(count)
			block := heatmapStyles[level].Render(dayBlock)
			if _, err := fmt.Fprint(stdout, block+" "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(stdout); err != nil {
			return err
		}
	}

	return nil
}

// Returns heatmap heat level based on given day commit count.
func commitCountHeatLevel(count int) heatLevel {
	switch {
	case count >= 20:
		return heatLevel4
	case count >= 10:
		return heatLevel3
	case count >= 5:
		return heatLevel2
	case count >= 1:
		return heatLevel1
	default:
		return heatLevel0
	}
}

// Returns time formatted as YYYY-MM-DD.
func timeToDateString(t time.Time) string {
	return t.Format("2006-01-02")
}
