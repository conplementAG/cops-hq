package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"time"

	"github.com/mattn/go-colorable"
)

var DefaultLogLevel = logrus.InfoLevel

// Init will initialize the Logrus system, which per default sets up logging to
// console and to file at the same time. Features are file rotation,
// fixed colors on Windows etc.
// Init is a global initialization (without return value), since Logrus and the
// logging systems are also globally available singletons
func Init(logFileName string) {
	var logLevel = DefaultLogLevel

	// main reason we use the prefixed library TextFormatter, instead of the default logrus.TextFormatter,
	// it to have the "ForceFormatting" option, which enables the same format in all TTY and non-TTY
	// execution environments.
	consoleFormatter := &prefixed.TextFormatter{
		// required to show colors in build for example, which would otherwise
		// show plain output because build is not a TTY session
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
		// make sure the format is kept in all execution envs
		ForceFormatting: true,
	}

	fileFormatter := &prefixed.TextFormatter{
		// files will show colors in raw format, which pollutes the files too much.
		// Problematic output look like: "[36mINFO[0m", which we don't want.
		// Colors are only available as console output.
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
		// make sure the format is kept in all execution envs
		ForceFormatting: true,
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logFileName,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     90, // days
		Level:      logLevel,
		Formatter:  fileFormatter,
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logrus.SetLevel(logLevel)

	// special stdout io.Writer capable of colors on Windows
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(consoleFormatter)

	// this hook will also route logs to file
	logrus.AddHook(rotateFileHook)
}
