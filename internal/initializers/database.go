package initializers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB creates a pool to connect into a postgres database
// and returns it. Please always remember to defer close it
// after using.
func ConnectDB() *pgxpool.Pool {
	configString := fmt.Sprintf(
		"user=%v host=%v port=%v dbname=%v pool_max_conns=%v",
		os.Getenv("GO_POSTGRES_USER"),
		os.Getenv("GO_POSTGRES_HOST"),
		os.Getenv("GO_POSTGRES_PORT"),
		os.Getenv("GO_POSTGRES_DBNAME"),
		os.Getenv("GO_POSTGRES_POOL"),
	)

	config, err := pgxpool.ParseConfig(configString)

	if err != nil {
		log.Fatal("Error on the config string, please check env vars")
	}

	config.ConnConfig.Tracer = NewMultiQueryTracer(NewLoggingQueryTracer(slog.Default()))

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatal("Error connecting to database, please check env vars or values to connect")
	}

	return dbpool
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

var (
	replaceTabs                      = regexp.MustCompile(`\t+`)
	replaceSpacesBeforeOpeningParens = regexp.MustCompile(`\s+\(`)
	replaceSpacesAfterOpeningParens  = regexp.MustCompile(`\(\s+`)
	replaceSpacesBeforeClosingParens = regexp.MustCompile(`\s+\)`)
	replaceSpacesAfterClosingParens  = regexp.MustCompile(`\)\s+`)
	replaceSpaces                    = regexp.MustCompile(`\s+`)
)

// prettyPrintSQL removes empty lines and trims spaces.
func prettyPrintSQL(sql string) string {
	lines := strings.Split(sql, "\n")

	pretty := strings.Join(lines, " ")
	pretty = replaceTabs.ReplaceAllString(pretty, "")
	pretty = replaceSpacesBeforeOpeningParens.ReplaceAllString(pretty, "(")
	pretty = replaceSpacesAfterOpeningParens.ReplaceAllString(pretty, "(")
	pretty = replaceSpacesAfterClosingParens.ReplaceAllString(pretty, ")")
	pretty = replaceSpacesBeforeClosingParens.ReplaceAllString(pretty, ")")

	// Finally, replace multiple spaces with a single space
	pretty = replaceSpaces.ReplaceAllString(pretty, " ")

	return strings.TrimSpace(pretty)
}

// https://github.com/jackc/pgx/issues/1061#issuecomment-1186250809
type LoggingQueryTracer struct {
	logger *slog.Logger
}

func NewLoggingQueryTracer(logger *slog.Logger) *LoggingQueryTracer {
	return &LoggingQueryTracer{logger: logger}
}

func (l *LoggingQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	l.logger.
		Info("query start",
			slog.String("sql", prettyPrintSQL(data.SQL)),
			slog.Any("args", data.Args),
		)
	return ctx
}

func (l *LoggingQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// Failure
	if data.Err != nil {
		l.logger.
			Error("query end",
				slog.String("error", data.Err.Error()),
				slog.String("command_tag", data.CommandTag.String()),
			)
		return
	}

	// Success
	l.logger.
		Info("query end",
			slog.String("command_tag", data.CommandTag.String()),
		)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// https://github.com/jackc/pgx/discussions/1677#discussioncomment-8815982
type MultiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

func NewMultiQueryTracer(tracers ...pgx.QueryTracer) *MultiQueryTracer {
	return &MultiQueryTracer{Tracers: tracers}
}

func (m *MultiQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, t := range m.Tracers {
		ctx = t.TraceQueryStart(ctx, conn, data)
	}
	return ctx
}

func (m *MultiQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, t := range m.Tracers {
		t.TraceQueryEnd(ctx, conn, data)
	}
}
