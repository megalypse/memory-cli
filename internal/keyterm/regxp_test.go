package keyterm

import "testing"

func TestGetRegexMatchesCapitalCaseTerms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "single word", input: "Memória", want: "Memória"},
		{name: "multiple words", input: "Memória Muscular", want: "Memória Muscular"},
		{name: "connector", input: "Rio de Janeiro", want: "Rio de Janeiro"},
		{name: "multiple connectors", input: "Senhor dos Anéis de Poder", want: "Senhor dos Anéis de Poder"},
		{name: "prefix", input: "o Senhor dos Anéis", want: "o Senhor dos Anéis"},
		{name: "plural prefix", input: "as Crônicas de Nárnia", want: "as Crônicas de Nárnia"},
		{name: "term in text", input: "leia o Senhor dos Anéis hoje.", want: "o Senhor dos Anéis"},
	}

	rgx := GetRegex()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := rgx.FindString(test.input)
			if got != test.want {
				t.Fatalf("FindString(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestGetRegexStopsBeforeInvalidContinuation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "lowercase word", input: "Rio bonito", want: "Rio"},
		{name: "connector without capitalized word", input: "Rio de janeiro", want: "Rio"},
		{name: "prefix without capitalized word", input: "o senhor", want: ""},
	}

	rgx := GetRegex()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := rgx.FindString(test.input)
			if got != test.want {
				t.Fatalf("FindString(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
