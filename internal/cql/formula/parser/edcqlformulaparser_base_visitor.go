// Code generated from EDCqlFormulaParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // EDCqlFormulaParser

import "github.com/antlr4-go/antlr/v4"

type BaseEDCqlFormulaParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseEDCqlFormulaParserVisitor) VisitFormula(ctx *FormulaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitFullExpression(ctx *FullExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitPart(ctx *PartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitNestablePart(ctx *NestablePartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitNestableGroup(ctx *NestableGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitOuterFunction(ctx *OuterFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitNestableFunction(ctx *NestableFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitTimeshift(ctx *TimeshiftContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitMovingAverage(ctx *MovingAverageContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitGroup(ctx *GroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitRawGroup(ctx *RawGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitRawPart(ctx *RawPartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitRawExpr(ctx *RawExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitMovingAverageExpr(ctx *MovingAverageExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitShift(ctx *ShiftContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseEDCqlFormulaParserVisitor) VisitMovingAveragePeriod(ctx *MovingAveragePeriodContext) interface{} {
	return v.VisitChildren(ctx)
}
