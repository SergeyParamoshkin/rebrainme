package tabledriven

import "testing"

func TestToUpper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "hello", "HELLO"},
		{"uppercase", "WORLD", "WORLD"},
		{"mixed", "Hello World", "HELLO WORLD"},
		{"empty", "", ""},
		{"numbers", "test123", "TEST123"},
		{"cyrillic", "привет", "ПРИВЕТ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			got := ToUpper(tt.input)
			if got != tt.want {
				t.Errorf("ToUpper(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "hello", "olleh"},
		{"single char", "a", "a"},
		{"empty", "", ""},
		{"palindrome", "noon", "noon"},
		{"with spaces", "hello world", "dlrow olleh"},
		{"unicode", "привет", "тевирп"},
		{"emoji", "Hello👋World", "dlroW👋olleH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reverse(tt.input)
			if got != tt.want {
				t.Errorf("Reverse(%q) = %q; want %q", tt.input, got, tt.want)
			}

			doubleReverse := Reverse(got)
			if doubleReverse != tt.input {
				t.Errorf("Reverse(Reverse(%q)) = %q; want %q", tt.input, doubleReverse, tt.input)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"simple palindrome", "noon", true},
		{"not palindrome", "hello", false},
		{"single char", "a", true},
		{"empty", "", true},
		{"with spaces", "race car", true},
		{"case insensitive", "RaceCar", true},
		{"phrase", "A man a plan a canal Panama", true},
		{"cyrillic", "казак", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(tt.input)
			if got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCountVowels(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"all vowels", "aeiou", 5},
		{"no vowels", "bcdfg", 0},
		{"mixed", "hello", 2},
		{"empty", "", 0},
		{"uppercase", "HELLO", 2},
		{"sentence", "The quick brown fox", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountVowels(tt.input)
			if got != tt.want {
				t.Errorf("CountVowels(%q) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"no @", "testexample.com", false},
		{"no domain", "test@", false},
		{"no username", "@example.com", false},
		{"no dot in domain", "test@example", false},
		{"multiple @", "test@@example.com", false},
		{"valid subdomain", "user@mail.example.com", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidEmail(tt.email)
			if got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v; want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short string", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"needs truncation", "hello world", 8, "hello..."},
		{"very short maxLen", "hello", 3, "hel"},
		{"unicode", "привет мир", 8, "приве..."},
		{"emoji", "Hello👋World", 8, "Hello..."},
		{"zero maxLen", "hello", 0, ""},
		{"negative maxLen", "hello", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("TruncateString(%q, %d) = %q; want %q",
					tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestTableDrivenWithSetup(t *testing.T) {
	type testData struct {
		value    string
		expected int
	}

	tests := []struct {
		name string
		data []testData
	}{
		{
			name: "vowels",
			data: []testData{
				{"hello", 2},
				{"world", 1},
			},
		},
		{
			name: "no vowels",
			data: []testData{
				{"bcdfg", 0},
				{"xyz", 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, td := range tt.data {
				got := CountVowels(td.value)
				if got != td.expected {
					t.Errorf("CountVowels(%q) = %d; want %d",
						td.value, got, td.expected)
				}
			}
		})
	}
}
