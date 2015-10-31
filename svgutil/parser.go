// http://play.golang.org/p/kyfff6Kg1c
package svgutil

import (
	"encoding/xml"
	"strconv"
)

type Path struct {
	Id string `xml:"id,attr"`
	D  string `xml:"d, attr"`
}

type Rect struct {
	Id string `xml:"id,attr"`
}

type Group struct {
	Id          string
	Stroke      string
	StrokeWidth int32
	Fill        string
	FillRule    string
	Elements    []interface{}
}

type Svg struct {
	Title  string  `xml:"title"`
	Groups []Group `xml:"g"`
}

// Implements encoding.xml.Unmarshaler interface
func (g *Group) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			g.Id = attr.Value
		case "stroke":
			g.Stroke = attr.Value
		case "stroke-width":
			if intValue, err := strconv.ParseInt(attr.Value, 10, 32); err != nil {
				return err
			} else {
				g.StrokeWidth = int32(intValue)
			}
		case "fill":
			g.Fill = attr.Value
		case "fill-rule":
			g.FillRule = attr.Value
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var elementStruct interface{}

			switch tok.Name.Local {
			case "rect":
				elementStruct = &Rect{}
			case "path":
				elementStruct = &Path{}
			}

			if err = decoder.DecodeElement(elementStruct, &tok); err != nil {
				return err
			} else {
				g.Elements = append(g.Elements, elementStruct)
			}

		case xml.EndElement:
			return nil
		}
	}
}

func ParseSvgString(str string) (*Svg, error) {
	svg := &Svg{}

	err := xml.Unmarshal([]byte(str), &svg)
	if err != nil {
		return nil, err
	}
	return svg, nil
}

func ParseSvgBytes(data []byte) (*Svg, error) {
	svg := &Svg{}

	err := xml.Unmarshal(data, &svg)
	if err != nil {
		return nil, err
	}
	return svg, nil
}
