package xmlhandler

import "encoding/xml"

type XMLElementList struct {
	XMLName  xml.Name               `xml:"elements"`
	Elements []XMLElementDefinition `xml:"element"`
}

type XMLElementDefinition struct {
	XMLName   xml.Name         `xml:"element"`
	Name      string           `xml:"name,attr"`
	Role      string           `xml:"role,attr"`
	Display   *XMLDisplay      `xml:"display"`
	Material  *XMLMaterialData `xml:"material"`
	Reactions *XMLReactions    `xml:"reactions"`

	Air            *XMLAirData
	ImmovableSolid *XMLImmovableSolidData
	MovableSolid   *XMLMovableSolidData
	Liquid         *XMLLiquidData
	Gas            *XMLGasData
	Dust           *XMLDustData
}

type XMLDisplay struct {
	XMLName    xml.Name `xml:"display"`
	Name       string   `xml:"name"`
	Color      string   `xml:"color"`
	Selectable bool     `xml:"selectable"`
}

type XMLAirData struct {
	XMLName xml.Name `xml:"air"`
}

type XMLImmovableSolidData struct {
	XMLName xml.Name `xml:"immovable-solid"`
}

type XMLMovableSolidData struct {
	XMLName xml.Name `xml:"movable-solid"`
}

type XMLLiquidData struct {
	XMLName xml.Name `xml:"liquid"`
}

type XMLGasData struct {
	XMLName xml.Name `xml:"gas"`
	Weight  float32  `xml:"weight"`
}

type XMLDustData struct {
	XMLName xml.Name `xml:"dust"`
	Weight  float32  `xml:"weight"`
}

type XMLMaterialData struct {
	XMLName xml.Name `xml:"material"`
	Density float32  `xml:"density"`
}

type XMLReactions struct {
	XMLName   xml.Name      `xml:"reactions"`
	Reactions []XMLReaction `xml:"reaction"`
}

type XMLReaction struct {
	XMLName xml.Name       `xml:"reaction"`
	Steps   []ReactionStep `xml:",any"`
}

type ReactionStep struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}
