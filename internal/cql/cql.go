// Package cql checks Edge Delta query syntax offline, using the same ANTLR grammars
// the backend parses queries with.
//
// The subdirectories are generated parsers vendored verbatim from
// pkg/antlrcql/{logparser,metric/parser,formula/parser} in the edgedelta monorepo.
// They depend only on the ANTLR runtime, so they compile here unmodified. Refresh them
// with `make sync-cql-parsers`.
//
// # What this does and does not check
//
// This is a syntax check: it reports what the grammar rejects and nothing more. The
// backend's semantic layer (pkg/antlrcql/{log,metric}) additionally resolves facet
// names, scopes and top-level columns against the org's data, which needs a live
// backend, so it cannot run here. A query that passes this check is well-formed; it may
// still reference a facet or metric that does not exist.
//
// Queries in a dashboard are templates containing `$variable` references, which the
// frontend substitutes before sending them. `$foo` lexes as an ordinary TERM, so
// templates parse as-is — but this validates the template, not the substituted result.
package cql

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	formulaparser "github.com/edgedelta/edx/internal/cql/formula/parser"
	"github.com/edgedelta/edx/internal/cql/logparser"
	metricparser "github.com/edgedelta/edx/internal/cql/metric/parser"
)

// Dialect selects a grammar. Which one a query needs depends on where it appears in a
// dashboard; see dashboards.QuerySites.
type Dialect string

const (
	// DialectLog is EDCqlLogParser: filters with optional group-by and rollup. The
	// backend parses log, event, pattern, trace and facet-option queries with it.
	DialectLog Dialect = "log"
	// DialectMetric is EDCqlMetricParser: `agg:metric{filter} by {...}.fill().rollup()`.
	DialectMetric Dialect = "metric"
	// DialectFormula is EDCqlFormulaParser: arithmetic over other queries' results.
	DialectFormula Dialect = "formula"
)

// SyntaxError is one problem the grammar found, located within the query string.
type SyntaxError struct {
	// Line is 1-based, Column 0-based, matching ANTLR.
	Line, Column int
	Message      string
}

func (e SyntaxError) String() string {
	return fmt.Sprintf("at column %d: %s", e.Column, e.Message)
}

// errorListener collects syntax errors instead of ANTLR's default of printing them to
// stderr and continuing.
type errorListener struct {
	*antlr.DefaultErrorListener
	errs []SyntaxError
}

func (l *errorListener) SyntaxError(_ antlr.Recognizer, _ any, line, column int, msg string, _ antlr.RecognitionException) {
	l.errs = append(l.errs, SyntaxError{Line: line, Column: column, Message: msg})
}

// Validate parses query with the given dialect's grammar and returns every syntax error
// found. An empty result means the query is well-formed.
//
// An empty query is valid: every query field in a dashboard definition is optional, and
// the backend treats a missing filter as "match all".
func Validate(dialect Dialect, query string) ([]SyntaxError, error) {
	if query == "" {
		return nil, nil
	}

	l := &errorListener{}
	input := antlr.NewInputStream(query)

	// Each grammar's entry rule ends in EOF, so trailing garbage is a syntax error
	// rather than silently ignored input.
	switch dialect {
	case DialectLog:
		lexer := logparser.NewEDCqlLexer(input)
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(l)
		p := logparser.NewEDCqlLogParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
		p.RemoveErrorListeners()
		p.AddErrorListener(l)
		p.FinalQuery()
	case DialectMetric:
		lexer := metricparser.NewEDCqlLexer(input)
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(l)
		p := metricparser.NewEDCqlMetricParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
		p.RemoveErrorListeners()
		p.AddErrorListener(l)
		p.FinalQuery()
	case DialectFormula:
		lexer := formulaparser.NewEDFormulaLexer(input)
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(l)
		p := formulaparser.NewEDCqlFormulaParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
		p.RemoveErrorListeners()
		p.AddErrorListener(l)
		p.Formula()
	default:
		// Returned rather than ignored: a dialect this package does not know about must
		// never read as "no problems found".
		return nil, fmt.Errorf("unknown query dialect %q", dialect)
	}

	return l.errs, nil
}
