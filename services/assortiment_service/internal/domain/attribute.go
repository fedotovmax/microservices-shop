package domain

type AttributeDataType uint8

const (
	AttributeDataTypeString AttributeDataType = iota + 1
	AttributeDataTypeNumber
	AttributeDataTypeBoolean
)

type AttributeValueMetaType uint8

const (
	AttributeValueMetaTypeString AttributeValueMetaType = iota + 1
	AttributeValueMetaTypeNumber
	AttributeValueMetaTypeColor
)
