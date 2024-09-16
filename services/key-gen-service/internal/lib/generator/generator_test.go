package generator

import "testing"

func TestGenerator_Generate(t *testing.T) {
	type fields struct {
		length int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test lengtth of 7",
			fields: fields{
				length: 7,
			},
		},
		{
			name: "Test lengtth of 7",
			fields: fields{
				length: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator(tt.fields.length)
			if got := g.Generate(); len(got) != tt.fields.length {
				t.Errorf("Generator.Generate() = %d, want %d", len(got), tt.fields.length)
			}
		})
	}
}
