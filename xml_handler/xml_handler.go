package xmlhandler

import "encoding/xml"

type XMLElementList struct {
	XMLName  xml.Name               `xml:"elements"`
	Elements []XMLElementDefinition `xml:"element"`
}

type XMLElementDefinition struct {
	XMLName  xml.Name         `xml:"element"`
	Name     string           `xml:"name,attr"`
	Role     string           `xml:"role,attr"`
	Display  *XMLDisplay      `xml:"display"`
	Material *XMLMaterialData `xml:"material"`

	Air            *XMLAirData
	ImmovableSolid *XMLImmovableSolidData
	MovableSolid   *XMLMovableSolidData
	Liquid         *XMLLiquidData
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

type XMLMaterialData struct {
	XMLName xml.Name `xml:"material"`
	Density float32  `xml:"density"`
}
