package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

func ConfigLogger() {
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// Fallback: fixed UTC+05:30
		ist = time.FixedZone("IST", 5*60*60+30*60)
	}

	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, strings.ToUpper(t.In(ist).Format("02-Jan-06 15:04:05")))
				}
			}
			return a
		},
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
}
