package logging

import "os"

// LogFileAppender simple io.writer adapter, which will append to the configured log file
type LogFileAppender struct {
	logFileName string
}

func NewLogFileAppender(logFileName string) *LogFileAppender {
	return &LogFileAppender{logFileName: logFileName}
}

func (w *LogFileAppender) Write(p []byte) (int, error) {
	// should work if file already exists and opened by another process
	f, err := os.OpenFile(w.logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return 0, err
	}

	defer f.Close()

	if _, err := f.Write(p); err != nil {
		return 0, err
	}

	return len(p), nil
}
