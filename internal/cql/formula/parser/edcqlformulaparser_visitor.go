// Code generated from EDCqlFormulaParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlFormulaParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by EDCqlFormulaParser.
type EDCqlFormulaParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by EDCqlFormulaParser#formula.
	VisitFormula(ctx *FormulaContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#fullExpression.
	VisitFullExpression(ctx *FullExpressionContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#part.
	VisitPart(ctx *PartContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#nestablePart.
	VisitNestablePart(ctx *NestablePartContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#nestableGroup.
	VisitNestableGroup(ctx *NestableGroupContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#outerFunction.
	VisitOuterFunction(ctx *OuterFunctionContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#nestableFunction.
	VisitNestableFunction(ctx *NestableFunctionContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#timeshift.
	VisitTimeshift(ctx *TimeshiftContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#movingAverage.
	VisitMovingAverage(ctx *MovingAverageContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#group.
	VisitGroup(ctx *GroupContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#rawGroup.
	VisitRawGroup(ctx *RawGroupContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#rawPart.
	VisitRawPart(ctx *RawPartContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#rawExpr.
	VisitRawExpr(ctx *RawExprContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#movingAverageExpr.
	VisitMovingAverageExpr(ctx *MovingAverageExprContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#shift.
	VisitShift(ctx *ShiftContext) interface{}

	// Visit a parse tree produced by EDCqlFormulaParser#movingAveragePeriod.
	VisitMovingAveragePeriod(ctx *MovingAveragePeriodContext) interface{}
}
