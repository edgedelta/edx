// Code generated from EDCqlLogParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package logparser // EDCqlLogParser
import "github.com/antlr4-go/antlr/v4"

type BaseEDCqlLogParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseEDCqlLogParserVisitor) VisitFinalQuery(ctx *FinalQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitTopLevelQuery(ctx *TopLevelQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitRollupMethod(ctx *RollupMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitRollupField(ctx *RollupFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitRollupSection(ctx *RollupSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitCountUniqueFields(ctx *CountUniqueFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitGroupByFields(ctx *GroupByFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitGroupBySection(ctx *GroupBySectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitCountUniqueSection(ctx *CountUniqueSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitQueryOperation(ctx *QueryOperationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitDisjQuery(ctx *DisjQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitConjQuery(ctx *ConjQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitModClause(ctx *ModClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitModifier(ctx *ModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitClause(ctx *ClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitTerm(ctx *TermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitGroupingExpr(ctx *GroupingExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitFieldName(ctx *FieldNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitQuotedTerm(ctx *QuotedTermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlLogParserVisitor) VisitOperatorColon(ctx *OperatorColonContext) interface{} {
	return v.VisitChildren(ctx)
}
