package initializers

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

var (
	replaceTabs                      = regexp.MustCompile(`\t+`)
	replaceSpacesBeforeOpeningParens = regexp.MustCompile(`\s+\(`)
	replaceSpacesAfterOpeningParens  = regexp.MustCompile(`\(\s+`)
	replaceSpacesBeforeClosingParens = regexp.MustCompile(`\s+\)`)
	replaceSpacesAfterClosingParens  = regexp.MustCompile(`\)\s+`)
	replaceSpaces                    = regexp.MustCompile(`\s+`)
)

type LoggingQueryTracer struct {
	logger *zerolog.Logger
}

type MultiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

// ConnectDB creates a pool to connect into a postgres database
// and returns it. Please always remember to defer close it
// after using.
func ConnectDB(logger *zerolog.Logger) *pgxpool.Pool {
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
		logger.Error().Err(err).Msg("Error on the config string, please check env vars")
	}

	config.ConnConfig.Tracer = NewMultiQueryTracer(NewLoggingQueryTracer(logger))

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		logger.Error().Err(err).Msg("Error connecting to database, please check env vars or values to connect")
	}

	return dbpool
}

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
func NewLoggingQueryTracer(logger *zerolog.Logger) *LoggingQueryTracer {
	return &LoggingQueryTracer{logger: logger}
}

func (l *LoggingQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	l.logger.Log().
		Str("query", prettyPrintSQL(data.SQL)).
		Any("args", data.Args).
		Msg("Query Start")

	return ctx
}

func (l *LoggingQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// Failure
	if data.Err != nil {
		l.logger.Error().
			Str("return", data.CommandTag.String()).
			Str("error", data.Err.Error()).
			Msg("Query End")

		return
	}

	// Success
	l.logger.Log().
		Str("return", data.CommandTag.String()).
		Msg("Query End")
}

// https://github.com/jackc/pgx/discussions/1677#discussioncomment-8815982
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
