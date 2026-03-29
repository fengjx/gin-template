package log

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

type dailyFileWriter struct {
	baseFilename string
	maxBackups   int
	maxAgeDays   int
	compress     bool
	nowFunc      func() time.Time

	mu          sync.Mutex
	currentDate string
	currentPath string
	file        *os.File
}

func newDailyFileWriter(baseFilename string, maxBackups int, maxAgeDays int, compress bool) *dailyFileWriter {
	return &dailyFileWriter{
		baseFilename: baseFilename,
		maxBackups:   maxBackups,
		maxAgeDays:   maxAgeDays,
		compress:     compress,
		nowFunc:      time.Now,
	}
}

func (w *dailyFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeededLocked(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *dailyFileWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}
	return w.file.Sync()
}

func (w *dailyFileWriter) rotateIfNeededLocked() error {
	now := w.nowFunc().In(time.Local)
	dateKey := now.Format(time.DateOnly)
	targetPath := datedFilename(w.baseFilename, now)
	if w.file != nil && w.currentDate == dateKey && w.currentPath == targetPath {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	if w.file != nil {
		_ = w.file.Sync()
		_ = w.file.Close()
	}

	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	w.file = file
	w.currentDate = dateKey
	w.currentPath = targetPath
	return w.cleanupLocked(now)
}

func (w *dailyFileWriter) cleanupLocked(now time.Time) error {
	entries, err := os.ReadDir(filepath.Dir(w.baseFilename))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	type logFile struct {
		path       string
		date       time.Time
		compressed bool
	}

	files := make([]logFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		date, compressed, ok := parseDatedFilename(w.baseFilename, entry.Name())
		if !ok {
			continue
		}
		files = append(files, logFile{
			path:       filepath.Join(filepath.Dir(w.baseFilename), entry.Name()),
			date:       date,
			compressed: compressed,
		})
	}

	slices.SortFunc(files, func(a, b logFile) int {
		switch {
		case a.date.Before(b.date):
			return -1
		case a.date.After(b.date):
			return 1
		default:
			return strings.Compare(a.path, b.path)
		}
	})

	if w.compress {
		for _, item := range files {
			if item.compressed || item.path == w.currentPath {
				continue
			}
			if !item.date.Before(now.Truncate(24 * time.Hour)) {
				continue
			}
			if err := compressFile(item.path); err != nil {
				return err
			}
		}
		files = files[:0]
		entries, err = os.ReadDir(filepath.Dir(w.baseFilename))
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			date, compressed, ok := parseDatedFilename(w.baseFilename, entry.Name())
			if !ok {
				continue
			}
			files = append(files, logFile{
				path:       filepath.Join(filepath.Dir(w.baseFilename), entry.Name()),
				date:       date,
				compressed: compressed,
			})
		}
		slices.SortFunc(files, func(a, b logFile) int {
			switch {
			case a.date.Before(b.date):
				return -1
			case a.date.After(b.date):
				return 1
			default:
				return strings.Compare(a.path, b.path)
			}
		})
	}

	if w.maxAgeDays > 0 {
		expireBefore := now.AddDate(0, 0, -w.maxAgeDays)
		for _, item := range files {
			if item.path == w.currentPath {
				continue
			}
			if item.date.Before(truncateDate(expireBefore)) {
				_ = os.Remove(item.path)
			}
		}
	}

	if w.maxBackups > 0 {
		latestByDate := make(map[string]logFile)
		dateKeys := make([]string, 0, len(files))
		for _, item := range files {
			key := item.date.Format(time.DateOnly)
			existing, ok := latestByDate[key]
			if !ok {
				dateKeys = append(dateKeys, key)
				latestByDate[key] = item
				continue
			}
			if existing.compressed && !item.compressed {
				latestByDate[key] = item
			}
		}
		slices.Sort(dateKeys)
		if len(dateKeys) > w.maxBackups {
			for _, key := range dateKeys[:len(dateKeys)-w.maxBackups] {
				target := latestByDate[key]
				_ = os.Remove(target.path)
				raw := strings.TrimSuffix(target.path, ".gz")
				if raw != target.path {
					_ = os.Remove(raw)
				}
				_ = os.Remove(raw + ".gz")
			}
		}
	}

	return nil
}

func datedFilename(baseFilename string, now time.Time) string {
	dir := filepath.Dir(baseFilename)
	name := filepath.Base(baseFilename)
	ext := filepath.Ext(name)
	prefix := strings.TrimSuffix(name, ext)
	return filepath.Join(dir, prefix+"-"+now.In(time.Local).Format(time.DateOnly)+ext)
}

func parseDatedFilename(baseFilename string, entryName string) (time.Time, bool, bool) {
	name := filepath.Base(baseFilename)
	ext := filepath.Ext(name)
	prefix := strings.TrimSuffix(name, ext) + "-"
	suffix := ext
	compressed := false
	if strings.HasSuffix(entryName, ".gz") {
		entryName = strings.TrimSuffix(entryName, ".gz")
		compressed = true
	}
	if !strings.HasPrefix(entryName, prefix) || !strings.HasSuffix(entryName, suffix) {
		return time.Time{}, false, false
	}
	datePart := strings.TrimSuffix(strings.TrimPrefix(entryName, prefix), suffix)
	parsed, err := time.ParseInLocation(time.DateOnly, datePart, time.Local)
	if err != nil {
		return time.Time{}, false, false
	}
	return parsed, compressed, true
}

func compressFile(path string) error {
	source, err := os.Open(path)
	if err != nil {
		return err
	}
	defer source.Close()

	targetPath := path + ".gz"
	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	gzipWriter := gzip.NewWriter(target)
	if _, err := io.Copy(gzipWriter, source); err != nil {
		_ = gzipWriter.Close()
		_ = target.Close()
		return err
	}
	if err := gzipWriter.Close(); err != nil {
		_ = target.Close()
		return err
	}
	if err := target.Close(); err != nil {
		return err
	}
	return os.Remove(path)
}

func truncateDate(t time.Time) time.Time {
	local := t.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local)
}
