package layout

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"
)

type ILayout = domain.ILayout
type IObject = domain.IObject

type PositionData = domain.PositionData
type BoxData = domain.BoxData

var (
	Position = enginetype.GetComponentType("pos")
	Box      = enginetype.GetComponentType("box")
)

const (
	DirRow        = domain.DirRow
	DirColumn     = domain.DirColumn
	DirMask       = domain.DirMask
	AlignStart    = domain.AlignStart
	AlignCenter   = domain.AlignCenter
	AlignEnd      = domain.AlignEnd
	AlignMask     = domain.AlignMask
	JustifyStart  = domain.JustifyStart
	JustifyCenter = domain.JustifyCenter
	JustifyEnd    = domain.JustifyEnd
	JustifyMask   = domain.JustifyMask
)
