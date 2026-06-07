package keyterm

import (
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "multiple terms",
			raw:  "estude Rio de Janeiro e leia o Senhor dos Anéis.",
			want: []string{"Rio de Janeiro", "o Senhor dos Anéis"},
		},
		{
			name: "single term",
			raw:  "entenda Memória Muscular.",
			want: []string{"Memória Muscular"},
		},
		{
			name: "no terms",
			raw:  "texto sem termos capitalizados",
		},
		{
			name: "empty input",
			raw:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := Extract(test.raw)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Extract(%q) = %#v, want %#v", test.raw, got, test.want)
			}
		})
	}
}
