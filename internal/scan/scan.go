package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Recursively scans the given directory and and saves the all git
// repo paths to the output path.
func Scan(dir string, ignore []string, cachePath string, cacheOverride bool) error {
	repos, err := scanRepos(dir, ignore, make([]string, 0))
	if err != nil {
		return err
	}
	if cacheOverride {
		clearCache(cachePath)
	}
	err = cacheRepos(repos, cachePath)
	if err != nil {
		return err
	}
	return nil
}

// Returns git repo paths from the given file.
// If the file does not exist, it returns an empty slice.
func Load(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []string{}, nil
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
}

// Returns a list of git repo paths.
func scanRepos(dir string, ignore []string, repos []string) ([]string, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := file.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || slices.Contains(ignore, entry.Name()) {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if entry.Name() == ".git" {
			repos = append(repos, strings.TrimSuffix(path, "/.git"))
			continue
		}

		repos, err = scanRepos(path, ignore, repos)
		if err != nil {
			return nil, err
		}
	}
	return repos, nil
}

// Caches git repository paths to given file.
func cacheRepos(repos []string, cacheFile string) error {
	existingData, err := os.ReadFile(cacheFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var existingRepos []string
	if len(existingData) > 0 {
		existingRepos = strings.Split(strings.TrimSpace(string(existingData)), "\n")
	}

	data := []byte(strings.Join(joinSlicesUnique(repos, existingRepos), "\n"))
	return os.WriteFile(cacheFile, data, 0644)
}

// Truncates git repository cache file.
func clearCache(cacheFile string) error {
	if err := os.Truncate(cacheFile, 0); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

// Joins two slices and removes duplicates.
func joinSlicesUnique(a, b []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(a)+len(b))

	for _, s := range a {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	for _, s := range b {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}

	return result
}
