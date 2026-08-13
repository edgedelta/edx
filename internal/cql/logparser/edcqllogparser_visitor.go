// Code generated from EDCqlLogParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package logparser // EDCqlLogParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by EDCqlLogParser.
type EDCqlLogParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by EDCqlLogParser#finalQuery.
	VisitFinalQuery(ctx *FinalQueryContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#topLevelQuery.
	VisitTopLevelQuery(ctx *TopLevelQueryContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#rollupMethod.
	VisitRollupMethod(ctx *RollupMethodContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#rollupField.
	VisitRollupField(ctx *RollupFieldContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#rollupSection.
	VisitRollupSection(ctx *RollupSectionContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#countUniqueFields.
	VisitCountUniqueFields(ctx *CountUniqueFieldsContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#groupByFields.
	VisitGroupByFields(ctx *GroupByFieldsContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#groupBySection.
	VisitGroupBySection(ctx *GroupBySectionContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#countUniqueSection.
	VisitCountUniqueSection(ctx *CountUniqueSectionContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#queryOperation.
	VisitQueryOperation(ctx *QueryOperationContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#query.
	VisitQuery(ctx *QueryContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#disjQuery.
	VisitDisjQuery(ctx *DisjQueryContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#conjQuery.
	VisitConjQuery(ctx *ConjQueryContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#modClause.
	VisitModClause(ctx *ModClauseContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#modifier.
	VisitModifier(ctx *ModifierContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#clause.
	VisitClause(ctx *ClauseContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#term.
	VisitTerm(ctx *TermContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#groupingExpr.
	VisitGroupingExpr(ctx *GroupingExprContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#fieldName.
	VisitFieldName(ctx *FieldNameContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#quotedTerm.
	VisitQuotedTerm(ctx *QuotedTermContext) interface{}

	// Visit a parse tree produced by EDCqlLogParser#operatorColon.
	VisitOperatorColon(ctx *OperatorColonContext) interface{}
}
