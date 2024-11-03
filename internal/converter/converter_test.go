package converter

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Converter_Convert(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		calls  []testRead
		source string
		want   Result
	}{
		{
			name:   "should convert markdown",
			path:   "/testdata/Markdown.md",
			source: "# Hello, World!\n",
			want:   Result{Markdown: "# Hello, World!\n"},
		},
		{
			name: "images",
			path: "/testdata/Images.md",
			calls: []testRead{
				{path: "/testdata/scream.png", content: "IMAGE CONTENT", err: nil},
			},
			source: "![Scream](./scream.png)\n",
			want: Result{
				Markdown: "![Scream](@media/22abb8f07c02970e.png)\n",
				Attachments: []Attachment{
					{Bytes: []byte("IMAGE CONTENT"), Filename: "22abb8f07c02970e.png"},
				},
			},
		},
		{
			name:   "video",
			path:   "/testdata/Video.md",
			source: "![@video](youtube.com/VIDEO)\n",
			want: Result{
				Markdown: "<iframe src=\"youtube.com/VIDEO?rel=0&amp;autoplay=0&amp;showinfo=0&amp;enablejsapi=0\" frameborder=\"0\" loading=\"lazy\" gesture=\"media\" allow=\"autoplay; fullscreen\" allowautoplay=\"true\" allowfullscreen=\"true\" style=\"aspect-ratio:16/9;height:100%;width:100%;\"></iframe>\n",
			},
		},
		{
			name:   "mermaid",
			source: "```mermaid\ngraph TD;\n    A-->B;\n    A-->C;\n    B-->D;\n    C-->D;\n```\n",
			want: Result{
				Markdown: "<div class=\"mermaid\"><svg aria-roledescription=\"flowchart-v2\" role=\"graphics-document document\" viewBox=\"0 0 204.640625 278\" style=\"max-width: 204.641px; background-color: white;\" class=\"flowchart\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" xmlns=\"http://www.w3.org/2000/svg\" width=\"100%\" id=\"my-svg\"><style>#my-svg{font-family:\"trebuchet ms\",verdana,arial,sans-serif;font-size:16px;fill:#333;}#my-svg .error-icon{fill:#552222;}#my-svg .error-text{fill:#552222;stroke:#552222;}#my-svg .edge-thickness-normal{stroke-width:1px;}#my-svg .edge-thickness-thick{stroke-width:3.5px;}#my-svg .edge-pattern-solid{stroke-dasharray:0;}#my-svg .edge-thickness-invisible{stroke-width:0;fill:none;}#my-svg .edge-pattern-dashed{stroke-dasharray:3;}#my-svg .edge-pattern-dotted{stroke-dasharray:2;}#my-svg .marker{fill:#333333;stroke:#333333;}#my-svg .marker.cross{stroke:#333333;}#my-svg svg{font-family:\"trebuchet ms\",verdana,arial,sans-serif;font-size:16px;}#my-svg p{margin:0;}#my-svg .label{font-family:\"trebuchet ms\",verdana,arial,sans-serif;color:#333;}#my-svg .cluster-label text{fill:#333;}#my-svg .cluster-label span{color:#333;}#my-svg .cluster-label span p{background-color:transparent;}#my-svg .label text,#my-svg span{fill:#333;color:#333;}#my-svg .node rect,#my-svg .node circle,#my-svg .node ellipse,#my-svg .node polygon,#my-svg .node path{fill:#ECECFF;stroke:#9370DB;stroke-width:1px;}#my-svg .rough-node .label text,#my-svg .node .label text,#my-svg .image-shape .label,#my-svg .icon-shape .label{text-anchor:middle;}#my-svg .node .katex path{fill:#000;stroke:#000;stroke-width:1px;}#my-svg .rough-node .label,#my-svg .node .label,#my-svg .image-shape .label,#my-svg .icon-shape .label{text-align:center;}#my-svg .node.clickable{cursor:pointer;}#my-svg .root .anchor path{fill:#333333!important;stroke-width:0;stroke:#333333;}#my-svg .arrowheadPath{fill:#333333;}#my-svg .edgePath .path{stroke:#333333;stroke-width:2.0px;}#my-svg .flowchart-link{stroke:#333333;fill:none;}#my-svg .edgeLabel{background-color:rgba(232,232,232, 0.8);text-align:center;}#my-svg .edgeLabel p{background-color:rgba(232,232,232, 0.8);}#my-svg .edgeLabel rect{opacity:0.5;background-color:rgba(232,232,232, 0.8);fill:rgba(232,232,232, 0.8);}#my-svg .labelBkg{background-color:rgba(232, 232, 232, 0.5);}#my-svg .cluster rect{fill:#ffffde;stroke:#aaaa33;stroke-width:1px;}#my-svg .cluster text{fill:#333;}#my-svg .cluster span{color:#333;}#my-svg div.mermaidTooltip{position:absolute;text-align:center;max-width:200px;padding:2px;font-family:\"trebuchet ms\",verdana,arial,sans-serif;font-size:12px;background:hsl(80, 100%, 96.2745098039%);border:1px solid #aaaa33;border-radius:2px;pointer-events:none;z-index:100;}#my-svg .flowchartTitleText{text-anchor:middle;font-size:18px;fill:#333;}#my-svg rect.text{fill:none;stroke-width:0;}#my-svg .icon-shape,#my-svg .image-shape{background-color:rgba(232,232,232, 0.8);text-align:center;}#my-svg .icon-shape p,#my-svg .image-shape p{background-color:rgba(232,232,232, 0.8);padding:2px;}#my-svg .icon-shape rect,#my-svg .image-shape rect{opacity:0.5;background-color:rgba(232,232,232, 0.8);fill:rgba(232,232,232, 0.8);}#my-svg :root{--mermaid-font-family:\"trebuchet ms\",verdana,arial,sans-serif;}</style><g><marker orient=\"auto\" markerHeight=\"8\" markerWidth=\"8\" markerUnits=\"userSpaceOnUse\" refY=\"5\" refX=\"5\" viewBox=\"0 0 10 10\" class=\"marker flowchart-v2\" id=\"my-svg_flowchart-v2-pointEnd\"><path style=\"stroke-width: 1; stroke-dasharray: 1, 0;\" class=\"arrowMarkerPath\" d=\"M 0 0 L 10 5 L 0 10 z\"/></marker><marker orient=\"auto\" markerHeight=\"8\" markerWidth=\"8\" markerUnits=\"userSpaceOnUse\" refY=\"5\" refX=\"4.5\" viewBox=\"0 0 10 10\" class=\"marker flowchart-v2\" id=\"my-svg_flowchart-v2-pointStart\"><path style=\"stroke-width: 1; stroke-dasharray: 1, 0;\" class=\"arrowMarkerPath\" d=\"M 0 5 L 10 10 L 10 0 z\"/></marker><marker orient=\"auto\" markerHeight=\"11\" markerWidth=\"11\" markerUnits=\"userSpaceOnUse\" refY=\"5\" refX=\"11\" viewBox=\"0 0 10 10\" class=\"marker flowchart-v2\" id=\"my-svg_flowchart-v2-circleEnd\"><circle style=\"stroke-width: 1; stroke-dasharray: 1, 0;\" class=\"arrowMarkerPath\" r=\"5\" cy=\"5\" cx=\"5\"/></marker><marker orient=\"auto\" markerHeight=\"11\" markerWidth=\"11\" markerUnits=\"userSpaceOnUse\" refY=\"5\" refX=\"-1\" viewBox=\"0 0 10 10\" class=\"marker flowchart-v2\" id=\"my-svg_flowchart-v2-circleStart\"><circle style=\"stroke-width: 1; stroke-dasharray: 1, 0;\" class=\"arrowMarkerPath\" r=\"5\" cy=\"5\" cx=\"5\"/></marker><marker orient=\"auto\" markerHeight=\"11\" markerWidth=\"11\" markerUnits=\"userSpaceOnUse\" refY=\"5.2\" refX=\"12\" viewBox=\"0 0 11 11\" class=\"marker cross flowchart-v2\" id=\"my-svg_flowchart-v2-crossEnd\"><path style=\"stroke-width: 2; stroke-dasharray: 1, 0;\" class=\"arrowMarkerPath\" d=\"M 1,1 l 9,9 M 10,1 l -9,9\"/></marker><marker orient=\"auto\" markerHeight=\"11\" markerWidth=\"11\" markerUnits=\"userSpaceOnUse\" refY=\"5.2\" refX=\"-1\" viewBox=\"0 0 11 11\" class=\"marker cross flowchart-v2\" id=\"my-svg_flowchart-v2-crossStart\"><path style=\"stroke-width: 2; stroke-dasharray: 1, 0;\" class=\"arrowMarkerPath\" d=\"M 1,1 l 9,9 M 10,1 l -9,9\"/></marker><g class=\"root\"><g class=\"clusters\"/><g class=\"edgePaths\"><path marker-end=\"url(#my-svg_flowchart-v2-pointEnd)\" style=\"\" class=\"edge-thickness-normal edge-pattern-solid edge-thickness-normal edge-pattern-solid flowchart-link\" id=\"L_A_B_0\" d=\"M71.214,62L66.434,66.167C61.653,70.333,52.092,78.667,47.312,86.333C42.531,94,42.531,101,42.531,104.5L42.531,108\"/><path marker-end=\"url(#my-svg_flowchart-v2-pointEnd)\" style=\"\" class=\"edge-thickness-normal edge-pattern-solid edge-thickness-normal edge-pattern-solid flowchart-link\" id=\"L_A_C_1\" d=\"M133.169,62L137.949,66.167C142.73,70.333,152.291,78.667,157.071,86.333C161.852,94,161.852,101,161.852,104.5L161.852,108\"/><path marker-end=\"url(#my-svg_flowchart-v2-pointEnd)\" style=\"\" class=\"edge-thickness-normal edge-pattern-solid edge-thickness-normal edge-pattern-solid flowchart-link\" id=\"L_B_D_2\" d=\"M42.531,166L42.531,170.167C42.531,174.333,42.531,182.667,46.809,190.562C51.087,198.457,59.643,205.915,63.921,209.643L68.199,213.372\"/><path marker-end=\"url(#my-svg_flowchart-v2-pointEnd)\" style=\"\" class=\"edge-thickness-normal edge-pattern-solid edge-thickness-normal edge-pattern-solid flowchart-link\" id=\"L_C_D_3\" d=\"M161.852,166L161.852,170.167C161.852,174.333,161.852,182.667,157.574,190.562C153.296,198.457,144.74,205.915,140.462,209.643L136.184,213.372\"/></g><g class=\"edgeLabels\"><g class=\"edgeLabel\"><g transform=\"translate(0, 0)\" class=\"label\"><foreignObject height=\"0\" width=\"0\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" class=\"labelBkg\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"edgeLabel\"></span></div></foreignObject></g></g><g class=\"edgeLabel\"><g transform=\"translate(0, 0)\" class=\"label\"><foreignObject height=\"0\" width=\"0\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" class=\"labelBkg\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"edgeLabel\"></span></div></foreignObject></g></g><g class=\"edgeLabel\"><g transform=\"translate(0, 0)\" class=\"label\"><foreignObject height=\"0\" width=\"0\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" class=\"labelBkg\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"edgeLabel\"></span></div></foreignObject></g></g><g class=\"edgeLabel\"><g transform=\"translate(0, 0)\" class=\"label\"><foreignObject height=\"0\" width=\"0\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" class=\"labelBkg\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"edgeLabel\"></span></div></foreignObject></g></g></g><g class=\"nodes\"><g transform=\"translate(102.19140625, 35)\" id=\"flowchart-A-0\" class=\"node default\"><rect height=\"54\" width=\"69.4375\" y=\"-27\" x=\"-34.71875\" style=\"\" class=\"basic label-container\"/><g transform=\"translate(-4.71875, -12)\" style=\"\" class=\"label\"><rect/><foreignObject height=\"24\" width=\"9.4375\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"nodeLabel\"><p>A</p></span></div></foreignObject></g></g><g transform=\"translate(42.53125, 139)\" id=\"flowchart-B-1\" class=\"node default\"><rect height=\"54\" width=\"69.0625\" y=\"-27\" x=\"-34.53125\" style=\"\" class=\"basic label-container\"/><g transform=\"translate(-4.53125, -12)\" style=\"\" class=\"label\"><rect/><foreignObject height=\"24\" width=\"9.0625\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"nodeLabel\"><p>B</p></span></div></foreignObject></g></g><g transform=\"translate(161.8515625, 139)\" id=\"flowchart-C-3\" class=\"node default\"><rect height=\"54\" width=\"69.578125\" y=\"-27\" x=\"-34.7890625\" style=\"\" class=\"basic label-container\"/><g transform=\"translate(-4.7890625, -12)\" style=\"\" class=\"label\"><rect/><foreignObject height=\"24\" width=\"9.578125\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"nodeLabel\"><p>C</p></span></div></foreignObject></g></g><g transform=\"translate(102.19140625, 243)\" id=\"flowchart-D-5\" class=\"node default\"><rect height=\"54\" width=\"69.8125\" y=\"-27\" x=\"-34.90625\" style=\"\" class=\"basic label-container\"/><g transform=\"translate(-4.90625, -12)\" style=\"\" class=\"label\"><rect/><foreignObject height=\"24\" width=\"9.8125\"><div style=\"display: table-cell; white-space: nowrap; line-height: 1.5; max-width: 200px; text-align: center;\" xmlns=\"http://www.w3.org/1999/xhtml\"><span class=\"nodeLabel\"><p>D</p></span></div></foreignObject></g></g></g></g></g></svg></div>\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newMockReader(tt.calls)
			c := New()
			got, err := c.Convert(r, tt.path, tt.source)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			r.AssertExpectations(t)
		})
	}
}

type testRead struct {
	path    string
	content string
	err     error
}

type mockFile struct {
	mock.Mock
}

func newMockReader(calls []testRead) *mockFile {
	m := new(mockFile)
	for _, call := range calls {
		m.On("Read", call.path).Return(call.content, call.err)
	}
	return m
}

func (m *mockFile) Read(p string) (io.ReadCloser, error) {
	args := m.Mock.Called(p)
	rc := strings.NewReader(args.String(0))
	return io.NopCloser(rc), args.Error(1)
}
