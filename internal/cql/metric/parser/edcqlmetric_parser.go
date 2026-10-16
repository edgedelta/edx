// Code generated from EDCqlMetricParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlMetricParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type EDCqlMetricParser struct {
	*antlr.BaseParser
}

var EDCqlMetricParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func edcqlmetricparserParserInit() {
	staticData := &EDCqlMetricParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "'{'", "'}'", "'('", "')'", "", "", "", "'*'", "", "'.rollup'",
		"", "'count_unique'", "'.fill'", "'last'",
	}
	staticData.SymbolicNames = []string{
		"", "AND", "OR", "NOT", "LBRACE", "RBRACE", "LPAREN", "RPAREN", "OP_COLON",
		"OPERATOR", "REGEX_OPERATOR", "MATCH_ALL", "COMMA", "ROLLUP_START",
		"BY", "COUNT_UNIQUE", "FILL_START", "COMMON_OPERATOR", "FILL_OPERATOR",
		"AGGREGATION_OPERATOR", "COUNTER_OPERATOR", "QUOTED", "NUMBER", "TERM",
		"DEFAULT_SKIP", "UNKNOWN",
	}
	staticData.RuleNames = []string{
		"finalQuery", "topLevelQuery", "queryOperation", "metricName", "aggregationFilter",
		"aggregationMethod", "operatorColon", "groupByFields", "countUniqueFields",
		"groupBySection", "countUniqueSection", "rollupWindow", "rollupSection",
		"aggregation", "fillMethod", "fillLimit", "fillSection", "query", "disjQuery",
		"conjQuery", "modClause", "modifier", "clause", "term", "groupingExpr",
		"fieldName", "quotedTerm",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 25, 205, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 3, 2, 62, 8, 2, 1, 3,
		1, 3, 1, 4, 1, 4, 3, 4, 68, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 74, 8,
		5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 5, 7, 81, 8, 7, 10, 7, 12, 7, 84, 9, 7,
		1, 8, 1, 8, 1, 8, 5, 8, 89, 8, 8, 10, 8, 12, 8, 92, 9, 8, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 3, 10, 102, 8, 10, 1, 10, 3, 10, 105,
		8, 10, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1,
		13, 3, 13, 117, 8, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 124, 8,
		13, 1, 13, 3, 13, 127, 8, 13, 1, 13, 3, 13, 130, 8, 13, 1, 14, 1, 14, 1,
		15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 141, 8, 16, 1, 16,
		1, 16, 1, 17, 4, 17, 146, 8, 17, 11, 17, 12, 17, 147, 1, 18, 1, 18, 1,
		18, 5, 18, 153, 8, 18, 10, 18, 12, 18, 156, 9, 18, 1, 19, 1, 19, 3, 19,
		160, 8, 19, 1, 19, 5, 19, 163, 8, 19, 10, 19, 12, 19, 166, 9, 19, 1, 20,
		3, 20, 169, 8, 20, 1, 20, 1, 20, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 3,
		22, 178, 8, 22, 1, 22, 1, 22, 3, 22, 182, 8, 22, 1, 23, 1, 23, 1, 23, 1,
		23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 3, 23, 195, 8, 23,
		1, 24, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 26, 0, 0, 27,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36,
		38, 40, 42, 44, 46, 48, 50, 52, 0, 3, 2, 0, 20, 20, 23, 23, 1, 0, 8, 10,
		1, 0, 17, 18, 208, 0, 54, 1, 0, 0, 0, 2, 56, 1, 0, 0, 0, 4, 61, 1, 0, 0,
		0, 6, 63, 1, 0, 0, 0, 8, 67, 1, 0, 0, 0, 10, 73, 1, 0, 0, 0, 12, 75, 1,
		0, 0, 0, 14, 77, 1, 0, 0, 0, 16, 85, 1, 0, 0, 0, 18, 93, 1, 0, 0, 0, 20,
		98, 1, 0, 0, 0, 22, 106, 1, 0, 0, 0, 24, 108, 1, 0, 0, 0, 26, 116, 1, 0,
		0, 0, 28, 131, 1, 0, 0, 0, 30, 133, 1, 0, 0, 0, 32, 135, 1, 0, 0, 0, 34,
		145, 1, 0, 0, 0, 36, 149, 1, 0, 0, 0, 38, 157, 1, 0, 0, 0, 40, 168, 1,
		0, 0, 0, 42, 172, 1, 0, 0, 0, 44, 177, 1, 0, 0, 0, 46, 194, 1, 0, 0, 0,
		48, 196, 1, 0, 0, 0, 50, 200, 1, 0, 0, 0, 52, 202, 1, 0, 0, 0, 54, 55,
		3, 2, 1, 0, 55, 1, 1, 0, 0, 0, 56, 57, 3, 4, 2, 0, 57, 58, 5, 0, 0, 1,
		58, 3, 1, 0, 0, 0, 59, 62, 3, 26, 13, 0, 60, 62, 3, 34, 17, 0, 61, 59,
		1, 0, 0, 0, 61, 60, 1, 0, 0, 0, 62, 5, 1, 0, 0, 0, 63, 64, 7, 0, 0, 0,
		64, 7, 1, 0, 0, 0, 65, 68, 5, 11, 0, 0, 66, 68, 3, 34, 17, 0, 67, 65, 1,
		0, 0, 0, 67, 66, 1, 0, 0, 0, 68, 9, 1, 0, 0, 0, 69, 74, 5, 19, 0, 0, 70,
		74, 5, 17, 0, 0, 71, 74, 5, 20, 0, 0, 72, 74, 3, 20, 10, 0, 73, 69, 1,
		0, 0, 0, 73, 70, 1, 0, 0, 0, 73, 71, 1, 0, 0, 0, 73, 72, 1, 0, 0, 0, 74,
		11, 1, 0, 0, 0, 75, 76, 7, 1, 0, 0, 76, 13, 1, 0, 0, 0, 77, 82, 7, 0, 0,
		0, 78, 79, 5, 12, 0, 0, 79, 81, 7, 0, 0, 0, 80, 78, 1, 0, 0, 0, 81, 84,
		1, 0, 0, 0, 82, 80, 1, 0, 0, 0, 82, 83, 1, 0, 0, 0, 83, 15, 1, 0, 0, 0,
		84, 82, 1, 0, 0, 0, 85, 90, 7, 0, 0, 0, 86, 87, 5, 12, 0, 0, 87, 89, 7,
		0, 0, 0, 88, 86, 1, 0, 0, 0, 89, 92, 1, 0, 0, 0, 90, 88, 1, 0, 0, 0, 90,
		91, 1, 0, 0, 0, 91, 17, 1, 0, 0, 0, 92, 90, 1, 0, 0, 0, 93, 94, 5, 14,
		0, 0, 94, 95, 5, 4, 0, 0, 95, 96, 3, 14, 7, 0, 96, 97, 5, 5, 0, 0, 97,
		19, 1, 0, 0, 0, 98, 104, 5, 15, 0, 0, 99, 101, 5, 6, 0, 0, 100, 102, 3,
		16, 8, 0, 101, 100, 1, 0, 0, 0, 101, 102, 1, 0, 0, 0, 102, 103, 1, 0, 0,
		0, 103, 105, 5, 7, 0, 0, 104, 99, 1, 0, 0, 0, 104, 105, 1, 0, 0, 0, 105,
		21, 1, 0, 0, 0, 106, 107, 5, 22, 0, 0, 107, 23, 1, 0, 0, 0, 108, 109, 5,
		13, 0, 0, 109, 110, 5, 6, 0, 0, 110, 111, 3, 22, 11, 0, 111, 112, 5, 7,
		0, 0, 112, 25, 1, 0, 0, 0, 113, 114, 3, 10, 5, 0, 114, 115, 5, 8, 0, 0,
		115, 117, 1, 0, 0, 0, 116, 113, 1, 0, 0, 0, 116, 117, 1, 0, 0, 0, 117,
		118, 1, 0, 0, 0, 118, 119, 3, 6, 3, 0, 119, 120, 5, 4, 0, 0, 120, 121,
		3, 8, 4, 0, 121, 123, 5, 5, 0, 0, 122, 124, 3, 18, 9, 0, 123, 122, 1, 0,
		0, 0, 123, 124, 1, 0, 0, 0, 124, 126, 1, 0, 0, 0, 125, 127, 3, 32, 16,
		0, 126, 125, 1, 0, 0, 0, 126, 127, 1, 0, 0, 0, 127, 129, 1, 0, 0, 0, 128,
		130, 3, 24, 12, 0, 129, 128, 1, 0, 0, 0, 129, 130, 1, 0, 0, 0, 130, 27,
		1, 0, 0, 0, 131, 132, 7, 2, 0, 0, 132, 29, 1, 0, 0, 0, 133, 134, 5, 22,
		0, 0, 134, 31, 1, 0, 0, 0, 135, 136, 5, 16, 0, 0, 136, 137, 5, 6, 0, 0,
		137, 140, 3, 28, 14, 0, 138, 139, 5, 12, 0, 0, 139, 141, 3, 30, 15, 0,
		140, 138, 1, 0, 0, 0, 140, 141, 1, 0, 0, 0, 141, 142, 1, 0, 0, 0, 142,
		143, 5, 7, 0, 0, 143, 33, 1, 0, 0, 0, 144, 146, 3, 36, 18, 0, 145, 144,
		1, 0, 0, 0, 146, 147, 1, 0, 0, 0, 147, 145, 1, 0, 0, 0, 147, 148, 1, 0,
		0, 0, 148, 35, 1, 0, 0, 0, 149, 154, 3, 38, 19, 0, 150, 151, 5, 2, 0, 0,
		151, 153, 3, 38, 19, 0, 152, 150, 1, 0, 0, 0, 153, 156, 1, 0, 0, 0, 154,
		152, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0, 155, 37, 1, 0, 0, 0, 156, 154, 1,
		0, 0, 0, 157, 164, 3, 40, 20, 0, 158, 160, 5, 1, 0, 0, 159, 158, 1, 0,
		0, 0, 159, 160, 1, 0, 0, 0, 160, 161, 1, 0, 0, 0, 161, 163, 3, 40, 20,
		0, 162, 159, 1, 0, 0, 0, 163, 166, 1, 0, 0, 0, 164, 162, 1, 0, 0, 0, 164,
		165, 1, 0, 0, 0, 165, 39, 1, 0, 0, 0, 166, 164, 1, 0, 0, 0, 167, 169, 3,
		42, 21, 0, 168, 167, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 170, 1, 0,
		0, 0, 170, 171, 3, 44, 22, 0, 171, 41, 1, 0, 0, 0, 172, 173, 5, 3, 0, 0,
		173, 43, 1, 0, 0, 0, 174, 175, 3, 50, 25, 0, 175, 176, 3, 12, 6, 0, 176,
		178, 1, 0, 0, 0, 177, 174, 1, 0, 0, 0, 177, 178, 1, 0, 0, 0, 178, 181,
		1, 0, 0, 0, 179, 182, 3, 46, 23, 0, 180, 182, 3, 48, 24, 0, 181, 179, 1,
		0, 0, 0, 181, 180, 1, 0, 0, 0, 182, 45, 1, 0, 0, 0, 183, 195, 3, 52, 26,
		0, 184, 195, 5, 19, 0, 0, 185, 195, 5, 20, 0, 0, 186, 195, 5, 17, 0, 0,
		187, 195, 5, 18, 0, 0, 188, 195, 5, 15, 0, 0, 189, 195, 5, 13, 0, 0, 190,
		195, 5, 16, 0, 0, 191, 195, 5, 14, 0, 0, 192, 195, 5, 22, 0, 0, 193, 195,
		5, 23, 0, 0, 194, 183, 1, 0, 0, 0, 194, 184, 1, 0, 0, 0, 194, 185, 1, 0,
		0, 0, 194, 186, 1, 0, 0, 0, 194, 187, 1, 0, 0, 0, 194, 188, 1, 0, 0, 0,
		194, 189, 1, 0, 0, 0, 194, 190, 1, 0, 0, 0, 194, 191, 1, 0, 0, 0, 194,
		192, 1, 0, 0, 0, 194, 193, 1, 0, 0, 0, 195, 47, 1, 0, 0, 0, 196, 197, 5,
		6, 0, 0, 197, 198, 3, 34, 17, 0, 198, 199, 5, 7, 0, 0, 199, 49, 1, 0, 0,
		0, 200, 201, 7, 0, 0, 0, 201, 51, 1, 0, 0, 0, 202, 203, 5, 21, 0, 0, 203,
		53, 1, 0, 0, 0, 20, 61, 67, 73, 82, 90, 101, 104, 116, 123, 126, 129, 140,
		147, 154, 159, 164, 168, 177, 181, 194,
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

// EDCqlMetricParserInit initializes any static state used to implement EDCqlMetricParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewEDCqlMetricParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func EDCqlMetricParserInit() {
	staticData := &EDCqlMetricParserParserStaticData
	staticData.once.Do(edcqlmetricparserParserInit)
}

// NewEDCqlMetricParser produces a new parser instance for the optional input antlr.TokenStream.
func NewEDCqlMetricParser(input antlr.TokenStream) *EDCqlMetricParser {
	EDCqlMetricParserInit()
	this := new(EDCqlMetricParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &EDCqlMetricParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "EDCqlMetricParser.g4"

	return this
}

// EDCqlMetricParser tokens.
const (
	EDCqlMetricParserEOF                  = antlr.TokenEOF
	EDCqlMetricParserAND                  = 1
	EDCqlMetricParserOR                   = 2
	EDCqlMetricParserNOT                  = 3
	EDCqlMetricParserLBRACE               = 4
	EDCqlMetricParserRBRACE               = 5
	EDCqlMetricParserLPAREN               = 6
	EDCqlMetricParserRPAREN               = 7
	EDCqlMetricParserOP_COLON             = 8
	EDCqlMetricParserOPERATOR             = 9
	EDCqlMetricParserREGEX_OPERATOR       = 10
	EDCqlMetricParserMATCH_ALL            = 11
	EDCqlMetricParserCOMMA                = 12
	EDCqlMetricParserROLLUP_START         = 13
	EDCqlMetricParserBY                   = 14
	EDCqlMetricParserCOUNT_UNIQUE         = 15
	EDCqlMetricParserFILL_START           = 16
	EDCqlMetricParserCOMMON_OPERATOR      = 17
	EDCqlMetricParserFILL_OPERATOR        = 18
	EDCqlMetricParserAGGREGATION_OPERATOR = 19
	EDCqlMetricParserCOUNTER_OPERATOR     = 20
	EDCqlMetricParserQUOTED               = 21
	EDCqlMetricParserNUMBER               = 22
	EDCqlMetricParserTERM                 = 23
	EDCqlMetricParserDEFAULT_SKIP         = 24
	EDCqlMetricParserUNKNOWN              = 25
)

// EDCqlMetricParser rules.
const (
	EDCqlMetricParserRULE_finalQuery         = 0
	EDCqlMetricParserRULE_topLevelQuery      = 1
	EDCqlMetricParserRULE_queryOperation     = 2
	EDCqlMetricParserRULE_metricName         = 3
	EDCqlMetricParserRULE_aggregationFilter  = 4
	EDCqlMetricParserRULE_aggregationMethod  = 5
	EDCqlMetricParserRULE_operatorColon      = 6
	EDCqlMetricParserRULE_groupByFields      = 7
	EDCqlMetricParserRULE_countUniqueFields  = 8
	EDCqlMetricParserRULE_groupBySection     = 9
	EDCqlMetricParserRULE_countUniqueSection = 10
	EDCqlMetricParserRULE_rollupWindow       = 11
	EDCqlMetricParserRULE_rollupSection      = 12
	EDCqlMetricParserRULE_aggregation        = 13
	EDCqlMetricParserRULE_fillMethod         = 14
	EDCqlMetricParserRULE_fillLimit          = 15
	EDCqlMetricParserRULE_fillSection        = 16
	EDCqlMetricParserRULE_query              = 17
	EDCqlMetricParserRULE_disjQuery          = 18
	EDCqlMetricParserRULE_conjQuery          = 19
	EDCqlMetricParserRULE_modClause          = 20
	EDCqlMetricParserRULE_modifier           = 21
	EDCqlMetricParserRULE_clause             = 22
	EDCqlMetricParserRULE_term               = 23
	EDCqlMetricParserRULE_groupingExpr       = 24
	EDCqlMetricParserRULE_fieldName          = 25
	EDCqlMetricParserRULE_quotedTerm         = 26
)

// IFinalQueryContext is an interface to support dynamic dispatch.
type IFinalQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TopLevelQuery() ITopLevelQueryContext

	// IsFinalQueryContext differentiates from other interfaces.
	IsFinalQueryContext()
}

type FinalQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFinalQueryContext() *FinalQueryContext {
	var p = new(FinalQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_finalQuery
	return p
}

func InitEmptyFinalQueryContext(p *FinalQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_finalQuery
}

func (*FinalQueryContext) IsFinalQueryContext() {}

func NewFinalQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FinalQueryContext {
	var p = new(FinalQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_finalQuery

	return p
}

func (s *FinalQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *FinalQueryContext) TopLevelQuery() ITopLevelQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITopLevelQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITopLevelQueryContext)
}

func (s *FinalQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FinalQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FinalQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterFinalQuery(s)
	}
}

func (s *FinalQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitFinalQuery(s)
	}
}

func (s *FinalQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitFinalQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) FinalQuery() (localctx IFinalQueryContext) {
	localctx = NewFinalQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, EDCqlMetricParserRULE_finalQuery)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(54)
		p.TopLevelQuery()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITopLevelQueryContext is an interface to support dynamic dispatch.
type ITopLevelQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QueryOperation() IQueryOperationContext
	EOF() antlr.TerminalNode

	// IsTopLevelQueryContext differentiates from other interfaces.
	IsTopLevelQueryContext()
}

type TopLevelQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTopLevelQueryContext() *TopLevelQueryContext {
	var p = new(TopLevelQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_topLevelQuery
	return p
}

func InitEmptyTopLevelQueryContext(p *TopLevelQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_topLevelQuery
}

func (*TopLevelQueryContext) IsTopLevelQueryContext() {}

func NewTopLevelQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelQueryContext {
	var p = new(TopLevelQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_topLevelQuery

	return p
}

func (s *TopLevelQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *TopLevelQueryContext) QueryOperation() IQueryOperationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryOperationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryOperationContext)
}

func (s *TopLevelQueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserEOF, 0)
}

func (s *TopLevelQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TopLevelQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TopLevelQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterTopLevelQuery(s)
	}
}

func (s *TopLevelQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitTopLevelQuery(s)
	}
}

func (s *TopLevelQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitTopLevelQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) TopLevelQuery() (localctx ITopLevelQueryContext) {
	localctx = NewTopLevelQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, EDCqlMetricParserRULE_topLevelQuery)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(56)
		p.QueryOperation()
	}
	{
		p.SetState(57)
		p.Match(EDCqlMetricParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQueryOperationContext is an interface to support dynamic dispatch.
type IQueryOperationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Aggregation() IAggregationContext
	Query() IQueryContext

	// IsQueryOperationContext differentiates from other interfaces.
	IsQueryOperationContext()
}

type QueryOperationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryOperationContext() *QueryOperationContext {
	var p = new(QueryOperationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_queryOperation
	return p
}

func InitEmptyQueryOperationContext(p *QueryOperationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_queryOperation
}

func (*QueryOperationContext) IsQueryOperationContext() {}

func NewQueryOperationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryOperationContext {
	var p = new(QueryOperationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_queryOperation

	return p
}

func (s *QueryOperationContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryOperationContext) Aggregation() IAggregationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationContext)
}

func (s *QueryOperationContext) Query() IQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryContext)
}

func (s *QueryOperationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryOperationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryOperationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterQueryOperation(s)
	}
}

func (s *QueryOperationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitQueryOperation(s)
	}
}

func (s *QueryOperationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitQueryOperation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) QueryOperation() (localctx IQueryOperationContext) {
	localctx = NewQueryOperationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, EDCqlMetricParserRULE_queryOperation)
	p.SetState(61)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(59)
			p.Aggregation()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(60)
			p.Query()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMetricNameContext is an interface to support dynamic dispatch.
type IMetricNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TERM() antlr.TerminalNode
	COUNTER_OPERATOR() antlr.TerminalNode

	// IsMetricNameContext differentiates from other interfaces.
	IsMetricNameContext()
}

type MetricNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetricNameContext() *MetricNameContext {
	var p = new(MetricNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_metricName
	return p
}

func InitEmptyMetricNameContext(p *MetricNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_metricName
}

func (*MetricNameContext) IsMetricNameContext() {}

func NewMetricNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetricNameContext {
	var p = new(MetricNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_metricName

	return p
}

func (s *MetricNameContext) GetParser() antlr.Parser { return s.parser }

func (s *MetricNameContext) TERM() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserTERM, 0)
}

func (s *MetricNameContext) COUNTER_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNTER_OPERATOR, 0)
}

func (s *MetricNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetricNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MetricNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterMetricName(s)
	}
}

func (s *MetricNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitMetricName(s)
	}
}

func (s *MetricNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitMetricName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) MetricName() (localctx IMetricNameContext) {
	localctx = NewMetricNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, EDCqlMetricParserRULE_metricName)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(63)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EDCqlMetricParserCOUNTER_OPERATOR || _la == EDCqlMetricParserTERM) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAggregationFilterContext is an interface to support dynamic dispatch.
type IAggregationFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MATCH_ALL() antlr.TerminalNode
	Query() IQueryContext

	// IsAggregationFilterContext differentiates from other interfaces.
	IsAggregationFilterContext()
}

type AggregationFilterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationFilterContext() *AggregationFilterContext {
	var p = new(AggregationFilterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_aggregationFilter
	return p
}

func InitEmptyAggregationFilterContext(p *AggregationFilterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_aggregationFilter
}

func (*AggregationFilterContext) IsAggregationFilterContext() {}

func NewAggregationFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationFilterContext {
	var p = new(AggregationFilterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_aggregationFilter

	return p
}

func (s *AggregationFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationFilterContext) MATCH_ALL() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserMATCH_ALL, 0)
}

func (s *AggregationFilterContext) Query() IQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryContext)
}

func (s *AggregationFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AggregationFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterAggregationFilter(s)
	}
}

func (s *AggregationFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitAggregationFilter(s)
	}
}

func (s *AggregationFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitAggregationFilter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) AggregationFilter() (localctx IAggregationFilterContext) {
	localctx = NewAggregationFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, EDCqlMetricParserRULE_aggregationFilter)
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlMetricParserMATCH_ALL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(65)
			p.Match(EDCqlMetricParserMATCH_ALL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserNOT, EDCqlMetricParserLPAREN, EDCqlMetricParserROLLUP_START, EDCqlMetricParserBY, EDCqlMetricParserCOUNT_UNIQUE, EDCqlMetricParserFILL_START, EDCqlMetricParserCOMMON_OPERATOR, EDCqlMetricParserFILL_OPERATOR, EDCqlMetricParserAGGREGATION_OPERATOR, EDCqlMetricParserCOUNTER_OPERATOR, EDCqlMetricParserQUOTED, EDCqlMetricParserNUMBER, EDCqlMetricParserTERM:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(66)
			p.Query()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAggregationMethodContext is an interface to support dynamic dispatch.
type IAggregationMethodContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AGGREGATION_OPERATOR() antlr.TerminalNode
	COMMON_OPERATOR() antlr.TerminalNode
	COUNTER_OPERATOR() antlr.TerminalNode
	CountUniqueSection() ICountUniqueSectionContext

	// IsAggregationMethodContext differentiates from other interfaces.
	IsAggregationMethodContext()
}

type AggregationMethodContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationMethodContext() *AggregationMethodContext {
	var p = new(AggregationMethodContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_aggregationMethod
	return p
}

func InitEmptyAggregationMethodContext(p *AggregationMethodContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_aggregationMethod
}

func (*AggregationMethodContext) IsAggregationMethodContext() {}

func NewAggregationMethodContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationMethodContext {
	var p = new(AggregationMethodContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_aggregationMethod

	return p
}

func (s *AggregationMethodContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationMethodContext) AGGREGATION_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserAGGREGATION_OPERATOR, 0)
}

func (s *AggregationMethodContext) COMMON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOMMON_OPERATOR, 0)
}

func (s *AggregationMethodContext) COUNTER_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNTER_OPERATOR, 0)
}

func (s *AggregationMethodContext) CountUniqueSection() ICountUniqueSectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICountUniqueSectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICountUniqueSectionContext)
}

func (s *AggregationMethodContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationMethodContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AggregationMethodContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterAggregationMethod(s)
	}
}

func (s *AggregationMethodContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitAggregationMethod(s)
	}
}

func (s *AggregationMethodContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitAggregationMethod(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) AggregationMethod() (localctx IAggregationMethodContext) {
	localctx = NewAggregationMethodContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, EDCqlMetricParserRULE_aggregationMethod)
	p.SetState(73)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlMetricParserAGGREGATION_OPERATOR:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(69)
			p.Match(EDCqlMetricParserAGGREGATION_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserCOMMON_OPERATOR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(70)
			p.Match(EDCqlMetricParserCOMMON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserCOUNTER_OPERATOR:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(71)
			p.Match(EDCqlMetricParserCOUNTER_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserCOUNT_UNIQUE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(72)
			p.CountUniqueSection()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOperatorColonContext is an interface to support dynamic dispatch.
type IOperatorColonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode
	OP_COLON() antlr.TerminalNode
	REGEX_OPERATOR() antlr.TerminalNode

	// IsOperatorColonContext differentiates from other interfaces.
	IsOperatorColonContext()
}

type OperatorColonContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperatorColonContext() *OperatorColonContext {
	var p = new(OperatorColonContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_operatorColon
	return p
}

func InitEmptyOperatorColonContext(p *OperatorColonContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_operatorColon
}

func (*OperatorColonContext) IsOperatorColonContext() {}

func NewOperatorColonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperatorColonContext {
	var p = new(OperatorColonContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_operatorColon

	return p
}

func (s *OperatorColonContext) GetParser() antlr.Parser { return s.parser }

func (s *OperatorColonContext) OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserOPERATOR, 0)
}

func (s *OperatorColonContext) OP_COLON() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserOP_COLON, 0)
}

func (s *OperatorColonContext) REGEX_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserREGEX_OPERATOR, 0)
}

func (s *OperatorColonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperatorColonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperatorColonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterOperatorColon(s)
	}
}

func (s *OperatorColonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitOperatorColon(s)
	}
}

func (s *OperatorColonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitOperatorColon(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) OperatorColon() (localctx IOperatorColonContext) {
	localctx = NewOperatorColonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, EDCqlMetricParserRULE_operatorColon)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1792) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGroupByFieldsContext is an interface to support dynamic dispatch.
type IGroupByFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTERM() []antlr.TerminalNode
	TERM(i int) antlr.TerminalNode
	AllCOUNTER_OPERATOR() []antlr.TerminalNode
	COUNTER_OPERATOR(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGroupByFieldsContext differentiates from other interfaces.
	IsGroupByFieldsContext()
}

type GroupByFieldsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByFieldsContext() *GroupByFieldsContext {
	var p = new(GroupByFieldsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_groupByFields
	return p
}

func InitEmptyGroupByFieldsContext(p *GroupByFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_groupByFields
}

func (*GroupByFieldsContext) IsGroupByFieldsContext() {}

func NewGroupByFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByFieldsContext {
	var p = new(GroupByFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_groupByFields

	return p
}

func (s *GroupByFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByFieldsContext) AllTERM() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserTERM)
}

func (s *GroupByFieldsContext) TERM(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserTERM, i)
}

func (s *GroupByFieldsContext) AllCOUNTER_OPERATOR() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserCOUNTER_OPERATOR)
}

func (s *GroupByFieldsContext) COUNTER_OPERATOR(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNTER_OPERATOR, i)
}

func (s *GroupByFieldsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserCOMMA)
}

func (s *GroupByFieldsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOMMA, i)
}

func (s *GroupByFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupByFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterGroupByFields(s)
	}
}

func (s *GroupByFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitGroupByFields(s)
	}
}

func (s *GroupByFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitGroupByFields(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) GroupByFields() (localctx IGroupByFieldsContext) {
	localctx = NewGroupByFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, EDCqlMetricParserRULE_groupByFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(77)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EDCqlMetricParserCOUNTER_OPERATOR || _la == EDCqlMetricParserTERM) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(82)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EDCqlMetricParserCOMMA {
		{
			p.SetState(78)
			p.Match(EDCqlMetricParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(79)
			_la = p.GetTokenStream().LA(1)

			if !(_la == EDCqlMetricParserCOUNTER_OPERATOR || _la == EDCqlMetricParserTERM) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(84)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICountUniqueFieldsContext is an interface to support dynamic dispatch.
type ICountUniqueFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTERM() []antlr.TerminalNode
	TERM(i int) antlr.TerminalNode
	AllCOUNTER_OPERATOR() []antlr.TerminalNode
	COUNTER_OPERATOR(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsCountUniqueFieldsContext differentiates from other interfaces.
	IsCountUniqueFieldsContext()
}

type CountUniqueFieldsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCountUniqueFieldsContext() *CountUniqueFieldsContext {
	var p = new(CountUniqueFieldsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_countUniqueFields
	return p
}

func InitEmptyCountUniqueFieldsContext(p *CountUniqueFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_countUniqueFields
}

func (*CountUniqueFieldsContext) IsCountUniqueFieldsContext() {}

func NewCountUniqueFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CountUniqueFieldsContext {
	var p = new(CountUniqueFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_countUniqueFields

	return p
}

func (s *CountUniqueFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *CountUniqueFieldsContext) AllTERM() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserTERM)
}

func (s *CountUniqueFieldsContext) TERM(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserTERM, i)
}

func (s *CountUniqueFieldsContext) AllCOUNTER_OPERATOR() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserCOUNTER_OPERATOR)
}

func (s *CountUniqueFieldsContext) COUNTER_OPERATOR(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNTER_OPERATOR, i)
}

func (s *CountUniqueFieldsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserCOMMA)
}

func (s *CountUniqueFieldsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOMMA, i)
}

func (s *CountUniqueFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CountUniqueFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CountUniqueFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterCountUniqueFields(s)
	}
}

func (s *CountUniqueFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitCountUniqueFields(s)
	}
}

func (s *CountUniqueFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitCountUniqueFields(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) CountUniqueFields() (localctx ICountUniqueFieldsContext) {
	localctx = NewCountUniqueFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, EDCqlMetricParserRULE_countUniqueFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EDCqlMetricParserCOUNTER_OPERATOR || _la == EDCqlMetricParserTERM) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EDCqlMetricParserCOMMA {
		{
			p.SetState(86)
			p.Match(EDCqlMetricParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(87)
			_la = p.GetTokenStream().LA(1)

			if !(_la == EDCqlMetricParserCOUNTER_OPERATOR || _la == EDCqlMetricParserTERM) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(92)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGroupBySectionContext is an interface to support dynamic dispatch.
type IGroupBySectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BY() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	GroupByFields() IGroupByFieldsContext
	RBRACE() antlr.TerminalNode

	// IsGroupBySectionContext differentiates from other interfaces.
	IsGroupBySectionContext()
}

type GroupBySectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupBySectionContext() *GroupBySectionContext {
	var p = new(GroupBySectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_groupBySection
	return p
}

func InitEmptyGroupBySectionContext(p *GroupBySectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_groupBySection
}

func (*GroupBySectionContext) IsGroupBySectionContext() {}

func NewGroupBySectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupBySectionContext {
	var p = new(GroupBySectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_groupBySection

	return p
}

func (s *GroupBySectionContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupBySectionContext) BY() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserBY, 0)
}

func (s *GroupBySectionContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserLBRACE, 0)
}

func (s *GroupBySectionContext) GroupByFields() IGroupByFieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByFieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByFieldsContext)
}

func (s *GroupBySectionContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserRBRACE, 0)
}

func (s *GroupBySectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupBySectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupBySectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterGroupBySection(s)
	}
}

func (s *GroupBySectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitGroupBySection(s)
	}
}

func (s *GroupBySectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitGroupBySection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) GroupBySection() (localctx IGroupBySectionContext) {
	localctx = NewGroupBySectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, EDCqlMetricParserRULE_groupBySection)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(93)
		p.Match(EDCqlMetricParserBY)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(94)
		p.Match(EDCqlMetricParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(95)
		p.GroupByFields()
	}
	{
		p.SetState(96)
		p.Match(EDCqlMetricParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICountUniqueSectionContext is an interface to support dynamic dispatch.
type ICountUniqueSectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COUNT_UNIQUE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	CountUniqueFields() ICountUniqueFieldsContext

	// IsCountUniqueSectionContext differentiates from other interfaces.
	IsCountUniqueSectionContext()
}

type CountUniqueSectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCountUniqueSectionContext() *CountUniqueSectionContext {
	var p = new(CountUniqueSectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_countUniqueSection
	return p
}

func InitEmptyCountUniqueSectionContext(p *CountUniqueSectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_countUniqueSection
}

func (*CountUniqueSectionContext) IsCountUniqueSectionContext() {}

func NewCountUniqueSectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CountUniqueSectionContext {
	var p = new(CountUniqueSectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_countUniqueSection

	return p
}

func (s *CountUniqueSectionContext) GetParser() antlr.Parser { return s.parser }

func (s *CountUniqueSectionContext) COUNT_UNIQUE() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNT_UNIQUE, 0)
}

func (s *CountUniqueSectionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserLPAREN, 0)
}

func (s *CountUniqueSectionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserRPAREN, 0)
}

func (s *CountUniqueSectionContext) CountUniqueFields() ICountUniqueFieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICountUniqueFieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICountUniqueFieldsContext)
}

func (s *CountUniqueSectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CountUniqueSectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CountUniqueSectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterCountUniqueSection(s)
	}
}

func (s *CountUniqueSectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitCountUniqueSection(s)
	}
}

func (s *CountUniqueSectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitCountUniqueSection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) CountUniqueSection() (localctx ICountUniqueSectionContext) {
	localctx = NewCountUniqueSectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, EDCqlMetricParserRULE_countUniqueSection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(98)
		p.Match(EDCqlMetricParserCOUNT_UNIQUE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(104)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlMetricParserLPAREN {
		{
			p.SetState(99)
			p.Match(EDCqlMetricParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(101)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == EDCqlMetricParserCOUNTER_OPERATOR || _la == EDCqlMetricParserTERM {
			{
				p.SetState(100)
				p.CountUniqueFields()
			}

		}
		{
			p.SetState(103)
			p.Match(EDCqlMetricParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRollupWindowContext is an interface to support dynamic dispatch.
type IRollupWindowContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode

	// IsRollupWindowContext differentiates from other interfaces.
	IsRollupWindowContext()
}

type RollupWindowContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRollupWindowContext() *RollupWindowContext {
	var p = new(RollupWindowContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_rollupWindow
	return p
}

func InitEmptyRollupWindowContext(p *RollupWindowContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_rollupWindow
}

func (*RollupWindowContext) IsRollupWindowContext() {}

func NewRollupWindowContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RollupWindowContext {
	var p = new(RollupWindowContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_rollupWindow

	return p
}

func (s *RollupWindowContext) GetParser() antlr.Parser { return s.parser }

func (s *RollupWindowContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserNUMBER, 0)
}

func (s *RollupWindowContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RollupWindowContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RollupWindowContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterRollupWindow(s)
	}
}

func (s *RollupWindowContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitRollupWindow(s)
	}
}

func (s *RollupWindowContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitRollupWindow(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) RollupWindow() (localctx IRollupWindowContext) {
	localctx = NewRollupWindowContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, EDCqlMetricParserRULE_rollupWindow)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(106)
		p.Match(EDCqlMetricParserNUMBER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRollupSectionContext is an interface to support dynamic dispatch.
type IRollupSectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ROLLUP_START() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RollupWindow() IRollupWindowContext
	RPAREN() antlr.TerminalNode

	// IsRollupSectionContext differentiates from other interfaces.
	IsRollupSectionContext()
}

type RollupSectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRollupSectionContext() *RollupSectionContext {
	var p = new(RollupSectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_rollupSection
	return p
}

func InitEmptyRollupSectionContext(p *RollupSectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_rollupSection
}

func (*RollupSectionContext) IsRollupSectionContext() {}

func NewRollupSectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RollupSectionContext {
	var p = new(RollupSectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_rollupSection

	return p
}

func (s *RollupSectionContext) GetParser() antlr.Parser { return s.parser }

func (s *RollupSectionContext) ROLLUP_START() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserROLLUP_START, 0)
}

func (s *RollupSectionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserLPAREN, 0)
}

func (s *RollupSectionContext) RollupWindow() IRollupWindowContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRollupWindowContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRollupWindowContext)
}

func (s *RollupSectionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserRPAREN, 0)
}

func (s *RollupSectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RollupSectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RollupSectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterRollupSection(s)
	}
}

func (s *RollupSectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitRollupSection(s)
	}
}

func (s *RollupSectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitRollupSection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) RollupSection() (localctx IRollupSectionContext) {
	localctx = NewRollupSectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, EDCqlMetricParserRULE_rollupSection)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(108)
		p.Match(EDCqlMetricParserROLLUP_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(109)
		p.Match(EDCqlMetricParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(110)
		p.RollupWindow()
	}
	{
		p.SetState(111)
		p.Match(EDCqlMetricParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAggregationContext is an interface to support dynamic dispatch.
type IAggregationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MetricName() IMetricNameContext
	LBRACE() antlr.TerminalNode
	AggregationFilter() IAggregationFilterContext
	RBRACE() antlr.TerminalNode
	AggregationMethod() IAggregationMethodContext
	OP_COLON() antlr.TerminalNode
	GroupBySection() IGroupBySectionContext
	FillSection() IFillSectionContext
	RollupSection() IRollupSectionContext

	// IsAggregationContext differentiates from other interfaces.
	IsAggregationContext()
}

type AggregationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationContext() *AggregationContext {
	var p = new(AggregationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_aggregation
	return p
}

func InitEmptyAggregationContext(p *AggregationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_aggregation
}

func (*AggregationContext) IsAggregationContext() {}

func NewAggregationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationContext {
	var p = new(AggregationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_aggregation

	return p
}

func (s *AggregationContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationContext) MetricName() IMetricNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMetricNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *AggregationContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserLBRACE, 0)
}

func (s *AggregationContext) AggregationFilter() IAggregationFilterContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationFilterContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationFilterContext)
}

func (s *AggregationContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserRBRACE, 0)
}

func (s *AggregationContext) AggregationMethod() IAggregationMethodContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationMethodContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationMethodContext)
}

func (s *AggregationContext) OP_COLON() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserOP_COLON, 0)
}

func (s *AggregationContext) GroupBySection() IGroupBySectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupBySectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupBySectionContext)
}

func (s *AggregationContext) FillSection() IFillSectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFillSectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFillSectionContext)
}

func (s *AggregationContext) RollupSection() IRollupSectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRollupSectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRollupSectionContext)
}

func (s *AggregationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AggregationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterAggregation(s)
	}
}

func (s *AggregationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitAggregation(s)
	}
}

func (s *AggregationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitAggregation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) Aggregation() (localctx IAggregationContext) {
	localctx = NewAggregationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, EDCqlMetricParserRULE_aggregation)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(116)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(113)
			p.AggregationMethod()
		}
		{
			p.SetState(114)
			p.Match(EDCqlMetricParserOP_COLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(118)
		p.MetricName()
	}
	{
		p.SetState(119)
		p.Match(EDCqlMetricParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(120)
		p.AggregationFilter()
	}
	{
		p.SetState(121)
		p.Match(EDCqlMetricParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(123)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlMetricParserBY {
		{
			p.SetState(122)
			p.GroupBySection()
		}

	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlMetricParserFILL_START {
		{
			p.SetState(125)
			p.FillSection()
		}

	}
	p.SetState(129)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlMetricParserROLLUP_START {
		{
			p.SetState(128)
			p.RollupSection()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFillMethodContext is an interface to support dynamic dispatch.
type IFillMethodContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FILL_OPERATOR() antlr.TerminalNode
	COMMON_OPERATOR() antlr.TerminalNode

	// IsFillMethodContext differentiates from other interfaces.
	IsFillMethodContext()
}

type FillMethodContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillMethodContext() *FillMethodContext {
	var p = new(FillMethodContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fillMethod
	return p
}

func InitEmptyFillMethodContext(p *FillMethodContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fillMethod
}

func (*FillMethodContext) IsFillMethodContext() {}

func NewFillMethodContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillMethodContext {
	var p = new(FillMethodContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_fillMethod

	return p
}

func (s *FillMethodContext) GetParser() antlr.Parser { return s.parser }

func (s *FillMethodContext) FILL_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserFILL_OPERATOR, 0)
}

func (s *FillMethodContext) COMMON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOMMON_OPERATOR, 0)
}

func (s *FillMethodContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillMethodContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FillMethodContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterFillMethod(s)
	}
}

func (s *FillMethodContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitFillMethod(s)
	}
}

func (s *FillMethodContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitFillMethod(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) FillMethod() (localctx IFillMethodContext) {
	localctx = NewFillMethodContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, EDCqlMetricParserRULE_fillMethod)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(131)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EDCqlMetricParserCOMMON_OPERATOR || _la == EDCqlMetricParserFILL_OPERATOR) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFillLimitContext is an interface to support dynamic dispatch.
type IFillLimitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode

	// IsFillLimitContext differentiates from other interfaces.
	IsFillLimitContext()
}

type FillLimitContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillLimitContext() *FillLimitContext {
	var p = new(FillLimitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fillLimit
	return p
}

func InitEmptyFillLimitContext(p *FillLimitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fillLimit
}

func (*FillLimitContext) IsFillLimitContext() {}

func NewFillLimitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillLimitContext {
	var p = new(FillLimitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_fillLimit

	return p
}

func (s *FillLimitContext) GetParser() antlr.Parser { return s.parser }

func (s *FillLimitContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserNUMBER, 0)
}

func (s *FillLimitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillLimitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FillLimitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterFillLimit(s)
	}
}

func (s *FillLimitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitFillLimit(s)
	}
}

func (s *FillLimitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitFillLimit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) FillLimit() (localctx IFillLimitContext) {
	localctx = NewFillLimitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, EDCqlMetricParserRULE_fillLimit)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(133)
		p.Match(EDCqlMetricParserNUMBER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFillSectionContext is an interface to support dynamic dispatch.
type IFillSectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FILL_START() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	FillMethod() IFillMethodContext
	RPAREN() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	FillLimit() IFillLimitContext

	// IsFillSectionContext differentiates from other interfaces.
	IsFillSectionContext()
}

type FillSectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillSectionContext() *FillSectionContext {
	var p = new(FillSectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fillSection
	return p
}

func InitEmptyFillSectionContext(p *FillSectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fillSection
}

func (*FillSectionContext) IsFillSectionContext() {}

func NewFillSectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillSectionContext {
	var p = new(FillSectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_fillSection

	return p
}

func (s *FillSectionContext) GetParser() antlr.Parser { return s.parser }

func (s *FillSectionContext) FILL_START() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserFILL_START, 0)
}

func (s *FillSectionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserLPAREN, 0)
}

func (s *FillSectionContext) FillMethod() IFillMethodContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFillMethodContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFillMethodContext)
}

func (s *FillSectionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserRPAREN, 0)
}

func (s *FillSectionContext) COMMA() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOMMA, 0)
}

func (s *FillSectionContext) FillLimit() IFillLimitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFillLimitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFillLimitContext)
}

func (s *FillSectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillSectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FillSectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterFillSection(s)
	}
}

func (s *FillSectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitFillSection(s)
	}
}

func (s *FillSectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitFillSection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) FillSection() (localctx IFillSectionContext) {
	localctx = NewFillSectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, EDCqlMetricParserRULE_fillSection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(135)
		p.Match(EDCqlMetricParserFILL_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(136)
		p.Match(EDCqlMetricParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(137)
		p.FillMethod()
	}
	p.SetState(140)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlMetricParserCOMMA {
		{
			p.SetState(138)
			p.Match(EDCqlMetricParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(139)
			p.FillLimit()
		}

	}
	{
		p.SetState(142)
		p.Match(EDCqlMetricParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDisjQuery() []IDisjQueryContext
	DisjQuery(i int) IDisjQueryContext

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) AllDisjQuery() []IDisjQueryContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDisjQueryContext); ok {
			len++
		}
	}

	tst := make([]IDisjQueryContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDisjQueryContext); ok {
			tst[i] = t.(IDisjQueryContext)
			i++
		}
	}

	return tst
}

func (s *QueryContext) DisjQuery(i int) IDisjQueryContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDisjQueryContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDisjQueryContext)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (s *QueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, EDCqlMetricParserRULE_query)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16769096) != 0) {
		{
			p.SetState(144)
			p.DisjQuery()
		}

		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDisjQueryContext is an interface to support dynamic dispatch.
type IDisjQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllConjQuery() []IConjQueryContext
	ConjQuery(i int) IConjQueryContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsDisjQueryContext differentiates from other interfaces.
	IsDisjQueryContext()
}

type DisjQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDisjQueryContext() *DisjQueryContext {
	var p = new(DisjQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_disjQuery
	return p
}

func InitEmptyDisjQueryContext(p *DisjQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_disjQuery
}

func (*DisjQueryContext) IsDisjQueryContext() {}

func NewDisjQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DisjQueryContext {
	var p = new(DisjQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_disjQuery

	return p
}

func (s *DisjQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *DisjQueryContext) AllConjQuery() []IConjQueryContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConjQueryContext); ok {
			len++
		}
	}

	tst := make([]IConjQueryContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConjQueryContext); ok {
			tst[i] = t.(IConjQueryContext)
			i++
		}
	}

	return tst
}

func (s *DisjQueryContext) ConjQuery(i int) IConjQueryContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConjQueryContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConjQueryContext)
}

func (s *DisjQueryContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserOR)
}

func (s *DisjQueryContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserOR, i)
}

func (s *DisjQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DisjQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DisjQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterDisjQuery(s)
	}
}

func (s *DisjQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitDisjQuery(s)
	}
}

func (s *DisjQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitDisjQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) DisjQuery() (localctx IDisjQueryContext) {
	localctx = NewDisjQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, EDCqlMetricParserRULE_disjQuery)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(149)
		p.ConjQuery()
	}
	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EDCqlMetricParserOR {
		{
			p.SetState(150)
			p.Match(EDCqlMetricParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(151)
			p.ConjQuery()
		}

		p.SetState(156)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConjQueryContext is an interface to support dynamic dispatch.
type IConjQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllModClause() []IModClauseContext
	ModClause(i int) IModClauseContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsConjQueryContext differentiates from other interfaces.
	IsConjQueryContext()
}

type ConjQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConjQueryContext() *ConjQueryContext {
	var p = new(ConjQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_conjQuery
	return p
}

func InitEmptyConjQueryContext(p *ConjQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_conjQuery
}

func (*ConjQueryContext) IsConjQueryContext() {}

func NewConjQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConjQueryContext {
	var p = new(ConjQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_conjQuery

	return p
}

func (s *ConjQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *ConjQueryContext) AllModClause() []IModClauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IModClauseContext); ok {
			len++
		}
	}

	tst := make([]IModClauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IModClauseContext); ok {
			tst[i] = t.(IModClauseContext)
			i++
		}
	}

	return tst
}

func (s *ConjQueryContext) ModClause(i int) IModClauseContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IModClauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IModClauseContext)
}

func (s *ConjQueryContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(EDCqlMetricParserAND)
}

func (s *ConjQueryContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserAND, i)
}

func (s *ConjQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConjQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConjQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterConjQuery(s)
	}
}

func (s *ConjQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitConjQuery(s)
	}
}

func (s *ConjQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitConjQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) ConjQuery() (localctx IConjQueryContext) {
	localctx = NewConjQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, EDCqlMetricParserRULE_conjQuery)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(157)
		p.ModClause()
	}
	p.SetState(164)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(159)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if _la == EDCqlMetricParserAND {
				{
					p.SetState(158)
					p.Match(EDCqlMetricParserAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			}
			{
				p.SetState(161)
				p.ModClause()
			}

		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IModClauseContext is an interface to support dynamic dispatch.
type IModClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Clause() IClauseContext
	Modifier() IModifierContext

	// IsModClauseContext differentiates from other interfaces.
	IsModClauseContext()
}

type ModClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModClauseContext() *ModClauseContext {
	var p = new(ModClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_modClause
	return p
}

func InitEmptyModClauseContext(p *ModClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_modClause
}

func (*ModClauseContext) IsModClauseContext() {}

func NewModClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModClauseContext {
	var p = new(ModClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_modClause

	return p
}

func (s *ModClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ModClauseContext) Clause() IClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClauseContext)
}

func (s *ModClauseContext) Modifier() IModifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IModifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IModifierContext)
}

func (s *ModClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterModClause(s)
	}
}

func (s *ModClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitModClause(s)
	}
}

func (s *ModClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitModClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) ModClause() (localctx IModClauseContext) {
	localctx = NewModClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, EDCqlMetricParserRULE_modClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(168)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlMetricParserNOT {
		{
			p.SetState(167)
			p.Modifier()
		}

	}
	{
		p.SetState(170)
		p.Clause()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IModifierContext is an interface to support dynamic dispatch.
type IModifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NOT() antlr.TerminalNode

	// IsModifierContext differentiates from other interfaces.
	IsModifierContext()
}

type ModifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModifierContext() *ModifierContext {
	var p = new(ModifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_modifier
	return p
}

func InitEmptyModifierContext(p *ModifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_modifier
}

func (*ModifierContext) IsModifierContext() {}

func NewModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModifierContext {
	var p = new(ModifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_modifier

	return p
}

func (s *ModifierContext) GetParser() antlr.Parser { return s.parser }

func (s *ModifierContext) NOT() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserNOT, 0)
}

func (s *ModifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterModifier(s)
	}
}

func (s *ModifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitModifier(s)
	}
}

func (s *ModifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitModifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) Modifier() (localctx IModifierContext) {
	localctx = NewModifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, EDCqlMetricParserRULE_modifier)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(172)
		p.Match(EDCqlMetricParserNOT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IClauseContext is an interface to support dynamic dispatch.
type IClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Term() ITermContext
	GroupingExpr() IGroupingExprContext
	FieldName() IFieldNameContext
	OperatorColon() IOperatorColonContext

	// IsClauseContext differentiates from other interfaces.
	IsClauseContext()
}

type ClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClauseContext() *ClauseContext {
	var p = new(ClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_clause
	return p
}

func InitEmptyClauseContext(p *ClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_clause
}

func (*ClauseContext) IsClauseContext() {}

func NewClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClauseContext {
	var p = new(ClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_clause

	return p
}

func (s *ClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ClauseContext) Term() ITermContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITermContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITermContext)
}

func (s *ClauseContext) GroupingExpr() IGroupingExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupingExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupingExprContext)
}

func (s *ClauseContext) FieldName() IFieldNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldNameContext)
}

func (s *ClauseContext) OperatorColon() IOperatorColonContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperatorColonContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperatorColonContext)
}

func (s *ClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterClause(s)
	}
}

func (s *ClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitClause(s)
	}
}

func (s *ClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) Clause() (localctx IClauseContext) {
	localctx = NewClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, EDCqlMetricParserRULE_clause)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(177)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 17, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(174)
			p.FieldName()
		}

		{
			p.SetState(175)
			p.OperatorColon()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(181)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlMetricParserROLLUP_START, EDCqlMetricParserBY, EDCqlMetricParserCOUNT_UNIQUE, EDCqlMetricParserFILL_START, EDCqlMetricParserCOMMON_OPERATOR, EDCqlMetricParserFILL_OPERATOR, EDCqlMetricParserAGGREGATION_OPERATOR, EDCqlMetricParserCOUNTER_OPERATOR, EDCqlMetricParserQUOTED, EDCqlMetricParserNUMBER, EDCqlMetricParserTERM:
		{
			p.SetState(179)
			p.Term()
		}

	case EDCqlMetricParserLPAREN:
		{
			p.SetState(180)
			p.GroupingExpr()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITermContext is an interface to support dynamic dispatch.
type ITermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QuotedTerm() IQuotedTermContext
	AGGREGATION_OPERATOR() antlr.TerminalNode
	COUNTER_OPERATOR() antlr.TerminalNode
	COMMON_OPERATOR() antlr.TerminalNode
	FILL_OPERATOR() antlr.TerminalNode
	COUNT_UNIQUE() antlr.TerminalNode
	ROLLUP_START() antlr.TerminalNode
	FILL_START() antlr.TerminalNode
	BY() antlr.TerminalNode
	NUMBER() antlr.TerminalNode
	TERM() antlr.TerminalNode

	// IsTermContext differentiates from other interfaces.
	IsTermContext()
}

type TermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTermContext() *TermContext {
	var p = new(TermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_term
	return p
}

func InitEmptyTermContext(p *TermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_term
}

func (*TermContext) IsTermContext() {}

func NewTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TermContext {
	var p = new(TermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_term

	return p
}

func (s *TermContext) GetParser() antlr.Parser { return s.parser }

func (s *TermContext) QuotedTerm() IQuotedTermContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQuotedTermContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQuotedTermContext)
}

func (s *TermContext) AGGREGATION_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserAGGREGATION_OPERATOR, 0)
}

func (s *TermContext) COUNTER_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNTER_OPERATOR, 0)
}

func (s *TermContext) COMMON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOMMON_OPERATOR, 0)
}

func (s *TermContext) FILL_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserFILL_OPERATOR, 0)
}

func (s *TermContext) COUNT_UNIQUE() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNT_UNIQUE, 0)
}

func (s *TermContext) ROLLUP_START() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserROLLUP_START, 0)
}

func (s *TermContext) FILL_START() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserFILL_START, 0)
}

func (s *TermContext) BY() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserBY, 0)
}

func (s *TermContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserNUMBER, 0)
}

func (s *TermContext) TERM() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserTERM, 0)
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterTerm(s)
	}
}

func (s *TermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitTerm(s)
	}
}

func (s *TermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) Term() (localctx ITermContext) {
	localctx = NewTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, EDCqlMetricParserRULE_term)
	p.SetState(194)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlMetricParserQUOTED:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(183)
			p.QuotedTerm()
		}

	case EDCqlMetricParserAGGREGATION_OPERATOR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(184)
			p.Match(EDCqlMetricParserAGGREGATION_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserCOUNTER_OPERATOR:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(185)
			p.Match(EDCqlMetricParserCOUNTER_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserCOMMON_OPERATOR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(186)
			p.Match(EDCqlMetricParserCOMMON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserFILL_OPERATOR:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(187)
			p.Match(EDCqlMetricParserFILL_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserCOUNT_UNIQUE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(188)
			p.Match(EDCqlMetricParserCOUNT_UNIQUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserROLLUP_START:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(189)
			p.Match(EDCqlMetricParserROLLUP_START)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserFILL_START:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(190)
			p.Match(EDCqlMetricParserFILL_START)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserBY:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(191)
			p.Match(EDCqlMetricParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserNUMBER:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(192)
			p.Match(EDCqlMetricParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlMetricParserTERM:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(193)
			p.Match(EDCqlMetricParserTERM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGroupingExprContext is an interface to support dynamic dispatch.
type IGroupingExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	Query() IQueryContext
	RPAREN() antlr.TerminalNode

	// IsGroupingExprContext differentiates from other interfaces.
	IsGroupingExprContext()
}

type GroupingExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupingExprContext() *GroupingExprContext {
	var p = new(GroupingExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_groupingExpr
	return p
}

func InitEmptyGroupingExprContext(p *GroupingExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_groupingExpr
}

func (*GroupingExprContext) IsGroupingExprContext() {}

func NewGroupingExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupingExprContext {
	var p = new(GroupingExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_groupingExpr

	return p
}

func (s *GroupingExprContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupingExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserLPAREN, 0)
}

func (s *GroupingExprContext) Query() IQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryContext)
}

func (s *GroupingExprContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserRPAREN, 0)
}

func (s *GroupingExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupingExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupingExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterGroupingExpr(s)
	}
}

func (s *GroupingExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitGroupingExpr(s)
	}
}

func (s *GroupingExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitGroupingExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) GroupingExpr() (localctx IGroupingExprContext) {
	localctx = NewGroupingExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, EDCqlMetricParserRULE_groupingExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(196)
		p.Match(EDCqlMetricParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(197)
		p.Query()
	}
	{
		p.SetState(198)
		p.Match(EDCqlMetricParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldNameContext is an interface to support dynamic dispatch.
type IFieldNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TERM() antlr.TerminalNode
	COUNTER_OPERATOR() antlr.TerminalNode

	// IsFieldNameContext differentiates from other interfaces.
	IsFieldNameContext()
}

type FieldNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldNameContext() *FieldNameContext {
	var p = new(FieldNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fieldName
	return p
}

func InitEmptyFieldNameContext(p *FieldNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_fieldName
}

func (*FieldNameContext) IsFieldNameContext() {}

func NewFieldNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldNameContext {
	var p = new(FieldNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_fieldName

	return p
}

func (s *FieldNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldNameContext) TERM() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserTERM, 0)
}

func (s *FieldNameContext) COUNTER_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserCOUNTER_OPERATOR, 0)
}

func (s *FieldNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterFieldName(s)
	}
}

func (s *FieldNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitFieldName(s)
	}
}

func (s *FieldNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitFieldName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) FieldName() (localctx IFieldNameContext) {
	localctx = NewFieldNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, EDCqlMetricParserRULE_fieldName)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(200)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EDCqlMetricParserCOUNTER_OPERATOR || _la == EDCqlMetricParserTERM) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQuotedTermContext is an interface to support dynamic dispatch.
type IQuotedTermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QUOTED() antlr.TerminalNode

	// IsQuotedTermContext differentiates from other interfaces.
	IsQuotedTermContext()
}

type QuotedTermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuotedTermContext() *QuotedTermContext {
	var p = new(QuotedTermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_quotedTerm
	return p
}

func InitEmptyQuotedTermContext(p *QuotedTermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlMetricParserRULE_quotedTerm
}

func (*QuotedTermContext) IsQuotedTermContext() {}

func NewQuotedTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QuotedTermContext {
	var p = new(QuotedTermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlMetricParserRULE_quotedTerm

	return p
}

func (s *QuotedTermContext) GetParser() antlr.Parser { return s.parser }

func (s *QuotedTermContext) QUOTED() antlr.TerminalNode {
	return s.GetToken(EDCqlMetricParserQUOTED, 0)
}

func (s *QuotedTermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuotedTermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QuotedTermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.EnterQuotedTerm(s)
	}
}

func (s *QuotedTermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlMetricParserListener); ok {
		listenerT.ExitQuotedTerm(s)
	}
}

func (s *QuotedTermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlMetricParserVisitor:
		return t.VisitQuotedTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlMetricParser) QuotedTerm() (localctx IQuotedTermContext) {
	localctx = NewQuotedTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, EDCqlMetricParserRULE_quotedTerm)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(202)
		p.Match(EDCqlMetricParserQUOTED)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
