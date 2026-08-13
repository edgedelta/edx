// Code generated from EDCqlMetricParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlMetricParser

import "github.com/antlr4-go/antlr/v4"

type BaseEDCqlMetricParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseEDCqlMetricParserVisitor) VisitFinalQuery(ctx *FinalQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitTopLevelQuery(ctx *TopLevelQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitQueryOperation(ctx *QueryOperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitMetricName(ctx *MetricNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitAggregationFilter(ctx *AggregationFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitAggregationMethod(ctx *AggregationMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitGroupByFields(ctx *GroupByFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitCountUniqueFields(ctx *CountUniqueFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitGroupBySection(ctx *GroupBySectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitCountUniqueSection(ctx *CountUniqueSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitRollupWindow(ctx *RollupWindowContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitRollupSection(ctx *RollupSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitAggregation(ctx *AggregationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitFillMethod(ctx *FillMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitFillLimit(ctx *FillLimitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitFillSection(ctx *FillSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitDisjQuery(ctx *DisjQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitConjQuery(ctx *ConjQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitModClause(ctx *ModClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitModifier(ctx *ModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitClause(ctx *ClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitTerm(ctx *TermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitGroupingExpr(ctx *GroupingExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitFieldName(ctx *FieldNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitQuotedTerm(ctx *QuotedTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlMetricParserVisitor) VisitOperatorColon(ctx *OperatorColonContext) interface{} {
	return v.VisitChildren(ctx)
}
