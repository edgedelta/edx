// Code generated from EDCqlLogParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package logparser // EDCqlLogParser
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

type EDCqlLogParser struct {
	*antlr.BaseParser
}

var EDCqlLogParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func edcqllogparserParserInit() {
	staticData := &EDCqlLogParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "'{'", "'}'", "'('", "')'", "", "", "'*'", "", "'.rollup'",
		"", "'count_unique'", "'.fill'", "'last'",
	}
	staticData.SymbolicNames = []string{
		"", "AND", "OR", "NOT", "LBRACE", "RBRACE", "LPAREN", "RPAREN", "OP_COLON",
		"OPERATOR", "MATCH_ALL", "COMMA", "ROLLUP_START", "BY", "COUNT_UNIQUE",
		"FILL_START", "COMMON_OPERATOR", "FILL_OPERATOR", "AGGREGATION_OPERATOR",
		"QUOTED", "NUMBER", "TERM", "DEFAULT_SKIP", "UNKNOWN",
	}
	staticData.RuleNames = []string{
		"finalQuery", "topLevelQuery", "rollupMethod", "rollupField", "rollupSection",
		"countUniqueFields", "groupByFields", "groupBySection", "countUniqueSection",
		"queryOperation", "query", "disjQuery", "conjQuery", "modClause", "modifier",
		"clause", "term", "groupingExpr", "fieldName", "quotedTerm", "operatorColon",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 23, 173, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 1,
		0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 49, 8, 1, 1, 2, 1, 2, 1, 3, 1, 3,
		1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 5, 5, 65, 8,
		5, 10, 5, 12, 5, 68, 9, 5, 1, 6, 1, 6, 1, 6, 5, 6, 73, 8, 6, 10, 6, 12,
		6, 76, 9, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 3, 8, 86,
		8, 8, 1, 8, 3, 8, 89, 8, 8, 1, 9, 1, 9, 1, 9, 3, 9, 94, 8, 9, 1, 9, 1,
		9, 1, 9, 3, 9, 99, 8, 9, 1, 9, 1, 9, 1, 9, 3, 9, 104, 8, 9, 1, 9, 3, 9,
		107, 8, 9, 1, 9, 3, 9, 110, 8, 9, 1, 10, 4, 10, 113, 8, 10, 11, 10, 12,
		10, 114, 1, 11, 1, 11, 1, 11, 5, 11, 120, 8, 11, 10, 11, 12, 11, 123, 9,
		11, 1, 12, 1, 12, 3, 12, 127, 8, 12, 1, 12, 5, 12, 130, 8, 12, 10, 12,
		12, 12, 133, 9, 12, 1, 13, 3, 13, 136, 8, 13, 1, 13, 1, 13, 1, 14, 1, 14,
		1, 15, 1, 15, 1, 15, 3, 15, 145, 8, 15, 1, 15, 1, 15, 3, 15, 149, 8, 15,
		1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 3,
		16, 161, 8, 16, 1, 17, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 19, 1, 19,
		1, 20, 1, 20, 1, 20, 0, 0, 21, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22,
		24, 26, 28, 30, 32, 34, 36, 38, 40, 0, 2, 2, 0, 16, 16, 18, 18, 1, 0, 8,
		9, 178, 0, 42, 1, 0, 0, 0, 2, 48, 1, 0, 0, 0, 4, 50, 1, 0, 0, 0, 6, 52,
		1, 0, 0, 0, 8, 54, 1, 0, 0, 0, 10, 61, 1, 0, 0, 0, 12, 69, 1, 0, 0, 0,
		14, 77, 1, 0, 0, 0, 16, 82, 1, 0, 0, 0, 18, 93, 1, 0, 0, 0, 20, 112, 1,
		0, 0, 0, 22, 116, 1, 0, 0, 0, 24, 124, 1, 0, 0, 0, 26, 135, 1, 0, 0, 0,
		28, 139, 1, 0, 0, 0, 30, 144, 1, 0, 0, 0, 32, 160, 1, 0, 0, 0, 34, 162,
		1, 0, 0, 0, 36, 166, 1, 0, 0, 0, 38, 168, 1, 0, 0, 0, 40, 170, 1, 0, 0,
		0, 42, 43, 3, 2, 1, 0, 43, 1, 1, 0, 0, 0, 44, 45, 3, 18, 9, 0, 45, 46,
		5, 0, 0, 1, 46, 49, 1, 0, 0, 0, 47, 49, 5, 0, 0, 1, 48, 44, 1, 0, 0, 0,
		48, 47, 1, 0, 0, 0, 49, 3, 1, 0, 0, 0, 50, 51, 7, 0, 0, 0, 51, 5, 1, 0,
		0, 0, 52, 53, 5, 21, 0, 0, 53, 7, 1, 0, 0, 0, 54, 55, 5, 12, 0, 0, 55,
		56, 5, 6, 0, 0, 56, 57, 3, 4, 2, 0, 57, 58, 5, 11, 0, 0, 58, 59, 3, 6,
		3, 0, 59, 60, 5, 7, 0, 0, 60, 9, 1, 0, 0, 0, 61, 66, 5, 21, 0, 0, 62, 63,
		5, 11, 0, 0, 63, 65, 5, 21, 0, 0, 64, 62, 1, 0, 0, 0, 65, 68, 1, 0, 0,
		0, 66, 64, 1, 0, 0, 0, 66, 67, 1, 0, 0, 0, 67, 11, 1, 0, 0, 0, 68, 66,
		1, 0, 0, 0, 69, 74, 5, 21, 0, 0, 70, 71, 5, 11, 0, 0, 71, 73, 5, 21, 0,
		0, 72, 70, 1, 0, 0, 0, 73, 76, 1, 0, 0, 0, 74, 72, 1, 0, 0, 0, 74, 75,
		1, 0, 0, 0, 75, 13, 1, 0, 0, 0, 76, 74, 1, 0, 0, 0, 77, 78, 5, 13, 0, 0,
		78, 79, 5, 4, 0, 0, 79, 80, 3, 12, 6, 0, 80, 81, 5, 5, 0, 0, 81, 15, 1,
		0, 0, 0, 82, 88, 5, 14, 0, 0, 83, 85, 5, 6, 0, 0, 84, 86, 3, 10, 5, 0,
		85, 84, 1, 0, 0, 0, 85, 86, 1, 0, 0, 0, 86, 87, 1, 0, 0, 0, 87, 89, 5,
		7, 0, 0, 88, 83, 1, 0, 0, 0, 88, 89, 1, 0, 0, 0, 89, 17, 1, 0, 0, 0, 90,
		91, 3, 16, 8, 0, 91, 92, 5, 8, 0, 0, 92, 94, 1, 0, 0, 0, 93, 90, 1, 0,
		0, 0, 93, 94, 1, 0, 0, 0, 94, 103, 1, 0, 0, 0, 95, 98, 5, 4, 0, 0, 96,
		99, 3, 20, 10, 0, 97, 99, 5, 10, 0, 0, 98, 96, 1, 0, 0, 0, 98, 97, 1, 0,
		0, 0, 99, 100, 1, 0, 0, 0, 100, 104, 5, 5, 0, 0, 101, 104, 3, 20, 10, 0,
		102, 104, 5, 10, 0, 0, 103, 95, 1, 0, 0, 0, 103, 101, 1, 0, 0, 0, 103,
		102, 1, 0, 0, 0, 104, 106, 1, 0, 0, 0, 105, 107, 3, 14, 7, 0, 106, 105,
		1, 0, 0, 0, 106, 107, 1, 0, 0, 0, 107, 109, 1, 0, 0, 0, 108, 110, 3, 8,
		4, 0, 109, 108, 1, 0, 0, 0, 109, 110, 1, 0, 0, 0, 110, 19, 1, 0, 0, 0,
		111, 113, 3, 22, 11, 0, 112, 111, 1, 0, 0, 0, 113, 114, 1, 0, 0, 0, 114,
		112, 1, 0, 0, 0, 114, 115, 1, 0, 0, 0, 115, 21, 1, 0, 0, 0, 116, 121, 3,
		24, 12, 0, 117, 118, 5, 2, 0, 0, 118, 120, 3, 24, 12, 0, 119, 117, 1, 0,
		0, 0, 120, 123, 1, 0, 0, 0, 121, 119, 1, 0, 0, 0, 121, 122, 1, 0, 0, 0,
		122, 23, 1, 0, 0, 0, 123, 121, 1, 0, 0, 0, 124, 131, 3, 26, 13, 0, 125,
		127, 5, 1, 0, 0, 126, 125, 1, 0, 0, 0, 126, 127, 1, 0, 0, 0, 127, 128,
		1, 0, 0, 0, 128, 130, 3, 26, 13, 0, 129, 126, 1, 0, 0, 0, 130, 133, 1,
		0, 0, 0, 131, 129, 1, 0, 0, 0, 131, 132, 1, 0, 0, 0, 132, 25, 1, 0, 0,
		0, 133, 131, 1, 0, 0, 0, 134, 136, 3, 28, 14, 0, 135, 134, 1, 0, 0, 0,
		135, 136, 1, 0, 0, 0, 136, 137, 1, 0, 0, 0, 137, 138, 3, 30, 15, 0, 138,
		27, 1, 0, 0, 0, 139, 140, 5, 3, 0, 0, 140, 29, 1, 0, 0, 0, 141, 142, 3,
		36, 18, 0, 142, 143, 3, 40, 20, 0, 143, 145, 1, 0, 0, 0, 144, 141, 1, 0,
		0, 0, 144, 145, 1, 0, 0, 0, 145, 148, 1, 0, 0, 0, 146, 149, 3, 32, 16,
		0, 147, 149, 3, 34, 17, 0, 148, 146, 1, 0, 0, 0, 148, 147, 1, 0, 0, 0,
		149, 31, 1, 0, 0, 0, 150, 161, 3, 38, 19, 0, 151, 161, 5, 18, 0, 0, 152,
		161, 5, 16, 0, 0, 153, 161, 5, 17, 0, 0, 154, 161, 5, 14, 0, 0, 155, 161,
		5, 12, 0, 0, 156, 161, 5, 15, 0, 0, 157, 161, 5, 13, 0, 0, 158, 161, 5,
		20, 0, 0, 159, 161, 5, 21, 0, 0, 160, 150, 1, 0, 0, 0, 160, 151, 1, 0,
		0, 0, 160, 152, 1, 0, 0, 0, 160, 153, 1, 0, 0, 0, 160, 154, 1, 0, 0, 0,
		160, 155, 1, 0, 0, 0, 160, 156, 1, 0, 0, 0, 160, 157, 1, 0, 0, 0, 160,
		158, 1, 0, 0, 0, 160, 159, 1, 0, 0, 0, 161, 33, 1, 0, 0, 0, 162, 163, 5,
		6, 0, 0, 163, 164, 3, 20, 10, 0, 164, 165, 5, 7, 0, 0, 165, 35, 1, 0, 0,
		0, 166, 167, 5, 21, 0, 0, 167, 37, 1, 0, 0, 0, 168, 169, 5, 19, 0, 0, 169,
		39, 1, 0, 0, 0, 170, 171, 7, 1, 0, 0, 171, 41, 1, 0, 0, 0, 18, 48, 66,
		74, 85, 88, 93, 98, 103, 106, 109, 114, 121, 126, 131, 135, 144, 148, 160,
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

// EDCqlLogParserInit initializes any static state used to implement EDCqlLogParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewEDCqlLogParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func EDCqlLogParserInit() {
	staticData := &EDCqlLogParserParserStaticData
	staticData.once.Do(edcqllogparserParserInit)
}

// NewEDCqlLogParser produces a new parser instance for the optional input antlr.TokenStream.
func NewEDCqlLogParser(input antlr.TokenStream) *EDCqlLogParser {
	EDCqlLogParserInit()
	this := new(EDCqlLogParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &EDCqlLogParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "EDCqlLogParser.g4"

	return this
}

// EDCqlLogParser tokens.
const (
	EDCqlLogParserEOF                  = antlr.TokenEOF
	EDCqlLogParserAND                  = 1
	EDCqlLogParserOR                   = 2
	EDCqlLogParserNOT                  = 3
	EDCqlLogParserLBRACE               = 4
	EDCqlLogParserRBRACE               = 5
	EDCqlLogParserLPAREN               = 6
	EDCqlLogParserRPAREN               = 7
	EDCqlLogParserOP_COLON             = 8
	EDCqlLogParserOPERATOR             = 9
	EDCqlLogParserMATCH_ALL            = 10
	EDCqlLogParserCOMMA                = 11
	EDCqlLogParserROLLUP_START         = 12
	EDCqlLogParserBY                   = 13
	EDCqlLogParserCOUNT_UNIQUE         = 14
	EDCqlLogParserFILL_START           = 15
	EDCqlLogParserCOMMON_OPERATOR      = 16
	EDCqlLogParserFILL_OPERATOR        = 17
	EDCqlLogParserAGGREGATION_OPERATOR = 18
	EDCqlLogParserQUOTED               = 19
	EDCqlLogParserNUMBER               = 20
	EDCqlLogParserTERM                 = 21
	EDCqlLogParserDEFAULT_SKIP         = 22
	EDCqlLogParserUNKNOWN              = 23
)

// EDCqlLogParser rules.
const (
	EDCqlLogParserRULE_finalQuery         = 0
	EDCqlLogParserRULE_topLevelQuery      = 1
	EDCqlLogParserRULE_rollupMethod       = 2
	EDCqlLogParserRULE_rollupField        = 3
	EDCqlLogParserRULE_rollupSection      = 4
	EDCqlLogParserRULE_countUniqueFields  = 5
	EDCqlLogParserRULE_groupByFields      = 6
	EDCqlLogParserRULE_groupBySection     = 7
	EDCqlLogParserRULE_countUniqueSection = 8
	EDCqlLogParserRULE_queryOperation     = 9
	EDCqlLogParserRULE_query              = 10
	EDCqlLogParserRULE_disjQuery          = 11
	EDCqlLogParserRULE_conjQuery          = 12
	EDCqlLogParserRULE_modClause          = 13
	EDCqlLogParserRULE_modifier           = 14
	EDCqlLogParserRULE_clause             = 15
	EDCqlLogParserRULE_term               = 16
	EDCqlLogParserRULE_groupingExpr       = 17
	EDCqlLogParserRULE_fieldName          = 18
	EDCqlLogParserRULE_quotedTerm         = 19
	EDCqlLogParserRULE_operatorColon      = 20
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
	p.RuleIndex = EDCqlLogParserRULE_finalQuery
	return p
}

func InitEmptyFinalQueryContext(p *FinalQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_finalQuery
}

func (*FinalQueryContext) IsFinalQueryContext() {}

func NewFinalQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FinalQueryContext {
	var p = new(FinalQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_finalQuery

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
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterFinalQuery(s)
	}
}

func (s *FinalQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitFinalQuery(s)
	}
}

func (s *FinalQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitFinalQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) FinalQuery() (localctx IFinalQueryContext) {
	localctx = NewFinalQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, EDCqlLogParserRULE_finalQuery)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(42)
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
	p.RuleIndex = EDCqlLogParserRULE_topLevelQuery
	return p
}

func InitEmptyTopLevelQueryContext(p *TopLevelQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_topLevelQuery
}

func (*TopLevelQueryContext) IsTopLevelQueryContext() {}

func NewTopLevelQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelQueryContext {
	var p = new(TopLevelQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_topLevelQuery

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
	return s.GetToken(EDCqlLogParserEOF, 0)
}

func (s *TopLevelQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TopLevelQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TopLevelQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterTopLevelQuery(s)
	}
}

func (s *TopLevelQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitTopLevelQuery(s)
	}
}

func (s *TopLevelQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitTopLevelQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) TopLevelQuery() (localctx ITopLevelQueryContext) {
	localctx = NewTopLevelQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, EDCqlLogParserRULE_topLevelQuery)
	p.SetState(48)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlLogParserNOT, EDCqlLogParserLBRACE, EDCqlLogParserLPAREN, EDCqlLogParserMATCH_ALL, EDCqlLogParserROLLUP_START, EDCqlLogParserBY, EDCqlLogParserCOUNT_UNIQUE, EDCqlLogParserFILL_START, EDCqlLogParserCOMMON_OPERATOR, EDCqlLogParserFILL_OPERATOR, EDCqlLogParserAGGREGATION_OPERATOR, EDCqlLogParserQUOTED, EDCqlLogParserNUMBER, EDCqlLogParserTERM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(44)
			p.QueryOperation()
		}
		{
			p.SetState(45)
			p.Match(EDCqlLogParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserEOF:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(47)
			p.Match(EDCqlLogParserEOF)
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

// IRollupMethodContext is an interface to support dynamic dispatch.
type IRollupMethodContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AGGREGATION_OPERATOR() antlr.TerminalNode
	COMMON_OPERATOR() antlr.TerminalNode

	// IsRollupMethodContext differentiates from other interfaces.
	IsRollupMethodContext()
}

type RollupMethodContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRollupMethodContext() *RollupMethodContext {
	var p = new(RollupMethodContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_rollupMethod
	return p
}

func InitEmptyRollupMethodContext(p *RollupMethodContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_rollupMethod
}

func (*RollupMethodContext) IsRollupMethodContext() {}

func NewRollupMethodContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RollupMethodContext {
	var p = new(RollupMethodContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_rollupMethod

	return p
}

func (s *RollupMethodContext) GetParser() antlr.Parser { return s.parser }

func (s *RollupMethodContext) AGGREGATION_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserAGGREGATION_OPERATOR, 0)
}

func (s *RollupMethodContext) COMMON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserCOMMON_OPERATOR, 0)
}

func (s *RollupMethodContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RollupMethodContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RollupMethodContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterRollupMethod(s)
	}
}

func (s *RollupMethodContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitRollupMethod(s)
	}
}

func (s *RollupMethodContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitRollupMethod(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) RollupMethod() (localctx IRollupMethodContext) {
	localctx = NewRollupMethodContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, EDCqlLogParserRULE_rollupMethod)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(50)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EDCqlLogParserCOMMON_OPERATOR || _la == EDCqlLogParserAGGREGATION_OPERATOR) {
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

// IRollupFieldContext is an interface to support dynamic dispatch.
type IRollupFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TERM() antlr.TerminalNode

	// IsRollupFieldContext differentiates from other interfaces.
	IsRollupFieldContext()
}

type RollupFieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRollupFieldContext() *RollupFieldContext {
	var p = new(RollupFieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_rollupField
	return p
}

func InitEmptyRollupFieldContext(p *RollupFieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_rollupField
}

func (*RollupFieldContext) IsRollupFieldContext() {}

func NewRollupFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RollupFieldContext {
	var p = new(RollupFieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_rollupField

	return p
}

func (s *RollupFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *RollupFieldContext) TERM() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserTERM, 0)
}

func (s *RollupFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RollupFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RollupFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterRollupField(s)
	}
}

func (s *RollupFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitRollupField(s)
	}
}

func (s *RollupFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitRollupField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) RollupField() (localctx IRollupFieldContext) {
	localctx = NewRollupFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, EDCqlLogParserRULE_rollupField)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(52)
		p.Match(EDCqlLogParserTERM)
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
	RollupMethod() IRollupMethodContext
	COMMA() antlr.TerminalNode
	RollupField() IRollupFieldContext
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
	p.RuleIndex = EDCqlLogParserRULE_rollupSection
	return p
}

func InitEmptyRollupSectionContext(p *RollupSectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_rollupSection
}

func (*RollupSectionContext) IsRollupSectionContext() {}

func NewRollupSectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RollupSectionContext {
	var p = new(RollupSectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_rollupSection

	return p
}

func (s *RollupSectionContext) GetParser() antlr.Parser { return s.parser }

func (s *RollupSectionContext) ROLLUP_START() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserROLLUP_START, 0)
}

func (s *RollupSectionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserLPAREN, 0)
}

func (s *RollupSectionContext) RollupMethod() IRollupMethodContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRollupMethodContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRollupMethodContext)
}

func (s *RollupSectionContext) COMMA() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserCOMMA, 0)
}

func (s *RollupSectionContext) RollupField() IRollupFieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRollupFieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRollupFieldContext)
}

func (s *RollupSectionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserRPAREN, 0)
}

func (s *RollupSectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RollupSectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RollupSectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterRollupSection(s)
	}
}

func (s *RollupSectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitRollupSection(s)
	}
}

func (s *RollupSectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitRollupSection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) RollupSection() (localctx IRollupSectionContext) {
	localctx = NewRollupSectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, EDCqlLogParserRULE_rollupSection)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(54)
		p.Match(EDCqlLogParserROLLUP_START)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(55)
		p.Match(EDCqlLogParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(56)
		p.RollupMethod()
	}
	{
		p.SetState(57)
		p.Match(EDCqlLogParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(58)
		p.RollupField()
	}
	{
		p.SetState(59)
		p.Match(EDCqlLogParserRPAREN)
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

// ICountUniqueFieldsContext is an interface to support dynamic dispatch.
type ICountUniqueFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTERM() []antlr.TerminalNode
	TERM(i int) antlr.TerminalNode
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
	p.RuleIndex = EDCqlLogParserRULE_countUniqueFields
	return p
}

func InitEmptyCountUniqueFieldsContext(p *CountUniqueFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_countUniqueFields
}

func (*CountUniqueFieldsContext) IsCountUniqueFieldsContext() {}

func NewCountUniqueFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CountUniqueFieldsContext {
	var p = new(CountUniqueFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_countUniqueFields

	return p
}

func (s *CountUniqueFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *CountUniqueFieldsContext) AllTERM() []antlr.TerminalNode {
	return s.GetTokens(EDCqlLogParserTERM)
}

func (s *CountUniqueFieldsContext) TERM(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserTERM, i)
}

func (s *CountUniqueFieldsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EDCqlLogParserCOMMA)
}

func (s *CountUniqueFieldsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserCOMMA, i)
}

func (s *CountUniqueFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CountUniqueFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CountUniqueFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterCountUniqueFields(s)
	}
}

func (s *CountUniqueFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitCountUniqueFields(s)
	}
}

func (s *CountUniqueFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitCountUniqueFields(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) CountUniqueFields() (localctx ICountUniqueFieldsContext) {
	localctx = NewCountUniqueFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, EDCqlLogParserRULE_countUniqueFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(61)
		p.Match(EDCqlLogParserTERM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(66)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EDCqlLogParserCOMMA {
		{
			p.SetState(62)
			p.Match(EDCqlLogParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(63)
			p.Match(EDCqlLogParserTERM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(68)
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

// IGroupByFieldsContext is an interface to support dynamic dispatch.
type IGroupByFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTERM() []antlr.TerminalNode
	TERM(i int) antlr.TerminalNode
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
	p.RuleIndex = EDCqlLogParserRULE_groupByFields
	return p
}

func InitEmptyGroupByFieldsContext(p *GroupByFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_groupByFields
}

func (*GroupByFieldsContext) IsGroupByFieldsContext() {}

func NewGroupByFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByFieldsContext {
	var p = new(GroupByFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_groupByFields

	return p
}

func (s *GroupByFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByFieldsContext) AllTERM() []antlr.TerminalNode {
	return s.GetTokens(EDCqlLogParserTERM)
}

func (s *GroupByFieldsContext) TERM(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserTERM, i)
}

func (s *GroupByFieldsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EDCqlLogParserCOMMA)
}

func (s *GroupByFieldsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserCOMMA, i)
}

func (s *GroupByFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupByFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterGroupByFields(s)
	}
}

func (s *GroupByFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitGroupByFields(s)
	}
}

func (s *GroupByFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitGroupByFields(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) GroupByFields() (localctx IGroupByFieldsContext) {
	localctx = NewGroupByFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, EDCqlLogParserRULE_groupByFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(EDCqlLogParserTERM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(74)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EDCqlLogParserCOMMA {
		{
			p.SetState(70)
			p.Match(EDCqlLogParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(71)
			p.Match(EDCqlLogParserTERM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(76)
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
	p.RuleIndex = EDCqlLogParserRULE_groupBySection
	return p
}

func InitEmptyGroupBySectionContext(p *GroupBySectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_groupBySection
}

func (*GroupBySectionContext) IsGroupBySectionContext() {}

func NewGroupBySectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupBySectionContext {
	var p = new(GroupBySectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_groupBySection

	return p
}

func (s *GroupBySectionContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupBySectionContext) BY() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserBY, 0)
}

func (s *GroupBySectionContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserLBRACE, 0)
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
	return s.GetToken(EDCqlLogParserRBRACE, 0)
}

func (s *GroupBySectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupBySectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupBySectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterGroupBySection(s)
	}
}

func (s *GroupBySectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitGroupBySection(s)
	}
}

func (s *GroupBySectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitGroupBySection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) GroupBySection() (localctx IGroupBySectionContext) {
	localctx = NewGroupBySectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, EDCqlLogParserRULE_groupBySection)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(77)
		p.Match(EDCqlLogParserBY)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(78)
		p.Match(EDCqlLogParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(79)
		p.GroupByFields()
	}
	{
		p.SetState(80)
		p.Match(EDCqlLogParserRBRACE)
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
	p.RuleIndex = EDCqlLogParserRULE_countUniqueSection
	return p
}

func InitEmptyCountUniqueSectionContext(p *CountUniqueSectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_countUniqueSection
}

func (*CountUniqueSectionContext) IsCountUniqueSectionContext() {}

func NewCountUniqueSectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CountUniqueSectionContext {
	var p = new(CountUniqueSectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_countUniqueSection

	return p
}

func (s *CountUniqueSectionContext) GetParser() antlr.Parser { return s.parser }

func (s *CountUniqueSectionContext) COUNT_UNIQUE() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserCOUNT_UNIQUE, 0)
}

func (s *CountUniqueSectionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserLPAREN, 0)
}

func (s *CountUniqueSectionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserRPAREN, 0)
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
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterCountUniqueSection(s)
	}
}

func (s *CountUniqueSectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitCountUniqueSection(s)
	}
}

func (s *CountUniqueSectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitCountUniqueSection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) CountUniqueSection() (localctx ICountUniqueSectionContext) {
	localctx = NewCountUniqueSectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, EDCqlLogParserRULE_countUniqueSection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(82)
		p.Match(EDCqlLogParserCOUNT_UNIQUE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(88)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlLogParserLPAREN {
		{
			p.SetState(83)
			p.Match(EDCqlLogParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(85)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == EDCqlLogParserTERM {
			{
				p.SetState(84)
				p.CountUniqueFields()
			}

		}
		{
			p.SetState(87)
			p.Match(EDCqlLogParserRPAREN)
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

// IQueryOperationContext is an interface to support dynamic dispatch.
type IQueryOperationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	Query() IQueryContext
	MATCH_ALL() antlr.TerminalNode
	CountUniqueSection() ICountUniqueSectionContext
	OP_COLON() antlr.TerminalNode
	GroupBySection() IGroupBySectionContext
	RollupSection() IRollupSectionContext

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
	p.RuleIndex = EDCqlLogParserRULE_queryOperation
	return p
}

func InitEmptyQueryOperationContext(p *QueryOperationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_queryOperation
}

func (*QueryOperationContext) IsQueryOperationContext() {}

func NewQueryOperationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryOperationContext {
	var p = new(QueryOperationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_queryOperation

	return p
}

func (s *QueryOperationContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryOperationContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserLBRACE, 0)
}

func (s *QueryOperationContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserRBRACE, 0)
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

func (s *QueryOperationContext) MATCH_ALL() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserMATCH_ALL, 0)
}

func (s *QueryOperationContext) CountUniqueSection() ICountUniqueSectionContext {
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

func (s *QueryOperationContext) OP_COLON() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserOP_COLON, 0)
}

func (s *QueryOperationContext) GroupBySection() IGroupBySectionContext {
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

func (s *QueryOperationContext) RollupSection() IRollupSectionContext {
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

func (s *QueryOperationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryOperationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryOperationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterQueryOperation(s)
	}
}

func (s *QueryOperationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitQueryOperation(s)
	}
}

func (s *QueryOperationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitQueryOperation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) QueryOperation() (localctx IQueryOperationContext) {
	localctx = NewQueryOperationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, EDCqlLogParserRULE_queryOperation)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(93)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(90)
			p.CountUniqueSection()
		}
		{
			p.SetState(91)
			p.Match(EDCqlLogParserOP_COLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(103)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlLogParserLBRACE:
		{
			p.SetState(95)
			p.Match(EDCqlLogParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(98)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case EDCqlLogParserNOT, EDCqlLogParserLPAREN, EDCqlLogParserROLLUP_START, EDCqlLogParserBY, EDCqlLogParserCOUNT_UNIQUE, EDCqlLogParserFILL_START, EDCqlLogParserCOMMON_OPERATOR, EDCqlLogParserFILL_OPERATOR, EDCqlLogParserAGGREGATION_OPERATOR, EDCqlLogParserQUOTED, EDCqlLogParserNUMBER, EDCqlLogParserTERM:
			{
				p.SetState(96)
				p.Query()
			}

		case EDCqlLogParserMATCH_ALL:
			{
				p.SetState(97)
				p.Match(EDCqlLogParserMATCH_ALL)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(100)
			p.Match(EDCqlLogParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserNOT, EDCqlLogParserLPAREN, EDCqlLogParserROLLUP_START, EDCqlLogParserBY, EDCqlLogParserCOUNT_UNIQUE, EDCqlLogParserFILL_START, EDCqlLogParserCOMMON_OPERATOR, EDCqlLogParserFILL_OPERATOR, EDCqlLogParserAGGREGATION_OPERATOR, EDCqlLogParserQUOTED, EDCqlLogParserNUMBER, EDCqlLogParserTERM:
		{
			p.SetState(101)
			p.Query()
		}

	case EDCqlLogParserMATCH_ALL:
		{
			p.SetState(102)
			p.Match(EDCqlLogParserMATCH_ALL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.SetState(106)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlLogParserBY {
		{
			p.SetState(105)
			p.GroupBySection()
		}

	}
	p.SetState(109)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlLogParserROLLUP_START {
		{
			p.SetState(108)
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
	p.RuleIndex = EDCqlLogParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_query

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
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (s *QueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, EDCqlLogParserRULE_query)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(111)
				p.DisjQuery()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(114)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext())
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
	p.RuleIndex = EDCqlLogParserRULE_disjQuery
	return p
}

func InitEmptyDisjQueryContext(p *DisjQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_disjQuery
}

func (*DisjQueryContext) IsDisjQueryContext() {}

func NewDisjQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DisjQueryContext {
	var p = new(DisjQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_disjQuery

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
	return s.GetTokens(EDCqlLogParserOR)
}

func (s *DisjQueryContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserOR, i)
}

func (s *DisjQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DisjQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DisjQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterDisjQuery(s)
	}
}

func (s *DisjQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitDisjQuery(s)
	}
}

func (s *DisjQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitDisjQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) DisjQuery() (localctx IDisjQueryContext) {
	localctx = NewDisjQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, EDCqlLogParserRULE_disjQuery)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(116)
		p.ConjQuery()
	}
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EDCqlLogParserOR {
		{
			p.SetState(117)
			p.Match(EDCqlLogParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(118)
			p.ConjQuery()
		}

		p.SetState(123)
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
	p.RuleIndex = EDCqlLogParserRULE_conjQuery
	return p
}

func InitEmptyConjQueryContext(p *ConjQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_conjQuery
}

func (*ConjQueryContext) IsConjQueryContext() {}

func NewConjQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConjQueryContext {
	var p = new(ConjQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_conjQuery

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
	return s.GetTokens(EDCqlLogParserAND)
}

func (s *ConjQueryContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserAND, i)
}

func (s *ConjQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConjQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConjQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterConjQuery(s)
	}
}

func (s *ConjQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitConjQuery(s)
	}
}

func (s *ConjQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitConjQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) ConjQuery() (localctx IConjQueryContext) {
	localctx = NewConjQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, EDCqlLogParserRULE_conjQuery)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(124)
		p.ModClause()
	}
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(126)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if _la == EDCqlLogParserAND {
				{
					p.SetState(125)
					p.Match(EDCqlLogParserAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			}
			{
				p.SetState(128)
				p.ModClause()
			}

		}
		p.SetState(133)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext())
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
	p.RuleIndex = EDCqlLogParserRULE_modClause
	return p
}

func InitEmptyModClauseContext(p *ModClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_modClause
}

func (*ModClauseContext) IsModClauseContext() {}

func NewModClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModClauseContext {
	var p = new(ModClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_modClause

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
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterModClause(s)
	}
}

func (s *ModClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitModClause(s)
	}
}

func (s *ModClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitModClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) ModClause() (localctx IModClauseContext) {
	localctx = NewModClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, EDCqlLogParserRULE_modClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(135)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == EDCqlLogParserNOT {
		{
			p.SetState(134)
			p.Modifier()
		}

	}
	{
		p.SetState(137)
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
	p.RuleIndex = EDCqlLogParserRULE_modifier
	return p
}

func InitEmptyModifierContext(p *ModifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_modifier
}

func (*ModifierContext) IsModifierContext() {}

func NewModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModifierContext {
	var p = new(ModifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_modifier

	return p
}

func (s *ModifierContext) GetParser() antlr.Parser { return s.parser }

func (s *ModifierContext) NOT() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserNOT, 0)
}

func (s *ModifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterModifier(s)
	}
}

func (s *ModifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitModifier(s)
	}
}

func (s *ModifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitModifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) Modifier() (localctx IModifierContext) {
	localctx = NewModifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, EDCqlLogParserRULE_modifier)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(139)
		p.Match(EDCqlLogParserNOT)
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
	p.RuleIndex = EDCqlLogParserRULE_clause
	return p
}

func InitEmptyClauseContext(p *ClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_clause
}

func (*ClauseContext) IsClauseContext() {}

func NewClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClauseContext {
	var p = new(ClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_clause

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
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterClause(s)
	}
}

func (s *ClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitClause(s)
	}
}

func (s *ClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) Clause() (localctx IClauseContext) {
	localctx = NewClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, EDCqlLogParserRULE_clause)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(144)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(141)
			p.FieldName()
		}

		{
			p.SetState(142)
			p.OperatorColon()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(148)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlLogParserROLLUP_START, EDCqlLogParserBY, EDCqlLogParserCOUNT_UNIQUE, EDCqlLogParserFILL_START, EDCqlLogParserCOMMON_OPERATOR, EDCqlLogParserFILL_OPERATOR, EDCqlLogParserAGGREGATION_OPERATOR, EDCqlLogParserQUOTED, EDCqlLogParserNUMBER, EDCqlLogParserTERM:
		{
			p.SetState(146)
			p.Term()
		}

	case EDCqlLogParserLPAREN:
		{
			p.SetState(147)
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
	p.RuleIndex = EDCqlLogParserRULE_term
	return p
}

func InitEmptyTermContext(p *TermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_term
}

func (*TermContext) IsTermContext() {}

func NewTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TermContext {
	var p = new(TermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_term

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
	return s.GetToken(EDCqlLogParserAGGREGATION_OPERATOR, 0)
}

func (s *TermContext) COMMON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserCOMMON_OPERATOR, 0)
}

func (s *TermContext) FILL_OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserFILL_OPERATOR, 0)
}

func (s *TermContext) COUNT_UNIQUE() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserCOUNT_UNIQUE, 0)
}

func (s *TermContext) ROLLUP_START() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserROLLUP_START, 0)
}

func (s *TermContext) FILL_START() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserFILL_START, 0)
}

func (s *TermContext) BY() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserBY, 0)
}

func (s *TermContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserNUMBER, 0)
}

func (s *TermContext) TERM() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserTERM, 0)
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterTerm(s)
	}
}

func (s *TermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitTerm(s)
	}
}

func (s *TermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) Term() (localctx ITermContext) {
	localctx = NewTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, EDCqlLogParserRULE_term)
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlLogParserQUOTED:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(150)
			p.QuotedTerm()
		}

	case EDCqlLogParserAGGREGATION_OPERATOR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(151)
			p.Match(EDCqlLogParserAGGREGATION_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserCOMMON_OPERATOR:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(152)
			p.Match(EDCqlLogParserCOMMON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserFILL_OPERATOR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(153)
			p.Match(EDCqlLogParserFILL_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserCOUNT_UNIQUE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(154)
			p.Match(EDCqlLogParserCOUNT_UNIQUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserROLLUP_START:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(155)
			p.Match(EDCqlLogParserROLLUP_START)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserFILL_START:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(156)
			p.Match(EDCqlLogParserFILL_START)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserBY:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(157)
			p.Match(EDCqlLogParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserNUMBER:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(158)
			p.Match(EDCqlLogParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlLogParserTERM:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(159)
			p.Match(EDCqlLogParserTERM)
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
	p.RuleIndex = EDCqlLogParserRULE_groupingExpr
	return p
}

func InitEmptyGroupingExprContext(p *GroupingExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_groupingExpr
}

func (*GroupingExprContext) IsGroupingExprContext() {}

func NewGroupingExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupingExprContext {
	var p = new(GroupingExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_groupingExpr

	return p
}

func (s *GroupingExprContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupingExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserLPAREN, 0)
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
	return s.GetToken(EDCqlLogParserRPAREN, 0)
}

func (s *GroupingExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupingExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupingExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterGroupingExpr(s)
	}
}

func (s *GroupingExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitGroupingExpr(s)
	}
}

func (s *GroupingExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitGroupingExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) GroupingExpr() (localctx IGroupingExprContext) {
	localctx = NewGroupingExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, EDCqlLogParserRULE_groupingExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(162)
		p.Match(EDCqlLogParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(163)
		p.Query()
	}
	{
		p.SetState(164)
		p.Match(EDCqlLogParserRPAREN)
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
	p.RuleIndex = EDCqlLogParserRULE_fieldName
	return p
}

func InitEmptyFieldNameContext(p *FieldNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_fieldName
}

func (*FieldNameContext) IsFieldNameContext() {}

func NewFieldNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldNameContext {
	var p = new(FieldNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_fieldName

	return p
}

func (s *FieldNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldNameContext) TERM() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserTERM, 0)
}

func (s *FieldNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterFieldName(s)
	}
}

func (s *FieldNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitFieldName(s)
	}
}

func (s *FieldNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitFieldName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) FieldName() (localctx IFieldNameContext) {
	localctx = NewFieldNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, EDCqlLogParserRULE_fieldName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(166)
		p.Match(EDCqlLogParserTERM)
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
	p.RuleIndex = EDCqlLogParserRULE_quotedTerm
	return p
}

func InitEmptyQuotedTermContext(p *QuotedTermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_quotedTerm
}

func (*QuotedTermContext) IsQuotedTermContext() {}

func NewQuotedTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QuotedTermContext {
	var p = new(QuotedTermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_quotedTerm

	return p
}

func (s *QuotedTermContext) GetParser() antlr.Parser { return s.parser }

func (s *QuotedTermContext) QUOTED() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserQUOTED, 0)
}

func (s *QuotedTermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuotedTermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QuotedTermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterQuotedTerm(s)
	}
}

func (s *QuotedTermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitQuotedTerm(s)
	}
}

func (s *QuotedTermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitQuotedTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) QuotedTerm() (localctx IQuotedTermContext) {
	localctx = NewQuotedTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, EDCqlLogParserRULE_quotedTerm)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(168)
		p.Match(EDCqlLogParserQUOTED)
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

// IOperatorColonContext is an interface to support dynamic dispatch.
type IOperatorColonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR() antlr.TerminalNode
	OP_COLON() antlr.TerminalNode

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
	p.RuleIndex = EDCqlLogParserRULE_operatorColon
	return p
}

func InitEmptyOperatorColonContext(p *OperatorColonContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlLogParserRULE_operatorColon
}

func (*OperatorColonContext) IsOperatorColonContext() {}

func NewOperatorColonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperatorColonContext {
	var p = new(OperatorColonContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlLogParserRULE_operatorColon

	return p
}

func (s *OperatorColonContext) GetParser() antlr.Parser { return s.parser }

func (s *OperatorColonContext) OPERATOR() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserOPERATOR, 0)
}

func (s *OperatorColonContext) OP_COLON() antlr.TerminalNode {
	return s.GetToken(EDCqlLogParserOP_COLON, 0)
}

func (s *OperatorColonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperatorColonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperatorColonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.EnterOperatorColon(s)
	}
}

func (s *OperatorColonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlLogParserListener); ok {
		listenerT.ExitOperatorColon(s)
	}
}

func (s *OperatorColonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlLogParserVisitor:
		return t.VisitOperatorColon(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlLogParser) OperatorColon() (localctx IOperatorColonContext) {
	localctx = NewOperatorColonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, EDCqlLogParserRULE_operatorColon)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(170)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EDCqlLogParserOP_COLON || _la == EDCqlLogParserOPERATOR) {
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
