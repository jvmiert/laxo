package laxo

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"github.com/urfave/negroni"
)

// LoggerEntry is the structure passed to the template.
type LoggerEntry struct {
	Status   int
	Duration time.Duration
	Hostname string
	Method   string
	Path     string
	Request  *http.Request
}

// LoggerDefaultFormat is the format logged used by the default Logger instance.
var LoggerDefaultFormat = "{{.Method}} {{.Path}} | {{.Status}} | {{.Duration}} | {{.Hostname}}"

// LoggerDefaultDateFormat is the format used for date by the default Logger instance.
var LoggerDefaultDateFormat = time.RFC3339

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type NegroniZapLogger struct {
	Logger     *Logger
	dateFormat string
	template   *template.Template
}

// NewLogger returns a new Logger instance
func NewNegroniZapLogger(laxoLogger *Logger) *NegroniZapLogger {
	logger := &NegroniZapLogger{Logger: laxoLogger, dateFormat: LoggerDefaultDateFormat}
	logger.SetFormat(LoggerDefaultFormat)
	return logger
}

func (l *NegroniZapLogger) SetFormat(format string) {
	l.template = template.Must(template.New("negroni_parser").Parse(format))
}

func (l *NegroniZapLogger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *NegroniZapLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	log := LoggerEntry{
		Status:   res.Status(),
		Duration: time.Since(start),
		Hostname: r.Host,
		Method:   r.Method,
		Path:     r.URL.Path,
		Request:  r,
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)

	//@TODO: In production we will probably want to use Infow an pass the values
	//       instead of templating the string.
	l.Logger.Info(buff.String())
}
