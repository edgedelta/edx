// Code generated from EDCqlFormulaParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlFormulaParser

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

type EDCqlFormulaParser struct {
	*antlr.BaseParser
}

var EDCqlFormulaParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func edcqlformulaparserParserInit() {
	staticData := &EDCqlFormulaParserParserStaticData
	staticData.LiteralNames = []string{
		"", "'('", "')'", "','", "'timeshift'", "'moving_average'",
	}
	staticData.SymbolicNames = []string{
		"", "LPAREN", "RPAREN", "COMMA", "TIMESHIFT", "MOVING_AVERAGE", "NUMERIC",
		"FREE_TOKEN", "ARITHMETIC_TOKEN", "DEFAULT_SKIP", "UNKNOWN",
	}
	staticData.RuleNames = []string{
		"formula", "fullExpression", "part", "nestablePart", "nestableGroup",
		"outerFunction", "nestableFunction", "timeshift", "movingAverage", "group",
		"rawGroup", "rawPart", "rawExpr", "movingAverageExpr", "shift", "movingAveragePeriod",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 10, 122, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		1, 0, 1, 0, 1, 0, 1, 1, 4, 1, 37, 8, 1, 11, 1, 12, 1, 38, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 2, 3, 2, 46, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 53,
		8, 3, 1, 4, 1, 4, 4, 4, 57, 8, 4, 11, 4, 12, 4, 58, 1, 4, 1, 4, 1, 5, 1,
		5, 3, 5, 65, 8, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10,
		1, 10, 4, 10, 89, 8, 10, 11, 10, 12, 10, 90, 1, 10, 1, 10, 1, 11, 1, 11,
		1, 11, 1, 11, 3, 11, 99, 8, 11, 1, 12, 1, 12, 1, 12, 1, 12, 4, 12, 105,
		8, 12, 11, 12, 12, 12, 106, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 4, 13, 114,
		8, 13, 11, 13, 12, 13, 115, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 0, 0, 16,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 0, 0, 129, 0,
		32, 1, 0, 0, 0, 2, 36, 1, 0, 0, 0, 4, 45, 1, 0, 0, 0, 6, 52, 1, 0, 0, 0,
		8, 54, 1, 0, 0, 0, 10, 64, 1, 0, 0, 0, 12, 66, 1, 0, 0, 0, 14, 68, 1, 0,
		0, 0, 16, 75, 1, 0, 0, 0, 18, 82, 1, 0, 0, 0, 20, 86, 1, 0, 0, 0, 22, 98,
		1, 0, 0, 0, 24, 104, 1, 0, 0, 0, 26, 113, 1, 0, 0, 0, 28, 117, 1, 0, 0,
		0, 30, 119, 1, 0, 0, 0, 32, 33, 3, 2, 1, 0, 33, 34, 5, 0, 0, 1, 34, 1,
		1, 0, 0, 0, 35, 37, 3, 4, 2, 0, 36, 35, 1, 0, 0, 0, 37, 38, 1, 0, 0, 0,
		38, 36, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39, 3, 1, 0, 0, 0, 40, 46, 3, 10,
		5, 0, 41, 46, 3, 18, 9, 0, 42, 46, 5, 7, 0, 0, 43, 46, 5, 8, 0, 0, 44,
		46, 5, 6, 0, 0, 45, 40, 1, 0, 0, 0, 45, 41, 1, 0, 0, 0, 45, 42, 1, 0, 0,
		0, 45, 43, 1, 0, 0, 0, 45, 44, 1, 0, 0, 0, 46, 5, 1, 0, 0, 0, 47, 53, 3,
		12, 6, 0, 48, 53, 3, 18, 9, 0, 49, 53, 5, 7, 0, 0, 50, 53, 5, 8, 0, 0,
		51, 53, 5, 6, 0, 0, 52, 47, 1, 0, 0, 0, 52, 48, 1, 0, 0, 0, 52, 49, 1,
		0, 0, 0, 52, 50, 1, 0, 0, 0, 52, 51, 1, 0, 0, 0, 53, 7, 1, 0, 0, 0, 54,
		56, 5, 1, 0, 0, 55, 57, 3, 6, 3, 0, 56, 55, 1, 0, 0, 0, 57, 58, 1, 0, 0,
		0, 58, 56, 1, 0, 0, 0, 58, 59, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0, 60, 61,
		5, 2, 0, 0, 61, 9, 1, 0, 0, 0, 62, 65, 3, 14, 7, 0, 63, 65, 3, 16, 8, 0,
		64, 62, 1, 0, 0, 0, 64, 63, 1, 0, 0, 0, 65, 11, 1, 0, 0, 0, 66, 67, 3,
		14, 7, 0, 67, 13, 1, 0, 0, 0, 68, 69, 5, 4, 0, 0, 69, 70, 5, 1, 0, 0, 70,
		71, 3, 24, 12, 0, 71, 72, 5, 3, 0, 0, 72, 73, 3, 28, 14, 0, 73, 74, 5,
		2, 0, 0, 74, 15, 1, 0, 0, 0, 75, 76, 5, 5, 0, 0, 76, 77, 5, 1, 0, 0, 77,
		78, 3, 26, 13, 0, 78, 79, 5, 3, 0, 0, 79, 80, 3, 30, 15, 0, 80, 81, 5,
		2, 0, 0, 81, 17, 1, 0, 0, 0, 82, 83, 5, 1, 0, 0, 83, 84, 3, 2, 1, 0, 84,
		85, 5, 2, 0, 0, 85, 19, 1, 0, 0, 0, 86, 88, 5, 1, 0, 0, 87, 89, 3, 22,
		11, 0, 88, 87, 1, 0, 0, 0, 89, 90, 1, 0, 0, 0, 90, 88, 1, 0, 0, 0, 90,
		91, 1, 0, 0, 0, 91, 92, 1, 0, 0, 0, 92, 93, 5, 2, 0, 0, 93, 21, 1, 0, 0,
		0, 94, 99, 3, 20, 10, 0, 95, 99, 5, 7, 0, 0, 96, 99, 5, 8, 0, 0, 97, 99,
		5, 6, 0, 0, 98, 94, 1, 0, 0, 0, 98, 95, 1, 0, 0, 0, 98, 96, 1, 0, 0, 0,
		98, 97, 1, 0, 0, 0, 99, 23, 1, 0, 0, 0, 100, 105, 3, 20, 10, 0, 101, 105,
		5, 7, 0, 0, 102, 105, 5, 8, 0, 0, 103, 105, 5, 6, 0, 0, 104, 100, 1, 0,
		0, 0, 104, 101, 1, 0, 0, 0, 104, 102, 1, 0, 0, 0, 104, 103, 1, 0, 0, 0,
		105, 106, 1, 0, 0, 0, 106, 104, 1, 0, 0, 0, 106, 107, 1, 0, 0, 0, 107,
		25, 1, 0, 0, 0, 108, 114, 3, 12, 6, 0, 109, 114, 3, 8, 4, 0, 110, 114,
		5, 7, 0, 0, 111, 114, 5, 8, 0, 0, 112, 114, 5, 6, 0, 0, 113, 108, 1, 0,
		0, 0, 113, 109, 1, 0, 0, 0, 113, 110, 1, 0, 0, 0, 113, 111, 1, 0, 0, 0,
		113, 112, 1, 0, 0, 0, 114, 115, 1, 0, 0, 0, 115, 113, 1, 0, 0, 0, 115,
		116, 1, 0, 0, 0, 116, 27, 1, 0, 0, 0, 117, 118, 5, 6, 0, 0, 118, 29, 1,
		0, 0, 0, 119, 120, 5, 6, 0, 0, 120, 31, 1, 0, 0, 0, 11, 38, 45, 52, 58,
		64, 90, 98, 104, 106, 113, 115,
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

// EDCqlFormulaParserInit initializes any static state used to implement EDCqlFormulaParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewEDCqlFormulaParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func EDCqlFormulaParserInit() {
	staticData := &EDCqlFormulaParserParserStaticData
	staticData.once.Do(edcqlformulaparserParserInit)
}

// NewEDCqlFormulaParser produces a new parser instance for the optional input antlr.TokenStream.
func NewEDCqlFormulaParser(input antlr.TokenStream) *EDCqlFormulaParser {
	EDCqlFormulaParserInit()
	this := new(EDCqlFormulaParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &EDCqlFormulaParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "EDCqlFormulaParser.g4"

	return this
}

// EDCqlFormulaParser tokens.
const (
	EDCqlFormulaParserEOF              = antlr.TokenEOF
	EDCqlFormulaParserLPAREN           = 1
	EDCqlFormulaParserRPAREN           = 2
	EDCqlFormulaParserCOMMA            = 3
	EDCqlFormulaParserTIMESHIFT        = 4
	EDCqlFormulaParserMOVING_AVERAGE   = 5
	EDCqlFormulaParserNUMERIC          = 6
	EDCqlFormulaParserFREE_TOKEN       = 7
	EDCqlFormulaParserARITHMETIC_TOKEN = 8
	EDCqlFormulaParserDEFAULT_SKIP     = 9
	EDCqlFormulaParserUNKNOWN          = 10
)

// EDCqlFormulaParser rules.
const (
	EDCqlFormulaParserRULE_formula             = 0
	EDCqlFormulaParserRULE_fullExpression      = 1
	EDCqlFormulaParserRULE_part                = 2
	EDCqlFormulaParserRULE_nestablePart        = 3
	EDCqlFormulaParserRULE_nestableGroup       = 4
	EDCqlFormulaParserRULE_outerFunction       = 5
	EDCqlFormulaParserRULE_nestableFunction    = 6
	EDCqlFormulaParserRULE_timeshift           = 7
	EDCqlFormulaParserRULE_movingAverage       = 8
	EDCqlFormulaParserRULE_group               = 9
	EDCqlFormulaParserRULE_rawGroup            = 10
	EDCqlFormulaParserRULE_rawPart             = 11
	EDCqlFormulaParserRULE_rawExpr             = 12
	EDCqlFormulaParserRULE_movingAverageExpr   = 13
	EDCqlFormulaParserRULE_shift               = 14
	EDCqlFormulaParserRULE_movingAveragePeriod = 15
)

// IFormulaContext is an interface to support dynamic dispatch.
type IFormulaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FullExpression() IFullExpressionContext
	EOF() antlr.TerminalNode

	// IsFormulaContext differentiates from other interfaces.
	IsFormulaContext()
}

type FormulaContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFormulaContext() *FormulaContext {
	var p = new(FormulaContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_formula
	return p
}

func InitEmptyFormulaContext(p *FormulaContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_formula
}

func (*FormulaContext) IsFormulaContext() {}

func NewFormulaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FormulaContext {
	var p = new(FormulaContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_formula

	return p
}

func (s *FormulaContext) GetParser() antlr.Parser { return s.parser }

func (s *FormulaContext) FullExpression() IFullExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFullExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFullExpressionContext)
}

func (s *FormulaContext) EOF() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserEOF, 0)
}

func (s *FormulaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FormulaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FormulaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterFormula(s)
	}
}

func (s *FormulaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitFormula(s)
	}
}

func (s *FormulaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitFormula(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) Formula() (localctx IFormulaContext) {
	localctx = NewFormulaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, EDCqlFormulaParserRULE_formula)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(32)
		p.FullExpression()
	}
	{
		p.SetState(33)
		p.Match(EDCqlFormulaParserEOF)
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

// IFullExpressionContext is an interface to support dynamic dispatch.
type IFullExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPart() []IPartContext
	Part(i int) IPartContext

	// IsFullExpressionContext differentiates from other interfaces.
	IsFullExpressionContext()
}

type FullExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFullExpressionContext() *FullExpressionContext {
	var p = new(FullExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_fullExpression
	return p
}

func InitEmptyFullExpressionContext(p *FullExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_fullExpression
}

func (*FullExpressionContext) IsFullExpressionContext() {}

func NewFullExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FullExpressionContext {
	var p = new(FullExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_fullExpression

	return p
}

func (s *FullExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *FullExpressionContext) AllPart() []IPartContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPartContext); ok {
			len++
		}
	}

	tst := make([]IPartContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPartContext); ok {
			tst[i] = t.(IPartContext)
			i++
		}
	}

	return tst
}

func (s *FullExpressionContext) Part(i int) IPartContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPartContext); ok {
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

	return t.(IPartContext)
}

func (s *FullExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FullExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FullExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterFullExpression(s)
	}
}

func (s *FullExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitFullExpression(s)
	}
}

func (s *FullExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitFullExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) FullExpression() (localctx IFullExpressionContext) {
	localctx = NewFullExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, EDCqlFormulaParserRULE_fullExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(36)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&498) != 0) {
		{
			p.SetState(35)
			p.Part()
		}

		p.SetState(38)
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

// IPartContext is an interface to support dynamic dispatch.
type IPartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OuterFunction() IOuterFunctionContext
	Group() IGroupContext
	FREE_TOKEN() antlr.TerminalNode
	ARITHMETIC_TOKEN() antlr.TerminalNode
	NUMERIC() antlr.TerminalNode

	// IsPartContext differentiates from other interfaces.
	IsPartContext()
}

type PartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPartContext() *PartContext {
	var p = new(PartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_part
	return p
}

func InitEmptyPartContext(p *PartContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_part
}

func (*PartContext) IsPartContext() {}

func NewPartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PartContext {
	var p = new(PartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_part

	return p
}

func (s *PartContext) GetParser() antlr.Parser { return s.parser }

func (s *PartContext) OuterFunction() IOuterFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOuterFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOuterFunctionContext)
}

func (s *PartContext) Group() IGroupContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupContext)
}

func (s *PartContext) FREE_TOKEN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserFREE_TOKEN, 0)
}

func (s *PartContext) ARITHMETIC_TOKEN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserARITHMETIC_TOKEN, 0)
}

func (s *PartContext) NUMERIC() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserNUMERIC, 0)
}

func (s *PartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterPart(s)
	}
}

func (s *PartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitPart(s)
	}
}

func (s *PartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitPart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) Part() (localctx IPartContext) {
	localctx = NewPartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, EDCqlFormulaParserRULE_part)
	p.SetState(45)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlFormulaParserTIMESHIFT, EDCqlFormulaParserMOVING_AVERAGE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(40)
			p.OuterFunction()
		}

	case EDCqlFormulaParserLPAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(41)
			p.Group()
		}

	case EDCqlFormulaParserFREE_TOKEN:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(42)
			p.Match(EDCqlFormulaParserFREE_TOKEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlFormulaParserARITHMETIC_TOKEN:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(43)
			p.Match(EDCqlFormulaParserARITHMETIC_TOKEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlFormulaParserNUMERIC:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(44)
			p.Match(EDCqlFormulaParserNUMERIC)
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

// INestablePartContext is an interface to support dynamic dispatch.
type INestablePartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NestableFunction() INestableFunctionContext
	Group() IGroupContext
	FREE_TOKEN() antlr.TerminalNode
	ARITHMETIC_TOKEN() antlr.TerminalNode
	NUMERIC() antlr.TerminalNode

	// IsNestablePartContext differentiates from other interfaces.
	IsNestablePartContext()
}

type NestablePartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNestablePartContext() *NestablePartContext {
	var p = new(NestablePartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_nestablePart
	return p
}

func InitEmptyNestablePartContext(p *NestablePartContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_nestablePart
}

func (*NestablePartContext) IsNestablePartContext() {}

func NewNestablePartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NestablePartContext {
	var p = new(NestablePartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_nestablePart

	return p
}

func (s *NestablePartContext) GetParser() antlr.Parser { return s.parser }

func (s *NestablePartContext) NestableFunction() INestableFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INestableFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INestableFunctionContext)
}

func (s *NestablePartContext) Group() IGroupContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupContext)
}

func (s *NestablePartContext) FREE_TOKEN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserFREE_TOKEN, 0)
}

func (s *NestablePartContext) ARITHMETIC_TOKEN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserARITHMETIC_TOKEN, 0)
}

func (s *NestablePartContext) NUMERIC() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserNUMERIC, 0)
}

func (s *NestablePartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NestablePartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NestablePartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterNestablePart(s)
	}
}

func (s *NestablePartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitNestablePart(s)
	}
}

func (s *NestablePartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitNestablePart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) NestablePart() (localctx INestablePartContext) {
	localctx = NewNestablePartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, EDCqlFormulaParserRULE_nestablePart)
	p.SetState(52)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlFormulaParserTIMESHIFT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(47)
			p.NestableFunction()
		}

	case EDCqlFormulaParserLPAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(48)
			p.Group()
		}

	case EDCqlFormulaParserFREE_TOKEN:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(49)
			p.Match(EDCqlFormulaParserFREE_TOKEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlFormulaParserARITHMETIC_TOKEN:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(50)
			p.Match(EDCqlFormulaParserARITHMETIC_TOKEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlFormulaParserNUMERIC:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(51)
			p.Match(EDCqlFormulaParserNUMERIC)
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

// INestableGroupContext is an interface to support dynamic dispatch.
type INestableGroupContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllNestablePart() []INestablePartContext
	NestablePart(i int) INestablePartContext

	// IsNestableGroupContext differentiates from other interfaces.
	IsNestableGroupContext()
}

type NestableGroupContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNestableGroupContext() *NestableGroupContext {
	var p = new(NestableGroupContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_nestableGroup
	return p
}

func InitEmptyNestableGroupContext(p *NestableGroupContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_nestableGroup
}

func (*NestableGroupContext) IsNestableGroupContext() {}

func NewNestableGroupContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NestableGroupContext {
	var p = new(NestableGroupContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_nestableGroup

	return p
}

func (s *NestableGroupContext) GetParser() antlr.Parser { return s.parser }

func (s *NestableGroupContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserLPAREN, 0)
}

func (s *NestableGroupContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserRPAREN, 0)
}

func (s *NestableGroupContext) AllNestablePart() []INestablePartContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INestablePartContext); ok {
			len++
		}
	}

	tst := make([]INestablePartContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INestablePartContext); ok {
			tst[i] = t.(INestablePartContext)
			i++
		}
	}

	return tst
}

func (s *NestableGroupContext) NestablePart(i int) INestablePartContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INestablePartContext); ok {
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

	return t.(INestablePartContext)
}

func (s *NestableGroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NestableGroupContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NestableGroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterNestableGroup(s)
	}
}

func (s *NestableGroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitNestableGroup(s)
	}
}

func (s *NestableGroupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitNestableGroup(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) NestableGroup() (localctx INestableGroupContext) {
	localctx = NewNestableGroupContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, EDCqlFormulaParserRULE_nestableGroup)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(54)
		p.Match(EDCqlFormulaParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(56)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&466) != 0) {
		{
			p.SetState(55)
			p.NestablePart()
		}

		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(60)
		p.Match(EDCqlFormulaParserRPAREN)
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

// IOuterFunctionContext is an interface to support dynamic dispatch.
type IOuterFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Timeshift() ITimeshiftContext
	MovingAverage() IMovingAverageContext

	// IsOuterFunctionContext differentiates from other interfaces.
	IsOuterFunctionContext()
}

type OuterFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOuterFunctionContext() *OuterFunctionContext {
	var p = new(OuterFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_outerFunction
	return p
}

func InitEmptyOuterFunctionContext(p *OuterFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_outerFunction
}

func (*OuterFunctionContext) IsOuterFunctionContext() {}

func NewOuterFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OuterFunctionContext {
	var p = new(OuterFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_outerFunction

	return p
}

func (s *OuterFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *OuterFunctionContext) Timeshift() ITimeshiftContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeshiftContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeshiftContext)
}

func (s *OuterFunctionContext) MovingAverage() IMovingAverageContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMovingAverageContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMovingAverageContext)
}

func (s *OuterFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OuterFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OuterFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterOuterFunction(s)
	}
}

func (s *OuterFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitOuterFunction(s)
	}
}

func (s *OuterFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitOuterFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) OuterFunction() (localctx IOuterFunctionContext) {
	localctx = NewOuterFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, EDCqlFormulaParserRULE_outerFunction)
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlFormulaParserTIMESHIFT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(62)
			p.Timeshift()
		}

	case EDCqlFormulaParserMOVING_AVERAGE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(63)
			p.MovingAverage()
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

// INestableFunctionContext is an interface to support dynamic dispatch.
type INestableFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Timeshift() ITimeshiftContext

	// IsNestableFunctionContext differentiates from other interfaces.
	IsNestableFunctionContext()
}

type NestableFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNestableFunctionContext() *NestableFunctionContext {
	var p = new(NestableFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_nestableFunction
	return p
}

func InitEmptyNestableFunctionContext(p *NestableFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_nestableFunction
}

func (*NestableFunctionContext) IsNestableFunctionContext() {}

func NewNestableFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NestableFunctionContext {
	var p = new(NestableFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_nestableFunction

	return p
}

func (s *NestableFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *NestableFunctionContext) Timeshift() ITimeshiftContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeshiftContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeshiftContext)
}

func (s *NestableFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NestableFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NestableFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterNestableFunction(s)
	}
}

func (s *NestableFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitNestableFunction(s)
	}
}

func (s *NestableFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitNestableFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) NestableFunction() (localctx INestableFunctionContext) {
	localctx = NewNestableFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, EDCqlFormulaParserRULE_nestableFunction)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(66)
		p.Timeshift()
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

// ITimeshiftContext is an interface to support dynamic dispatch.
type ITimeshiftContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TIMESHIFT() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RawExpr() IRawExprContext
	COMMA() antlr.TerminalNode
	Shift() IShiftContext
	RPAREN() antlr.TerminalNode

	// IsTimeshiftContext differentiates from other interfaces.
	IsTimeshiftContext()
}

type TimeshiftContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeshiftContext() *TimeshiftContext {
	var p = new(TimeshiftContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_timeshift
	return p
}

func InitEmptyTimeshiftContext(p *TimeshiftContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_timeshift
}

func (*TimeshiftContext) IsTimeshiftContext() {}

func NewTimeshiftContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeshiftContext {
	var p = new(TimeshiftContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_timeshift

	return p
}

func (s *TimeshiftContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeshiftContext) TIMESHIFT() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserTIMESHIFT, 0)
}

func (s *TimeshiftContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserLPAREN, 0)
}

func (s *TimeshiftContext) RawExpr() IRawExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRawExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRawExprContext)
}

func (s *TimeshiftContext) COMMA() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserCOMMA, 0)
}

func (s *TimeshiftContext) Shift() IShiftContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShiftContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShiftContext)
}

func (s *TimeshiftContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserRPAREN, 0)
}

func (s *TimeshiftContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeshiftContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeshiftContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterTimeshift(s)
	}
}

func (s *TimeshiftContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitTimeshift(s)
	}
}

func (s *TimeshiftContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitTimeshift(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) Timeshift() (localctx ITimeshiftContext) {
	localctx = NewTimeshiftContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, EDCqlFormulaParserRULE_timeshift)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(68)
		p.Match(EDCqlFormulaParserTIMESHIFT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(69)
		p.Match(EDCqlFormulaParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(70)
		p.RawExpr()
	}
	{
		p.SetState(71)
		p.Match(EDCqlFormulaParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(72)
		p.Shift()
	}
	{
		p.SetState(73)
		p.Match(EDCqlFormulaParserRPAREN)
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

// IMovingAverageContext is an interface to support dynamic dispatch.
type IMovingAverageContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MOVING_AVERAGE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	MovingAverageExpr() IMovingAverageExprContext
	COMMA() antlr.TerminalNode
	MovingAveragePeriod() IMovingAveragePeriodContext
	RPAREN() antlr.TerminalNode

	// IsMovingAverageContext differentiates from other interfaces.
	IsMovingAverageContext()
}

type MovingAverageContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMovingAverageContext() *MovingAverageContext {
	var p = new(MovingAverageContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_movingAverage
	return p
}

func InitEmptyMovingAverageContext(p *MovingAverageContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_movingAverage
}

func (*MovingAverageContext) IsMovingAverageContext() {}

func NewMovingAverageContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MovingAverageContext {
	var p = new(MovingAverageContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_movingAverage

	return p
}

func (s *MovingAverageContext) GetParser() antlr.Parser { return s.parser }

func (s *MovingAverageContext) MOVING_AVERAGE() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserMOVING_AVERAGE, 0)
}

func (s *MovingAverageContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserLPAREN, 0)
}

func (s *MovingAverageContext) MovingAverageExpr() IMovingAverageExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMovingAverageExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMovingAverageExprContext)
}

func (s *MovingAverageContext) COMMA() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserCOMMA, 0)
}

func (s *MovingAverageContext) MovingAveragePeriod() IMovingAveragePeriodContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMovingAveragePeriodContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMovingAveragePeriodContext)
}

func (s *MovingAverageContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserRPAREN, 0)
}

func (s *MovingAverageContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MovingAverageContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MovingAverageContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterMovingAverage(s)
	}
}

func (s *MovingAverageContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitMovingAverage(s)
	}
}

func (s *MovingAverageContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitMovingAverage(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) MovingAverage() (localctx IMovingAverageContext) {
	localctx = NewMovingAverageContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, EDCqlFormulaParserRULE_movingAverage)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		p.Match(EDCqlFormulaParserMOVING_AVERAGE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(76)
		p.Match(EDCqlFormulaParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(77)
		p.MovingAverageExpr()
	}
	{
		p.SetState(78)
		p.Match(EDCqlFormulaParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(79)
		p.MovingAveragePeriod()
	}
	{
		p.SetState(80)
		p.Match(EDCqlFormulaParserRPAREN)
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

// IGroupContext is an interface to support dynamic dispatch.
type IGroupContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	FullExpression() IFullExpressionContext
	RPAREN() antlr.TerminalNode

	// IsGroupContext differentiates from other interfaces.
	IsGroupContext()
}

type GroupContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupContext() *GroupContext {
	var p = new(GroupContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_group
	return p
}

func InitEmptyGroupContext(p *GroupContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_group
}

func (*GroupContext) IsGroupContext() {}

func NewGroupContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupContext {
	var p = new(GroupContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_group

	return p
}

func (s *GroupContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserLPAREN, 0)
}

func (s *GroupContext) FullExpression() IFullExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFullExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFullExpressionContext)
}

func (s *GroupContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserRPAREN, 0)
}

func (s *GroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterGroup(s)
	}
}

func (s *GroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitGroup(s)
	}
}

func (s *GroupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitGroup(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) Group() (localctx IGroupContext) {
	localctx = NewGroupContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, EDCqlFormulaParserRULE_group)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(82)
		p.Match(EDCqlFormulaParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(83)
		p.FullExpression()
	}
	{
		p.SetState(84)
		p.Match(EDCqlFormulaParserRPAREN)
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

// IRawGroupContext is an interface to support dynamic dispatch.
type IRawGroupContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllRawPart() []IRawPartContext
	RawPart(i int) IRawPartContext

	// IsRawGroupContext differentiates from other interfaces.
	IsRawGroupContext()
}

type RawGroupContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRawGroupContext() *RawGroupContext {
	var p = new(RawGroupContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_rawGroup
	return p
}

func InitEmptyRawGroupContext(p *RawGroupContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_rawGroup
}

func (*RawGroupContext) IsRawGroupContext() {}

func NewRawGroupContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawGroupContext {
	var p = new(RawGroupContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_rawGroup

	return p
}

func (s *RawGroupContext) GetParser() antlr.Parser { return s.parser }

func (s *RawGroupContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserLPAREN, 0)
}

func (s *RawGroupContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserRPAREN, 0)
}

func (s *RawGroupContext) AllRawPart() []IRawPartContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRawPartContext); ok {
			len++
		}
	}

	tst := make([]IRawPartContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRawPartContext); ok {
			tst[i] = t.(IRawPartContext)
			i++
		}
	}

	return tst
}

func (s *RawGroupContext) RawPart(i int) IRawPartContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRawPartContext); ok {
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

	return t.(IRawPartContext)
}

func (s *RawGroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RawGroupContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RawGroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterRawGroup(s)
	}
}

func (s *RawGroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitRawGroup(s)
	}
}

func (s *RawGroupContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitRawGroup(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) RawGroup() (localctx IRawGroupContext) {
	localctx = NewRawGroupContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, EDCqlFormulaParserRULE_rawGroup)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(86)
		p.Match(EDCqlFormulaParserLPAREN)
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

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&450) != 0) {
		{
			p.SetState(87)
			p.RawPart()
		}

		p.SetState(90)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(92)
		p.Match(EDCqlFormulaParserRPAREN)
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

// IRawPartContext is an interface to support dynamic dispatch.
type IRawPartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RawGroup() IRawGroupContext
	FREE_TOKEN() antlr.TerminalNode
	ARITHMETIC_TOKEN() antlr.TerminalNode
	NUMERIC() antlr.TerminalNode

	// IsRawPartContext differentiates from other interfaces.
	IsRawPartContext()
}

type RawPartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRawPartContext() *RawPartContext {
	var p = new(RawPartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_rawPart
	return p
}

func InitEmptyRawPartContext(p *RawPartContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_rawPart
}

func (*RawPartContext) IsRawPartContext() {}

func NewRawPartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawPartContext {
	var p = new(RawPartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_rawPart

	return p
}

func (s *RawPartContext) GetParser() antlr.Parser { return s.parser }

func (s *RawPartContext) RawGroup() IRawGroupContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRawGroupContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRawGroupContext)
}

func (s *RawPartContext) FREE_TOKEN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserFREE_TOKEN, 0)
}

func (s *RawPartContext) ARITHMETIC_TOKEN() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserARITHMETIC_TOKEN, 0)
}

func (s *RawPartContext) NUMERIC() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserNUMERIC, 0)
}

func (s *RawPartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RawPartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RawPartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterRawPart(s)
	}
}

func (s *RawPartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitRawPart(s)
	}
}

func (s *RawPartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitRawPart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) RawPart() (localctx IRawPartContext) {
	localctx = NewRawPartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, EDCqlFormulaParserRULE_rawPart)
	p.SetState(98)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EDCqlFormulaParserLPAREN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(94)
			p.RawGroup()
		}

	case EDCqlFormulaParserFREE_TOKEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(95)
			p.Match(EDCqlFormulaParserFREE_TOKEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlFormulaParserARITHMETIC_TOKEN:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(96)
			p.Match(EDCqlFormulaParserARITHMETIC_TOKEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EDCqlFormulaParserNUMERIC:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(97)
			p.Match(EDCqlFormulaParserNUMERIC)
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

// IRawExprContext is an interface to support dynamic dispatch.
type IRawExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRawGroup() []IRawGroupContext
	RawGroup(i int) IRawGroupContext
	AllFREE_TOKEN() []antlr.TerminalNode
	FREE_TOKEN(i int) antlr.TerminalNode
	AllARITHMETIC_TOKEN() []antlr.TerminalNode
	ARITHMETIC_TOKEN(i int) antlr.TerminalNode
	AllNUMERIC() []antlr.TerminalNode
	NUMERIC(i int) antlr.TerminalNode

	// IsRawExprContext differentiates from other interfaces.
	IsRawExprContext()
}

type RawExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRawExprContext() *RawExprContext {
	var p = new(RawExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_rawExpr
	return p
}

func InitEmptyRawExprContext(p *RawExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_rawExpr
}

func (*RawExprContext) IsRawExprContext() {}

func NewRawExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RawExprContext {
	var p = new(RawExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_rawExpr

	return p
}

func (s *RawExprContext) GetParser() antlr.Parser { return s.parser }

func (s *RawExprContext) AllRawGroup() []IRawGroupContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRawGroupContext); ok {
			len++
		}
	}

	tst := make([]IRawGroupContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRawGroupContext); ok {
			tst[i] = t.(IRawGroupContext)
			i++
		}
	}

	return tst
}

func (s *RawExprContext) RawGroup(i int) IRawGroupContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRawGroupContext); ok {
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

	return t.(IRawGroupContext)
}

func (s *RawExprContext) AllFREE_TOKEN() []antlr.TerminalNode {
	return s.GetTokens(EDCqlFormulaParserFREE_TOKEN)
}

func (s *RawExprContext) FREE_TOKEN(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserFREE_TOKEN, i)
}

func (s *RawExprContext) AllARITHMETIC_TOKEN() []antlr.TerminalNode {
	return s.GetTokens(EDCqlFormulaParserARITHMETIC_TOKEN)
}

func (s *RawExprContext) ARITHMETIC_TOKEN(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserARITHMETIC_TOKEN, i)
}

func (s *RawExprContext) AllNUMERIC() []antlr.TerminalNode {
	return s.GetTokens(EDCqlFormulaParserNUMERIC)
}

func (s *RawExprContext) NUMERIC(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserNUMERIC, i)
}

func (s *RawExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RawExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RawExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterRawExpr(s)
	}
}

func (s *RawExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitRawExpr(s)
	}
}

func (s *RawExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitRawExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) RawExpr() (localctx IRawExprContext) {
	localctx = NewRawExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, EDCqlFormulaParserRULE_rawExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(104)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&450) != 0) {
		p.SetState(104)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case EDCqlFormulaParserLPAREN:
			{
				p.SetState(100)
				p.RawGroup()
			}

		case EDCqlFormulaParserFREE_TOKEN:
			{
				p.SetState(101)
				p.Match(EDCqlFormulaParserFREE_TOKEN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case EDCqlFormulaParserARITHMETIC_TOKEN:
			{
				p.SetState(102)
				p.Match(EDCqlFormulaParserARITHMETIC_TOKEN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case EDCqlFormulaParserNUMERIC:
			{
				p.SetState(103)
				p.Match(EDCqlFormulaParserNUMERIC)
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

// IMovingAverageExprContext is an interface to support dynamic dispatch.
type IMovingAverageExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNestableFunction() []INestableFunctionContext
	NestableFunction(i int) INestableFunctionContext
	AllNestableGroup() []INestableGroupContext
	NestableGroup(i int) INestableGroupContext
	AllFREE_TOKEN() []antlr.TerminalNode
	FREE_TOKEN(i int) antlr.TerminalNode
	AllARITHMETIC_TOKEN() []antlr.TerminalNode
	ARITHMETIC_TOKEN(i int) antlr.TerminalNode
	AllNUMERIC() []antlr.TerminalNode
	NUMERIC(i int) antlr.TerminalNode

	// IsMovingAverageExprContext differentiates from other interfaces.
	IsMovingAverageExprContext()
}

type MovingAverageExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMovingAverageExprContext() *MovingAverageExprContext {
	var p = new(MovingAverageExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_movingAverageExpr
	return p
}

func InitEmptyMovingAverageExprContext(p *MovingAverageExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_movingAverageExpr
}

func (*MovingAverageExprContext) IsMovingAverageExprContext() {}

func NewMovingAverageExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MovingAverageExprContext {
	var p = new(MovingAverageExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_movingAverageExpr

	return p
}

func (s *MovingAverageExprContext) GetParser() antlr.Parser { return s.parser }

func (s *MovingAverageExprContext) AllNestableFunction() []INestableFunctionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INestableFunctionContext); ok {
			len++
		}
	}

	tst := make([]INestableFunctionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INestableFunctionContext); ok {
			tst[i] = t.(INestableFunctionContext)
			i++
		}
	}

	return tst
}

func (s *MovingAverageExprContext) NestableFunction(i int) INestableFunctionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INestableFunctionContext); ok {
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

	return t.(INestableFunctionContext)
}

func (s *MovingAverageExprContext) AllNestableGroup() []INestableGroupContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INestableGroupContext); ok {
			len++
		}
	}

	tst := make([]INestableGroupContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INestableGroupContext); ok {
			tst[i] = t.(INestableGroupContext)
			i++
		}
	}

	return tst
}

func (s *MovingAverageExprContext) NestableGroup(i int) INestableGroupContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INestableGroupContext); ok {
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

	return t.(INestableGroupContext)
}

func (s *MovingAverageExprContext) AllFREE_TOKEN() []antlr.TerminalNode {
	return s.GetTokens(EDCqlFormulaParserFREE_TOKEN)
}

func (s *MovingAverageExprContext) FREE_TOKEN(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserFREE_TOKEN, i)
}

func (s *MovingAverageExprContext) AllARITHMETIC_TOKEN() []antlr.TerminalNode {
	return s.GetTokens(EDCqlFormulaParserARITHMETIC_TOKEN)
}

func (s *MovingAverageExprContext) ARITHMETIC_TOKEN(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserARITHMETIC_TOKEN, i)
}

func (s *MovingAverageExprContext) AllNUMERIC() []antlr.TerminalNode {
	return s.GetTokens(EDCqlFormulaParserNUMERIC)
}

func (s *MovingAverageExprContext) NUMERIC(i int) antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserNUMERIC, i)
}

func (s *MovingAverageExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MovingAverageExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MovingAverageExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterMovingAverageExpr(s)
	}
}

func (s *MovingAverageExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitMovingAverageExpr(s)
	}
}

func (s *MovingAverageExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitMovingAverageExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) MovingAverageExpr() (localctx IMovingAverageExprContext) {
	localctx = NewMovingAverageExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, EDCqlFormulaParserRULE_movingAverageExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&466) != 0) {
		p.SetState(113)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case EDCqlFormulaParserTIMESHIFT:
			{
				p.SetState(108)
				p.NestableFunction()
			}

		case EDCqlFormulaParserLPAREN:
			{
				p.SetState(109)
				p.NestableGroup()
			}

		case EDCqlFormulaParserFREE_TOKEN:
			{
				p.SetState(110)
				p.Match(EDCqlFormulaParserFREE_TOKEN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case EDCqlFormulaParserARITHMETIC_TOKEN:
			{
				p.SetState(111)
				p.Match(EDCqlFormulaParserARITHMETIC_TOKEN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case EDCqlFormulaParserNUMERIC:
			{
				p.SetState(112)
				p.Match(EDCqlFormulaParserNUMERIC)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(115)
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

// IShiftContext is an interface to support dynamic dispatch.
type IShiftContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMERIC() antlr.TerminalNode

	// IsShiftContext differentiates from other interfaces.
	IsShiftContext()
}

type ShiftContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShiftContext() *ShiftContext {
	var p = new(ShiftContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_shift
	return p
}

func InitEmptyShiftContext(p *ShiftContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_shift
}

func (*ShiftContext) IsShiftContext() {}

func NewShiftContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShiftContext {
	var p = new(ShiftContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_shift

	return p
}

func (s *ShiftContext) GetParser() antlr.Parser { return s.parser }

func (s *ShiftContext) NUMERIC() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserNUMERIC, 0)
}

func (s *ShiftContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShiftContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShiftContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterShift(s)
	}
}

func (s *ShiftContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitShift(s)
	}
}

func (s *ShiftContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitShift(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) Shift() (localctx IShiftContext) {
	localctx = NewShiftContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, EDCqlFormulaParserRULE_shift)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		p.Match(EDCqlFormulaParserNUMERIC)
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

// IMovingAveragePeriodContext is an interface to support dynamic dispatch.
type IMovingAveragePeriodContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMERIC() antlr.TerminalNode

	// IsMovingAveragePeriodContext differentiates from other interfaces.
	IsMovingAveragePeriodContext()
}

type MovingAveragePeriodContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMovingAveragePeriodContext() *MovingAveragePeriodContext {
	var p = new(MovingAveragePeriodContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_movingAveragePeriod
	return p
}

func InitEmptyMovingAveragePeriodContext(p *MovingAveragePeriodContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EDCqlFormulaParserRULE_movingAveragePeriod
}

func (*MovingAveragePeriodContext) IsMovingAveragePeriodContext() {}

func NewMovingAveragePeriodContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MovingAveragePeriodContext {
	var p = new(MovingAveragePeriodContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EDCqlFormulaParserRULE_movingAveragePeriod

	return p
}

func (s *MovingAveragePeriodContext) GetParser() antlr.Parser { return s.parser }

func (s *MovingAveragePeriodContext) NUMERIC() antlr.TerminalNode {
	return s.GetToken(EDCqlFormulaParserNUMERIC, 0)
}

func (s *MovingAveragePeriodContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MovingAveragePeriodContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MovingAveragePeriodContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.EnterMovingAveragePeriod(s)
	}
}

func (s *MovingAveragePeriodContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EDCqlFormulaParserListener); ok {
		listenerT.ExitMovingAveragePeriod(s)
	}
}

func (s *MovingAveragePeriodContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case EDCqlFormulaParserVisitor:
		return t.VisitMovingAveragePeriod(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *EDCqlFormulaParser) MovingAveragePeriod() (localctx IMovingAveragePeriodContext) {
	localctx = NewMovingAveragePeriodContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, EDCqlFormulaParserRULE_movingAveragePeriod)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
		p.Match(EDCqlFormulaParserNUMERIC)
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
