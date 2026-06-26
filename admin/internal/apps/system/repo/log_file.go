package repo

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type LogFile struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LogQuery struct {
	FileName string
	Include  string
	Exclude  string
	StartAt  time.Time
	EndAt    time.Time
}

type LogLine struct {
	FileName string    `json:"file_name"`
	Line     int       `json:"line"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time,omitempty"`
}

func (d *Data) FindLogList() ([]LogFile, error) {
	entries, err := os.ReadDir(d.logDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []LogFile{}, nil
		}
		return nil, err
	}

	logs := make([]LogFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		logs = append(logs, LogFile{
			Name:      entry.Name(),
			Path:      filepath.Join(d.logDir, entry.Name()),
			Size:      info.Size(),
			UpdatedAt: info.ModTime(),
		})
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].UpdatedAt.After(logs[j].UpdatedAt)
	})
	return logs, nil
}

func (d *Data) SearchLogs(query LogQuery) ([]LogLine, error) {
	logs, err := d.FindLogList()
	if err != nil {
		return nil, err
	}

	results := make([]LogLine, 0)
	for _, logFile := range logs {
		if query.FileName != "" && logFile.Name != query.FileName {
			continue
		}

		lines, err := d.searchLogFile(logFile, query)
		if err != nil {
			return nil, err
		}
		results = append(results, lines...)
	}
	return results, nil
}

func (d *Data) searchLogFile(logFile LogFile, query LogQuery) ([]LogLine, error) {
	file, err := os.Open(logFile.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results := make([]LogLine, 0)
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		content := scanner.Text()
		if query.Include != "" && !strings.Contains(content, query.Include) {
			continue
		}
		if query.Exclude != "" && strings.Contains(content, query.Exclude) {
			continue
		}

		logTime := parseLogTime(content)
		if !query.StartAt.IsZero() && (logTime.IsZero() || logTime.Before(query.StartAt)) {
			continue
		}
		if !query.EndAt.IsZero() && (logTime.IsZero() || logTime.After(query.EndAt)) {
			continue
		}

		results = append(results, LogLine{
			FileName: logFile.Name,
			Line:     lineNo,
			Content:  content,
			Time:     logTime,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func parseLogTime(content string) time.Time {
	fields := strings.Fields(content)
	for _, field := range fields {
		field = strings.Trim(field, "[]\"")
		for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
			t, err := time.ParseInLocation(layout, field, time.Local)
			if err == nil {
				return t
			}
		}
	}
	return time.Time{}
}
