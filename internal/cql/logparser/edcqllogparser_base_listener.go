// Code generated from EDCqlLogParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package logparser // EDCqlLogParser
import "github.com/antlr4-go/antlr/v4"

// BaseEDCqlLogParserListener is a complete listener for a parse tree produced by EDCqlLogParser.
type BaseEDCqlLogParserListener struct{}

var _ EDCqlLogParserListener = &BaseEDCqlLogParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseEDCqlLogParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseEDCqlLogParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseEDCqlLogParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseEDCqlLogParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFinalQuery is called when production finalQuery is entered.
func (s *BaseEDCqlLogParserListener) EnterFinalQuery(ctx *FinalQueryContext) {}

// ExitFinalQuery is called when production finalQuery is exited.
func (s *BaseEDCqlLogParserListener) ExitFinalQuery(ctx *FinalQueryContext) {}

// EnterTopLevelQuery is called when production topLevelQuery is entered.
func (s *BaseEDCqlLogParserListener) EnterTopLevelQuery(ctx *TopLevelQueryContext) {}

// ExitTopLevelQuery is called when production topLevelQuery is exited.
func (s *BaseEDCqlLogParserListener) ExitTopLevelQuery(ctx *TopLevelQueryContext) {}

// EnterRollupMethod is called when production rollupMethod is entered.
func (s *BaseEDCqlLogParserListener) EnterRollupMethod(ctx *RollupMethodContext) {}

// ExitRollupMethod is called when production rollupMethod is exited.
func (s *BaseEDCqlLogParserListener) ExitRollupMethod(ctx *RollupMethodContext) {}

// EnterRollupField is called when production rollupField is entered.
func (s *BaseEDCqlLogParserListener) EnterRollupField(ctx *RollupFieldContext) {}

// ExitRollupField is called when production rollupField is exited.
func (s *BaseEDCqlLogParserListener) ExitRollupField(ctx *RollupFieldContext) {}

// EnterRollupSection is called when production rollupSection is entered.
func (s *BaseEDCqlLogParserListener) EnterRollupSection(ctx *RollupSectionContext) {}

// ExitRollupSection is called when production rollupSection is exited.
func (s *BaseEDCqlLogParserListener) ExitRollupSection(ctx *RollupSectionContext) {}

// EnterCountUniqueFields is called when production countUniqueFields is entered.
func (s *BaseEDCqlLogParserListener) EnterCountUniqueFields(ctx *CountUniqueFieldsContext) {}

// ExitCountUniqueFields is called when production countUniqueFields is exited.
func (s *BaseEDCqlLogParserListener) ExitCountUniqueFields(ctx *CountUniqueFieldsContext) {}

// EnterGroupByFields is called when production groupByFields is entered.
func (s *BaseEDCqlLogParserListener) EnterGroupByFields(ctx *GroupByFieldsContext) {}

// ExitGroupByFields is called when production groupByFields is exited.
func (s *BaseEDCqlLogParserListener) ExitGroupByFields(ctx *GroupByFieldsContext) {}

// EnterGroupBySection is called when production groupBySection is entered.
func (s *BaseEDCqlLogParserListener) EnterGroupBySection(ctx *GroupBySectionContext) {}

// ExitGroupBySection is called when production groupBySection is exited.
func (s *BaseEDCqlLogParserListener) ExitGroupBySection(ctx *GroupBySectionContext) {}

// EnterCountUniqueSection is called when production countUniqueSection is entered.
func (s *BaseEDCqlLogParserListener) EnterCountUniqueSection(ctx *CountUniqueSectionContext) {}

// ExitCountUniqueSection is called when production countUniqueSection is exited.
func (s *BaseEDCqlLogParserListener) ExitCountUniqueSection(ctx *CountUniqueSectionContext) {}

// EnterQueryOperation is called when production queryOperation is entered.
func (s *BaseEDCqlLogParserListener) EnterQueryOperation(ctx *QueryOperationContext) {}

// ExitQueryOperation is called when production queryOperation is exited.
func (s *BaseEDCqlLogParserListener) ExitQueryOperation(ctx *QueryOperationContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseEDCqlLogParserListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseEDCqlLogParserListener) ExitQuery(ctx *QueryContext) {}

// EnterDisjQuery is called when production disjQuery is entered.
func (s *BaseEDCqlLogParserListener) EnterDisjQuery(ctx *DisjQueryContext) {}

// ExitDisjQuery is called when production disjQuery is exited.
func (s *BaseEDCqlLogParserListener) ExitDisjQuery(ctx *DisjQueryContext) {}

// EnterConjQuery is called when production conjQuery is entered.
func (s *BaseEDCqlLogParserListener) EnterConjQuery(ctx *ConjQueryContext) {}

// ExitConjQuery is called when production conjQuery is exited.
func (s *BaseEDCqlLogParserListener) ExitConjQuery(ctx *ConjQueryContext) {}

// EnterModClause is called when production modClause is entered.
func (s *BaseEDCqlLogParserListener) EnterModClause(ctx *ModClauseContext) {}

// ExitModClause is called when production modClause is exited.
func (s *BaseEDCqlLogParserListener) ExitModClause(ctx *ModClauseContext) {}

// EnterModifier is called when production modifier is entered.
func (s *BaseEDCqlLogParserListener) EnterModifier(ctx *ModifierContext) {}

// ExitModifier is called when production modifier is exited.
func (s *BaseEDCqlLogParserListener) ExitModifier(ctx *ModifierContext) {}

// EnterClause is called when production clause is entered.
func (s *BaseEDCqlLogParserListener) EnterClause(ctx *ClauseContext) {}

// ExitClause is called when production clause is exited.
func (s *BaseEDCqlLogParserListener) ExitClause(ctx *ClauseContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseEDCqlLogParserListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseEDCqlLogParserListener) ExitTerm(ctx *TermContext) {}

// EnterGroupingExpr is called when production groupingExpr is entered.
func (s *BaseEDCqlLogParserListener) EnterGroupingExpr(ctx *GroupingExprContext) {}

// ExitGroupingExpr is called when production groupingExpr is exited.
func (s *BaseEDCqlLogParserListener) ExitGroupingExpr(ctx *GroupingExprContext) {}

// EnterFieldName is called when production fieldName is entered.
func (s *BaseEDCqlLogParserListener) EnterFieldName(ctx *FieldNameContext) {}

// ExitFieldName is called when production fieldName is exited.
func (s *BaseEDCqlLogParserListener) ExitFieldName(ctx *FieldNameContext) {}

// EnterQuotedTerm is called when production quotedTerm is entered.
func (s *BaseEDCqlLogParserListener) EnterQuotedTerm(ctx *QuotedTermContext) {}

// ExitQuotedTerm is called when production quotedTerm is exited.
func (s *BaseEDCqlLogParserListener) ExitQuotedTerm(ctx *QuotedTermContext) {}

// EnterOperatorColon is called when production operatorColon is entered.
func (s *BaseEDCqlLogParserListener) EnterOperatorColon(ctx *OperatorColonContext) {}

// ExitOperatorColon is called when production operatorColon is exited.
func (s *BaseEDCqlLogParserListener) ExitOperatorColon(ctx *OperatorColonContext) {}
