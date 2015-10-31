package svgutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parser(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	svgStr := `
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg width="79px" height="114px" viewBox="0 0 79 114" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
    <!-- Generator: Sketch 3.0.4 (8053) - http://www.bohemiancoding.com/sketch -->
    <title>ship</title>
    <desc>Created with Sketch.</desc>
    <defs></defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
        <path d="M70.7470703,54.9351921 C59.4438539,23.2101932 39.4404297,-0.0302734375 39.4404297,-0.0302734375 C39.4404297,-0.0302734375 19.8288957,22.9468825 8.09220641,54.9351916 C3.08063764,68.5942062 -0.495117188,83.8962169 -0.495117188,99.7539062 L19.5214844,108.566406 L59.1821013,108.566406 L79.4462891,100.046875 C79.4462891,100.046875 75.1865234,67.3955078 70.7470703,54.9351921 Z" id="Path-1" fill="#D8D8D8" sketch:type="MSShapeGroup"></path>
        <rect id="Rectangle-1" fill="#F6A623" sketch:type="MSShapeGroup" x="22" y="107" width="34" height="7"></rect>
    </g>
</svg>
`

	svg, err := ParseSvgString(svgStr)
	assert.NoError(err)
	assert.NotNil(svg)
	assert.Len(svg.Groups, 1)

	svg, err = ParseSvgBytes([]byte(svgStr))
	assert.NoError(err)
	assert.NotNil(svg)
	assert.Len(svg.Groups, 1)
}
