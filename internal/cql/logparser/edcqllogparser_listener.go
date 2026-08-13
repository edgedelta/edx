// Code generated from EDCqlLogParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package logparser // EDCqlLogParser
import "github.com/antlr4-go/antlr/v4"

// EDCqlLogParserListener is a complete listener for a parse tree produced by EDCqlLogParser.
type EDCqlLogParserListener interface {
	antlr.ParseTreeListener

	// EnterFinalQuery is called when entering the finalQuery production.
	EnterFinalQuery(c *FinalQueryContext)

	// EnterTopLevelQuery is called when entering the topLevelQuery production.
	EnterTopLevelQuery(c *TopLevelQueryContext)

	// EnterRollupMethod is called when entering the rollupMethod production.
	EnterRollupMethod(c *RollupMethodContext)

	// EnterRollupField is called when entering the rollupField production.
	EnterRollupField(c *RollupFieldContext)

	// EnterRollupSection is called when entering the rollupSection production.
	EnterRollupSection(c *RollupSectionContext)

	// EnterCountUniqueFields is called when entering the countUniqueFields production.
	EnterCountUniqueFields(c *CountUniqueFieldsContext)

	// EnterGroupByFields is called when entering the groupByFields production.
	EnterGroupByFields(c *GroupByFieldsContext)

	// EnterGroupBySection is called when entering the groupBySection production.
	EnterGroupBySection(c *GroupBySectionContext)

	// EnterCountUniqueSection is called when entering the countUniqueSection production.
	EnterCountUniqueSection(c *CountUniqueSectionContext)

	// EnterQueryOperation is called when entering the queryOperation production.
	EnterQueryOperation(c *QueryOperationContext)

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterDisjQuery is called when entering the disjQuery production.
	EnterDisjQuery(c *DisjQueryContext)

	// EnterConjQuery is called when entering the conjQuery production.
	EnterConjQuery(c *ConjQueryContext)

	// EnterModClause is called when entering the modClause production.
	EnterModClause(c *ModClauseContext)

	// EnterModifier is called when entering the modifier production.
	EnterModifier(c *ModifierContext)

	// EnterClause is called when entering the clause production.
	EnterClause(c *ClauseContext)

	// EnterTerm is called when entering the term production.
	EnterTerm(c *TermContext)

	// EnterGroupingExpr is called when entering the groupingExpr production.
	EnterGroupingExpr(c *GroupingExprContext)

	// EnterFieldName is called when entering the fieldName production.
	EnterFieldName(c *FieldNameContext)

	// EnterQuotedTerm is called when entering the quotedTerm production.
	EnterQuotedTerm(c *QuotedTermContext)

	// EnterOperatorColon is called when entering the operatorColon production.
	EnterOperatorColon(c *OperatorColonContext)

	// ExitFinalQuery is called when exiting the finalQuery production.
	ExitFinalQuery(c *FinalQueryContext)

	// ExitTopLevelQuery is called when exiting the topLevelQuery production.
	ExitTopLevelQuery(c *TopLevelQueryContext)

	// ExitRollupMethod is called when exiting the rollupMethod production.
	ExitRollupMethod(c *RollupMethodContext)

	// ExitRollupField is called when exiting the rollupField production.
	ExitRollupField(c *RollupFieldContext)

	// ExitRollupSection is called when exiting the rollupSection production.
	ExitRollupSection(c *RollupSectionContext)

	// ExitCountUniqueFields is called when exiting the countUniqueFields production.
	ExitCountUniqueFields(c *CountUniqueFieldsContext)

	// ExitGroupByFields is called when exiting the groupByFields production.
	ExitGroupByFields(c *GroupByFieldsContext)

	// ExitGroupBySection is called when exiting the groupBySection production.
	ExitGroupBySection(c *GroupBySectionContext)

	// ExitCountUniqueSection is called when exiting the countUniqueSection production.
	ExitCountUniqueSection(c *CountUniqueSectionContext)

	// ExitQueryOperation is called when exiting the queryOperation production.
	ExitQueryOperation(c *QueryOperationContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitDisjQuery is called when exiting the disjQuery production.
	ExitDisjQuery(c *DisjQueryContext)

	// ExitConjQuery is called when exiting the conjQuery production.
	ExitConjQuery(c *ConjQueryContext)

	// ExitModClause is called when exiting the modClause production.
	ExitModClause(c *ModClauseContext)

	// ExitModifier is called when exiting the modifier production.
	ExitModifier(c *ModifierContext)

	// ExitClause is called when exiting the clause production.
	ExitClause(c *ClauseContext)

	// ExitTerm is called when exiting the term production.
	ExitTerm(c *TermContext)

	// ExitGroupingExpr is called when exiting the groupingExpr production.
	ExitGroupingExpr(c *GroupingExprContext)

	// ExitFieldName is called when exiting the fieldName production.
	ExitFieldName(c *FieldNameContext)

	// ExitQuotedTerm is called when exiting the quotedTerm production.
	ExitQuotedTerm(c *QuotedTermContext)

	// ExitOperatorColon is called when exiting the operatorColon production.
	ExitOperatorColon(c *OperatorColonContext)
}
