// Code generated from EDCqlMetricParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlMetricParser

import "github.com/antlr4-go/antlr/v4"

// BaseEDCqlMetricParserListener is a complete listener for a parse tree produced by EDCqlMetricParser.
type BaseEDCqlMetricParserListener struct{}

var _ EDCqlMetricParserListener = &BaseEDCqlMetricParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseEDCqlMetricParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseEDCqlMetricParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseEDCqlMetricParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseEDCqlMetricParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFinalQuery is called when production finalQuery is entered.
func (s *BaseEDCqlMetricParserListener) EnterFinalQuery(ctx *FinalQueryContext) {}

// ExitFinalQuery is called when production finalQuery is exited.
func (s *BaseEDCqlMetricParserListener) ExitFinalQuery(ctx *FinalQueryContext) {}

// EnterTopLevelQuery is called when production topLevelQuery is entered.
func (s *BaseEDCqlMetricParserListener) EnterTopLevelQuery(ctx *TopLevelQueryContext) {}

// ExitTopLevelQuery is called when production topLevelQuery is exited.
func (s *BaseEDCqlMetricParserListener) ExitTopLevelQuery(ctx *TopLevelQueryContext) {}

// EnterQueryOperation is called when production queryOperation is entered.
func (s *BaseEDCqlMetricParserListener) EnterQueryOperation(ctx *QueryOperationContext) {}

// ExitQueryOperation is called when production queryOperation is exited.
func (s *BaseEDCqlMetricParserListener) ExitQueryOperation(ctx *QueryOperationContext) {}

// EnterMetricName is called when production metricName is entered.
func (s *BaseEDCqlMetricParserListener) EnterMetricName(ctx *MetricNameContext) {}

// ExitMetricName is called when production metricName is exited.
func (s *BaseEDCqlMetricParserListener) ExitMetricName(ctx *MetricNameContext) {}

// EnterAggregationFilter is called when production aggregationFilter is entered.
func (s *BaseEDCqlMetricParserListener) EnterAggregationFilter(ctx *AggregationFilterContext) {}

// ExitAggregationFilter is called when production aggregationFilter is exited.
func (s *BaseEDCqlMetricParserListener) ExitAggregationFilter(ctx *AggregationFilterContext) {}

// EnterAggregationMethod is called when production aggregationMethod is entered.
func (s *BaseEDCqlMetricParserListener) EnterAggregationMethod(ctx *AggregationMethodContext) {}

// ExitAggregationMethod is called when production aggregationMethod is exited.
func (s *BaseEDCqlMetricParserListener) ExitAggregationMethod(ctx *AggregationMethodContext) {}

// EnterGroupByFields is called when production groupByFields is entered.
func (s *BaseEDCqlMetricParserListener) EnterGroupByFields(ctx *GroupByFieldsContext) {}

// ExitGroupByFields is called when production groupByFields is exited.
func (s *BaseEDCqlMetricParserListener) ExitGroupByFields(ctx *GroupByFieldsContext) {}

// EnterCountUniqueFields is called when production countUniqueFields is entered.
func (s *BaseEDCqlMetricParserListener) EnterCountUniqueFields(ctx *CountUniqueFieldsContext) {}

// ExitCountUniqueFields is called when production countUniqueFields is exited.
func (s *BaseEDCqlMetricParserListener) ExitCountUniqueFields(ctx *CountUniqueFieldsContext) {}

// EnterGroupBySection is called when production groupBySection is entered.
func (s *BaseEDCqlMetricParserListener) EnterGroupBySection(ctx *GroupBySectionContext) {}

// ExitGroupBySection is called when production groupBySection is exited.
func (s *BaseEDCqlMetricParserListener) ExitGroupBySection(ctx *GroupBySectionContext) {}

// EnterCountUniqueSection is called when production countUniqueSection is entered.
func (s *BaseEDCqlMetricParserListener) EnterCountUniqueSection(ctx *CountUniqueSectionContext) {}

// ExitCountUniqueSection is called when production countUniqueSection is exited.
func (s *BaseEDCqlMetricParserListener) ExitCountUniqueSection(ctx *CountUniqueSectionContext) {}

// EnterRollupWindow is called when production rollupWindow is entered.
func (s *BaseEDCqlMetricParserListener) EnterRollupWindow(ctx *RollupWindowContext) {}

// ExitRollupWindow is called when production rollupWindow is exited.
func (s *BaseEDCqlMetricParserListener) ExitRollupWindow(ctx *RollupWindowContext) {}

// EnterRollupSection is called when production rollupSection is entered.
func (s *BaseEDCqlMetricParserListener) EnterRollupSection(ctx *RollupSectionContext) {}

// ExitRollupSection is called when production rollupSection is exited.
func (s *BaseEDCqlMetricParserListener) ExitRollupSection(ctx *RollupSectionContext) {}

// EnterAggregation is called when production aggregation is entered.
func (s *BaseEDCqlMetricParserListener) EnterAggregation(ctx *AggregationContext) {}

// ExitAggregation is called when production aggregation is exited.
func (s *BaseEDCqlMetricParserListener) ExitAggregation(ctx *AggregationContext) {}

// EnterFillMethod is called when production fillMethod is entered.
func (s *BaseEDCqlMetricParserListener) EnterFillMethod(ctx *FillMethodContext) {}

// ExitFillMethod is called when production fillMethod is exited.
func (s *BaseEDCqlMetricParserListener) ExitFillMethod(ctx *FillMethodContext) {}

// EnterFillLimit is called when production fillLimit is entered.
func (s *BaseEDCqlMetricParserListener) EnterFillLimit(ctx *FillLimitContext) {}

// ExitFillLimit is called when production fillLimit is exited.
func (s *BaseEDCqlMetricParserListener) ExitFillLimit(ctx *FillLimitContext) {}

// EnterFillSection is called when production fillSection is entered.
func (s *BaseEDCqlMetricParserListener) EnterFillSection(ctx *FillSectionContext) {}

// ExitFillSection is called when production fillSection is exited.
func (s *BaseEDCqlMetricParserListener) ExitFillSection(ctx *FillSectionContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseEDCqlMetricParserListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseEDCqlMetricParserListener) ExitQuery(ctx *QueryContext) {}

// EnterDisjQuery is called when production disjQuery is entered.
func (s *BaseEDCqlMetricParserListener) EnterDisjQuery(ctx *DisjQueryContext) {}

// ExitDisjQuery is called when production disjQuery is exited.
func (s *BaseEDCqlMetricParserListener) ExitDisjQuery(ctx *DisjQueryContext) {}

// EnterConjQuery is called when production conjQuery is entered.
func (s *BaseEDCqlMetricParserListener) EnterConjQuery(ctx *ConjQueryContext) {}

// ExitConjQuery is called when production conjQuery is exited.
func (s *BaseEDCqlMetricParserListener) ExitConjQuery(ctx *ConjQueryContext) {}

// EnterModClause is called when production modClause is entered.
func (s *BaseEDCqlMetricParserListener) EnterModClause(ctx *ModClauseContext) {}

// ExitModClause is called when production modClause is exited.
func (s *BaseEDCqlMetricParserListener) ExitModClause(ctx *ModClauseContext) {}

// EnterModifier is called when production modifier is entered.
func (s *BaseEDCqlMetricParserListener) EnterModifier(ctx *ModifierContext) {}

// ExitModifier is called when production modifier is exited.
func (s *BaseEDCqlMetricParserListener) ExitModifier(ctx *ModifierContext) {}

// EnterClause is called when production clause is entered.
func (s *BaseEDCqlMetricParserListener) EnterClause(ctx *ClauseContext) {}

// ExitClause is called when production clause is exited.
func (s *BaseEDCqlMetricParserListener) ExitClause(ctx *ClauseContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseEDCqlMetricParserListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseEDCqlMetricParserListener) ExitTerm(ctx *TermContext) {}

// EnterGroupingExpr is called when production groupingExpr is entered.
func (s *BaseEDCqlMetricParserListener) EnterGroupingExpr(ctx *GroupingExprContext) {}

// ExitGroupingExpr is called when production groupingExpr is exited.
func (s *BaseEDCqlMetricParserListener) ExitGroupingExpr(ctx *GroupingExprContext) {}

// EnterFieldName is called when production fieldName is entered.
func (s *BaseEDCqlMetricParserListener) EnterFieldName(ctx *FieldNameContext) {}

// ExitFieldName is called when production fieldName is exited.
func (s *BaseEDCqlMetricParserListener) ExitFieldName(ctx *FieldNameContext) {}

// EnterQuotedTerm is called when production quotedTerm is entered.
func (s *BaseEDCqlMetricParserListener) EnterQuotedTerm(ctx *QuotedTermContext) {}

// ExitQuotedTerm is called when production quotedTerm is exited.
func (s *BaseEDCqlMetricParserListener) ExitQuotedTerm(ctx *QuotedTermContext) {}

// EnterOperatorColon is called when production operatorColon is entered.
func (s *BaseEDCqlMetricParserListener) EnterOperatorColon(ctx *OperatorColonContext) {}

// ExitOperatorColon is called when production operatorColon is exited.
func (s *BaseEDCqlMetricParserListener) ExitOperatorColon(ctx *OperatorColonContext) {}
