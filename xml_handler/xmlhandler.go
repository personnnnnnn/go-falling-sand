package xmlhandler

import "encoding/xml"

type XMLElementList struct {
	XMLName  xml.Name               `xml:"elements"`
	Elements []XMLElementDefinition `xml:"element"`
}

type XMLElementDefinition struct {
	XMLName   xml.Name    `xml:"element"`
	Name      string      `xml:"name,attr"`
	IsDefault bool        `xml:"default,attr"`
	Display   *XMLDisplay `xml:"display"`
}

type XMLDisplay struct {
	XMLName xml.Name `xml:"display"`
	Name    string   `xml:"name"`
	Color   string   `xml:"color"`
}
