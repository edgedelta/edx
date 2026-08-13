// Code generated from EDFormulaLexer.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type EDFormulaLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var EDFormulaLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func edformulalexerLexerInit() {
	staticData := &EDFormulaLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'('", "')'", "','", "'timeshift'", "'moving_average'",
	}
	staticData.SymbolicNames = []string{
		"", "LPAREN", "RPAREN", "COMMA", "TIMESHIFT", "MOVING_AVERAGE", "NUMERIC",
		"FREE_TOKEN", "ARITHMETIC_TOKEN", "DEFAULT_SKIP", "UNKNOWN",
	}
	staticData.RuleNames = []string{
		"LPAREN", "RPAREN", "COMMA", "TIMESHIFT", "MOVING_AVERAGE", "NUMERIC",
		"FREE_TOKEN", "ARITHMETIC_TOKEN", "DEFAULT_SKIP", "UNKNOWN", "WHITESPACE",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 10, 80, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4,
		1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 3, 5, 56, 8,
		5, 1, 5, 4, 5, 59, 8, 5, 11, 5, 12, 5, 60, 1, 6, 4, 6, 64, 8, 6, 11, 6,
		12, 6, 65, 1, 7, 1, 7, 1, 8, 4, 8, 71, 8, 8, 11, 8, 12, 8, 72, 1, 8, 1,
		8, 1, 9, 1, 9, 1, 10, 1, 10, 0, 0, 11, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11,
		6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 0, 1, 0, 4, 1, 0, 48, 57, 9, 0, 9,
		10, 13, 13, 32, 32, 37, 38, 40, 45, 47, 47, 94, 94, 124, 124, 12288, 12288,
		6, 0, 37, 38, 42, 43, 45, 45, 47, 47, 94, 94, 124, 124, 4, 0, 9, 10, 13,
		13, 32, 32, 12288, 12288, 82, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5,
		1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13,
		1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 1,
		23, 1, 0, 0, 0, 3, 25, 1, 0, 0, 0, 5, 27, 1, 0, 0, 0, 7, 29, 1, 0, 0, 0,
		9, 39, 1, 0, 0, 0, 11, 55, 1, 0, 0, 0, 13, 63, 1, 0, 0, 0, 15, 67, 1, 0,
		0, 0, 17, 70, 1, 0, 0, 0, 19, 76, 1, 0, 0, 0, 21, 78, 1, 0, 0, 0, 23, 24,
		5, 40, 0, 0, 24, 2, 1, 0, 0, 0, 25, 26, 5, 41, 0, 0, 26, 4, 1, 0, 0, 0,
		27, 28, 5, 44, 0, 0, 28, 6, 1, 0, 0, 0, 29, 30, 5, 116, 0, 0, 30, 31, 5,
		105, 0, 0, 31, 32, 5, 109, 0, 0, 32, 33, 5, 101, 0, 0, 33, 34, 5, 115,
		0, 0, 34, 35, 5, 104, 0, 0, 35, 36, 5, 105, 0, 0, 36, 37, 5, 102, 0, 0,
		37, 38, 5, 116, 0, 0, 38, 8, 1, 0, 0, 0, 39, 40, 5, 109, 0, 0, 40, 41,
		5, 111, 0, 0, 41, 42, 5, 118, 0, 0, 42, 43, 5, 105, 0, 0, 43, 44, 5, 110,
		0, 0, 44, 45, 5, 103, 0, 0, 45, 46, 5, 95, 0, 0, 46, 47, 5, 97, 0, 0, 47,
		48, 5, 118, 0, 0, 48, 49, 5, 101, 0, 0, 49, 50, 5, 114, 0, 0, 50, 51, 5,
		97, 0, 0, 51, 52, 5, 103, 0, 0, 52, 53, 5, 101, 0, 0, 53, 10, 1, 0, 0,
		0, 54, 56, 5, 45, 0, 0, 55, 54, 1, 0, 0, 0, 55, 56, 1, 0, 0, 0, 56, 58,
		1, 0, 0, 0, 57, 59, 7, 0, 0, 0, 58, 57, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0,
		60, 58, 1, 0, 0, 0, 60, 61, 1, 0, 0, 0, 61, 12, 1, 0, 0, 0, 62, 64, 8,
		1, 0, 0, 63, 62, 1, 0, 0, 0, 64, 65, 1, 0, 0, 0, 65, 63, 1, 0, 0, 0, 65,
		66, 1, 0, 0, 0, 66, 14, 1, 0, 0, 0, 67, 68, 7, 2, 0, 0, 68, 16, 1, 0, 0,
		0, 69, 71, 3, 21, 10, 0, 70, 69, 1, 0, 0, 0, 71, 72, 1, 0, 0, 0, 72, 70,
		1, 0, 0, 0, 72, 73, 1, 0, 0, 0, 73, 74, 1, 0, 0, 0, 74, 75, 6, 8, 0, 0,
		75, 18, 1, 0, 0, 0, 76, 77, 9, 0, 0, 0, 77, 20, 1, 0, 0, 0, 78, 79, 7,
		3, 0, 0, 79, 22, 1, 0, 0, 0, 5, 0, 55, 60, 65, 72, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// EDFormulaLexerInit initializes any static state used to implement EDFormulaLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewEDFormulaLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func EDFormulaLexerInit() {
	staticData := &EDFormulaLexerLexerStaticData
	staticData.once.Do(edformulalexerLexerInit)
}

// NewEDFormulaLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewEDFormulaLexer(input antlr.CharStream) *EDFormulaLexer {
	EDFormulaLexerInit()
	l := new(EDFormulaLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &EDFormulaLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "EDFormulaLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// EDFormulaLexer tokens.
const (
	EDFormulaLexerLPAREN           = 1
	EDFormulaLexerRPAREN           = 2
	EDFormulaLexerCOMMA            = 3
	EDFormulaLexerTIMESHIFT        = 4
	EDFormulaLexerMOVING_AVERAGE   = 5
	EDFormulaLexerNUMERIC          = 6
	EDFormulaLexerFREE_TOKEN       = 7
	EDFormulaLexerARITHMETIC_TOKEN = 8
	EDFormulaLexerDEFAULT_SKIP     = 9
	EDFormulaLexerUNKNOWN          = 10
)
