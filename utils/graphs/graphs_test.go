package graphs

import (
	"bytes"
	"strings"
	"testing"
)

func sampleTable() Table {
	return Table{
		Headings: []string{"Month", "Sales", "Costs"},
		Rows: [][]string{
			{"Jan", "12", "5"},
			{"Feb", "18", "7"},
			{"Mar", "", "9"},
		},
	}
}

func TestRenderTableSupportsTabularChartTypes(t *testing.T) {
	for _, chartType := range []Type{Line, Area, Bar, HorizontalBar, Scatter, Pie, Doughnut} {
		t.Run(string(chartType), func(t *testing.T) {
			valueColumns := []int{1, 2}
			if chartType == Pie || chartType == Doughnut {
				valueColumns = []int{1}
			}
			data, err := RenderTable(sampleTable(), Spec{
				Type:              chartType,
				CategoryColumn:    0,
				ValueColumns:      valueColumns,
				SeriesNames:       []string{"Sales", "Costs"},
				Title:             "Monthly totals",
				CategoryAxisTitle: "Month",
				ValueAxisTitle:    "Amount",
			})
			if err != nil {
				t.Fatalf("RenderTable() error = %v", err)
			}
			if len(data) == 0 || !bytes.Contains(data, []byte("<svg")) {
				t.Fatalf("RenderTable() returned %d bytes, want SVG", len(data))
			}
		})
	}
}

func TestRenderTableSupportsPNG(t *testing.T) {
	data, err := RenderTable(sampleTable(), Spec{
		Type:           Bar,
		Format:         PNG,
		CategoryColumn: 0,
		ValueColumns:   []int{1},
	})
	if err != nil {
		t.Fatalf("RenderTable() error = %v", err)
	}
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		t.Fatalf("RenderTable() returned %d bytes, want PNG signature", len(data))
	}
}

func TestRenderTableValidatesTabularMapping(t *testing.T) {
	tests := []struct {
		name string
		spec Spec
		want string
	}{
		{
			name: "pie needs one value column",
			spec: Spec{Type: Pie, CategoryColumn: 0, ValueColumns: []int{1, 2}},
			want: "require exactly one value column",
		},
		{
			name: "value column range",
			spec: Spec{Type: Line, CategoryColumn: 0, ValueColumns: []int{3}},
			want: "value column 3 is out of range",
		},
		{
			name: "category cannot also be a value",
			spec: Spec{Type: Line, CategoryColumn: 0, ValueColumns: []int{0}},
			want: "cannot also be a value column",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := RenderTable(sampleTable(), test.spec)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("RenderTable() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestRenderTableReportsInvalidNumberLocation(t *testing.T) {
	table := sampleTable()
	table.Rows[1][1] = "not a number"
	_, err := RenderTable(table, Spec{Type: Line, CategoryColumn: 0, ValueColumns: []int{1}})
	if err == nil || !strings.Contains(err.Error(), `row 3, column "Sales"`) {
		t.Fatalf("RenderTable() error = %v, want row and column context", err)
	}
}

func TestParseNumberStripsCommasAndCurrencySymbols(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{input: "$1,234.50", want: 1234.5},
		{input: "€1,234.50", want: 1234.5},
		{input: "-£2,000", want: -2000},
		{input: "¥12,345", want: 12345},
		{input: "1,234.50₹", want: 1234.5},
	}
	for _, test := range tests {
		got, err := parseNumber(test.input)
		if err != nil {
			t.Errorf("parseNumber(%q) error = %v", test.input, err)
			continue
		}
		if got != test.want {
			t.Errorf("parseNumber(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}

func TestParseNumberStillRejectsMalformedText(t *testing.T) {
	if _, err := parseNumber("$1,2x3"); err == nil {
		t.Fatal("parseNumber() error = nil, want malformed numeric text rejected")
	}
}
