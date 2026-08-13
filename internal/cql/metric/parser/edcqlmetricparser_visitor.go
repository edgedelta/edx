// Code generated from EDCqlMetricParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlMetricParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by EDCqlMetricParser.
type EDCqlMetricParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by EDCqlMetricParser#finalQuery.
	VisitFinalQuery(ctx *FinalQueryContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#topLevelQuery.
	VisitTopLevelQuery(ctx *TopLevelQueryContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#queryOperation.
	VisitQueryOperation(ctx *QueryOperationContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#metricName.
	VisitMetricName(ctx *MetricNameContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#aggregationFilter.
	VisitAggregationFilter(ctx *AggregationFilterContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#aggregationMethod.
	VisitAggregationMethod(ctx *AggregationMethodContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#groupByFields.
	VisitGroupByFields(ctx *GroupByFieldsContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#countUniqueFields.
	VisitCountUniqueFields(ctx *CountUniqueFieldsContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#groupBySection.
	VisitGroupBySection(ctx *GroupBySectionContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#countUniqueSection.
	VisitCountUniqueSection(ctx *CountUniqueSectionContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#rollupWindow.
	VisitRollupWindow(ctx *RollupWindowContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#rollupSection.
	VisitRollupSection(ctx *RollupSectionContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#aggregation.
	VisitAggregation(ctx *AggregationContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#fillMethod.
	VisitFillMethod(ctx *FillMethodContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#fillLimit.
	VisitFillLimit(ctx *FillLimitContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#fillSection.
	VisitFillSection(ctx *FillSectionContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#query.
	VisitQuery(ctx *QueryContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#disjQuery.
	VisitDisjQuery(ctx *DisjQueryContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#conjQuery.
	VisitConjQuery(ctx *ConjQueryContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#modClause.
	VisitModClause(ctx *ModClauseContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#modifier.
	VisitModifier(ctx *ModifierContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#clause.
	VisitClause(ctx *ClauseContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#term.
	VisitTerm(ctx *TermContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#groupingExpr.
	VisitGroupingExpr(ctx *GroupingExprContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#fieldName.
	VisitFieldName(ctx *FieldNameContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#quotedTerm.
	VisitQuotedTerm(ctx *QuotedTermContext) interface{}

	// Visit a parse tree produced by EDCqlMetricParser#operatorColon.
	VisitOperatorColon(ctx *OperatorColonContext) interface{}
}
