// Code generated from EDCqlFormulaParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlFormulaParser

import "github.com/antlr4-go/antlr/v4"

// EDCqlFormulaParserListener is a complete listener for a parse tree produced by EDCqlFormulaParser.
type EDCqlFormulaParserListener interface {
	antlr.ParseTreeListener

	// EnterFormula is called when entering the formula production.
	EnterFormula(c *FormulaContext)

	// EnterFullExpression is called when entering the fullExpression production.
	EnterFullExpression(c *FullExpressionContext)

	// EnterPart is called when entering the part production.
	EnterPart(c *PartContext)

	// EnterNestablePart is called when entering the nestablePart production.
	EnterNestablePart(c *NestablePartContext)

	// EnterNestableGroup is called when entering the nestableGroup production.
	EnterNestableGroup(c *NestableGroupContext)

	// EnterOuterFunction is called when entering the outerFunction production.
	EnterOuterFunction(c *OuterFunctionContext)

	// EnterNestableFunction is called when entering the nestableFunction production.
	EnterNestableFunction(c *NestableFunctionContext)

	// EnterTimeshift is called when entering the timeshift production.
	EnterTimeshift(c *TimeshiftContext)

	// EnterMovingAverage is called when entering the movingAverage production.
	EnterMovingAverage(c *MovingAverageContext)

	// EnterGroup is called when entering the group production.
	EnterGroup(c *GroupContext)

	// EnterRawGroup is called when entering the rawGroup production.
	EnterRawGroup(c *RawGroupContext)

	// EnterRawPart is called when entering the rawPart production.
	EnterRawPart(c *RawPartContext)

	// EnterRawExpr is called when entering the rawExpr production.
	EnterRawExpr(c *RawExprContext)

	// EnterMovingAverageExpr is called when entering the movingAverageExpr production.
	EnterMovingAverageExpr(c *MovingAverageExprContext)

	// EnterShift is called when entering the shift production.
	EnterShift(c *ShiftContext)

	// EnterMovingAveragePeriod is called when entering the movingAveragePeriod production.
	EnterMovingAveragePeriod(c *MovingAveragePeriodContext)

	// ExitFormula is called when exiting the formula production.
	ExitFormula(c *FormulaContext)

	// ExitFullExpression is called when exiting the fullExpression production.
	ExitFullExpression(c *FullExpressionContext)

	// ExitPart is called when exiting the part production.
	ExitPart(c *PartContext)

	// ExitNestablePart is called when exiting the nestablePart production.
	ExitNestablePart(c *NestablePartContext)

	// ExitNestableGroup is called when exiting the nestableGroup production.
	ExitNestableGroup(c *NestableGroupContext)

	// ExitOuterFunction is called when exiting the outerFunction production.
	ExitOuterFunction(c *OuterFunctionContext)

	// ExitNestableFunction is called when exiting the nestableFunction production.
	ExitNestableFunction(c *NestableFunctionContext)

	// ExitTimeshift is called when exiting the timeshift production.
	ExitTimeshift(c *TimeshiftContext)

	// ExitMovingAverage is called when exiting the movingAverage production.
	ExitMovingAverage(c *MovingAverageContext)

	// ExitGroup is called when exiting the group production.
	ExitGroup(c *GroupContext)

	// ExitRawGroup is called when exiting the rawGroup production.
	ExitRawGroup(c *RawGroupContext)

	// ExitRawPart is called when exiting the rawPart production.
	ExitRawPart(c *RawPartContext)

	// ExitRawExpr is called when exiting the rawExpr production.
	ExitRawExpr(c *RawExprContext)

	// ExitMovingAverageExpr is called when exiting the movingAverageExpr production.
	ExitMovingAverageExpr(c *MovingAverageExprContext)

	// ExitShift is called when exiting the shift production.
	ExitShift(c *ShiftContext)

	// ExitMovingAveragePeriod is called when exiting the movingAveragePeriod production.
	ExitMovingAveragePeriod(c *MovingAveragePeriodContext)
}
