// Code generated from EDCqlFormulaParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlFormulaParser

import "github.com/antlr4-go/antlr/v4"

// BaseEDCqlFormulaParserListener is a complete listener for a parse tree produced by EDCqlFormulaParser.
type BaseEDCqlFormulaParserListener struct{}

var _ EDCqlFormulaParserListener = &BaseEDCqlFormulaParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseEDCqlFormulaParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseEDCqlFormulaParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseEDCqlFormulaParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseEDCqlFormulaParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterFormula is called when production formula is entered.
func (s *BaseEDCqlFormulaParserListener) EnterFormula(ctx *FormulaContext) {}

// ExitFormula is called when production formula is exited.
func (s *BaseEDCqlFormulaParserListener) ExitFormula(ctx *FormulaContext) {}

// EnterFullExpression is called when production fullExpression is entered.
func (s *BaseEDCqlFormulaParserListener) EnterFullExpression(ctx *FullExpressionContext) {}

// ExitFullExpression is called when production fullExpression is exited.
func (s *BaseEDCqlFormulaParserListener) ExitFullExpression(ctx *FullExpressionContext) {}

// EnterPart is called when production part is entered.
func (s *BaseEDCqlFormulaParserListener) EnterPart(ctx *PartContext) {}

// ExitPart is called when production part is exited.
func (s *BaseEDCqlFormulaParserListener) ExitPart(ctx *PartContext) {}

// EnterNestablePart is called when production nestablePart is entered.
func (s *BaseEDCqlFormulaParserListener) EnterNestablePart(ctx *NestablePartContext) {}

// ExitNestablePart is called when production nestablePart is exited.
func (s *BaseEDCqlFormulaParserListener) ExitNestablePart(ctx *NestablePartContext) {}

// EnterNestableGroup is called when production nestableGroup is entered.
func (s *BaseEDCqlFormulaParserListener) EnterNestableGroup(ctx *NestableGroupContext) {}

// ExitNestableGroup is called when production nestableGroup is exited.
func (s *BaseEDCqlFormulaParserListener) ExitNestableGroup(ctx *NestableGroupContext) {}

// EnterOuterFunction is called when production outerFunction is entered.
func (s *BaseEDCqlFormulaParserListener) EnterOuterFunction(ctx *OuterFunctionContext) {}

// ExitOuterFunction is called when production outerFunction is exited.
func (s *BaseEDCqlFormulaParserListener) ExitOuterFunction(ctx *OuterFunctionContext) {}

// EnterNestableFunction is called when production nestableFunction is entered.
func (s *BaseEDCqlFormulaParserListener) EnterNestableFunction(ctx *NestableFunctionContext) {}

// ExitNestableFunction is called when production nestableFunction is exited.
func (s *BaseEDCqlFormulaParserListener) ExitNestableFunction(ctx *NestableFunctionContext) {}

// EnterTimeshift is called when production timeshift is entered.
func (s *BaseEDCqlFormulaParserListener) EnterTimeshift(ctx *TimeshiftContext) {}

// ExitTimeshift is called when production timeshift is exited.
func (s *BaseEDCqlFormulaParserListener) ExitTimeshift(ctx *TimeshiftContext) {}

// EnterMovingAverage is called when production movingAverage is entered.
func (s *BaseEDCqlFormulaParserListener) EnterMovingAverage(ctx *MovingAverageContext) {}

// ExitMovingAverage is called when production movingAverage is exited.
func (s *BaseEDCqlFormulaParserListener) ExitMovingAverage(ctx *MovingAverageContext) {}

// EnterGroup is called when production group is entered.
func (s *BaseEDCqlFormulaParserListener) EnterGroup(ctx *GroupContext) {}

// ExitGroup is called when production group is exited.
func (s *BaseEDCqlFormulaParserListener) ExitGroup(ctx *GroupContext) {}

// EnterRawGroup is called when production rawGroup is entered.
func (s *BaseEDCqlFormulaParserListener) EnterRawGroup(ctx *RawGroupContext) {}

// ExitRawGroup is called when production rawGroup is exited.
func (s *BaseEDCqlFormulaParserListener) ExitRawGroup(ctx *RawGroupContext) {}

// EnterRawPart is called when production rawPart is entered.
func (s *BaseEDCqlFormulaParserListener) EnterRawPart(ctx *RawPartContext) {}

// ExitRawPart is called when production rawPart is exited.
func (s *BaseEDCqlFormulaParserListener) ExitRawPart(ctx *RawPartContext) {}

// EnterRawExpr is called when production rawExpr is entered.
func (s *BaseEDCqlFormulaParserListener) EnterRawExpr(ctx *RawExprContext) {}

// ExitRawExpr is called when production rawExpr is exited.
func (s *BaseEDCqlFormulaParserListener) ExitRawExpr(ctx *RawExprContext) {}

// EnterMovingAverageExpr is called when production movingAverageExpr is entered.
func (s *BaseEDCqlFormulaParserListener) EnterMovingAverageExpr(ctx *MovingAverageExprContext) {}

// ExitMovingAverageExpr is called when production movingAverageExpr is exited.
func (s *BaseEDCqlFormulaParserListener) ExitMovingAverageExpr(ctx *MovingAverageExprContext) {}

// EnterShift is called when production shift is entered.
func (s *BaseEDCqlFormulaParserListener) EnterShift(ctx *ShiftContext) {}

// ExitShift is called when production shift is exited.
func (s *BaseEDCqlFormulaParserListener) ExitShift(ctx *ShiftContext) {}

// EnterMovingAveragePeriod is called when production movingAveragePeriod is entered.
func (s *BaseEDCqlFormulaParserListener) EnterMovingAveragePeriod(ctx *MovingAveragePeriodContext) {}

// ExitMovingAveragePeriod is called when production movingAveragePeriod is exited.
func (s *BaseEDCqlFormulaParserListener) ExitMovingAveragePeriod(ctx *MovingAveragePeriodContext) {}
