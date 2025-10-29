package fuzzing

import (
	"testing"
	"unicode/utf8"
)

func FuzzReverse(f *testing.F) {
	seeds := []string{
		"",
		"a",
		"hello",
		"Hello, World!",
		"привет",
		"Hello👋World",
		"12345",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)

		doubleRev := Reverse(rev)
		if orig != doubleRev {
			t.Errorf("Reverse is not involutive: orig=%q, rev=%q, doubleRev=%q",
				orig, rev, doubleRev)
		}

		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8: orig=%q, rev=%q", orig, rev)
		}

		if utf8.RuneCountInString(orig) != utf8.RuneCountInString(rev) {
			t.Errorf("Rune count changed: orig=%d, rev=%d",
				utf8.RuneCountInString(orig), utf8.RuneCountInString(rev))
		}
	})
}

func FuzzIsPalindrome(f *testing.F) {
	seeds := []string{
		"",
		"a",
		"aa",
		"aba",
		"racecar",
		"A man a plan a canal Panama",
		"not a palindrome",
		"12321",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsPalindrome panicked on input %q: %v", s, r)
			}
		}()

		result := IsPalindrome(s)

		reversed := Reverse(s)
		reversedResult := IsPalindrome(reversed)

		if result != reversedResult {
			t.Errorf("IsPalindrome not consistent with reverse: s=%q, result=%v, reversed=%q, reversedResult=%v",
				s, result, reversed, reversedResult)
		}
	})
}

func FuzzParseConfig(f *testing.F) {
	seeds := [][]byte{
		[]byte(`{"host":"localhost","port":8080,"enabled":true}`),
		[]byte(`{}`),
		[]byte(`{"host":"example.com"}`),
		[]byte(`{"port":80}`),
		[]byte(`{"enabled":false}`),
		[]byte(`{"port":-1}`),
		[]byte(`{"port":99999}`),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseConfig panicked on input %q: %v", data, r)
			}
		}()

		cfg, err := ParseConfig(data)

		if err != nil {
			return
		}

		if cfg.Port < 0 || cfg.Port > 65535 {
			t.Errorf("ParseConfig accepted invalid port: %d", cfg.Port)
		}
	})
}

func FuzzParseURL(f *testing.F) {
	seeds := []string{
		"https://example.com",
		"http://example.com/path",
		"//example.com",
		"http://user:pass@host:8080/path?key=value#fragment",
		"",
		"not a url",
		"ftp://example.com",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, urlStr string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseURL panicked on input %q: %v", urlStr, r)
			}
		}()

		u, err := ParseURL(urlStr)
		if err != nil {
			return
		}

		if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" {
			t.Errorf("ParseURL accepted unsupported scheme: %q", u.Scheme)
		}

		serialized := u.String()
		_, err = ParseURL(serialized)
		if err != nil {
			t.Errorf("ParseURL->String->ParseURL failed: orig=%q, serialized=%q, err=%v",
				urlStr, serialized, err)
		}
	})
}

func FuzzSafeDivide(f *testing.F) {
	f.Add(10, 2)
	f.Add(0, 1)
	f.Add(-10, 2)
	f.Add(10, 0)
	f.Add(-10, -2)

	f.Fuzz(func(t *testing.T, a, b int) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("SafeDivide panicked on input (%d, %d): %v", a, b, r)
			}
		}()

		result, err := SafeDivide(a, b)

		if b == 0 {
			if err == nil {
				t.Errorf("SafeDivide(%d, %d): expected error for division by zero, got result=%d",
					a, b, result)
			}
			return
		}

		if err != nil {
			t.Errorf("SafeDivide(%d, %d): unexpected error: %v", a, b, err)
			return
		}

		if b != 0 && result*b != a && result*b != a-1 && result*b != a+1 {
			t.Logf("Note: integer division truncation: %d / %d = %d (check: %d * %d = %d)",
				a, b, result, result, b, result*b)
		}
	})
}

func FuzzValidateUsername(f *testing.F) {
	seeds := []string{
		"user",
		"user123",
		"user_name",
		"user-name",
		"ab",
		"a_very_long_username_that_exceeds_limit",
		"user@name",
		"user name",
		"",
		"пользователь",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, username string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ValidateUsername panicked on input %q: %v", username, r)
			}
		}()

		err := ValidateUsername(username)

		if len(username) < 3 && err == nil {
			t.Errorf("ValidateUsername accepted too short username: %q", username)
		}

		if len(username) > 20 && err == nil {
			t.Errorf("ValidateUsername accepted too long username: %q", username)
		}

		if !utf8.ValidString(username) && err == nil {
			t.Errorf("ValidateUsername accepted invalid UTF-8: %q", username)
		}
	})
}

func FuzzParseKeyValue(f *testing.F) {
	seeds := []string{
		"key=value",
		"key1=value1&key2=value2",
		"",
		"key=",
		"=value",
		"key",
		"key1=value1&key2=value2&key3=value3",
		"key=value=extra",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseKeyValue panicked on input %q: %v", input, r)
			}
		}()

		result, err := ParseKeyValue(input)

		if err != nil {
			return
		}

		if input == "" && len(result) != 0 {
			t.Errorf("ParseKeyValue(%q) returned non-empty map for empty input: %v",
				input, result)
		}

		for key := range result {
			if key == "" {
				t.Errorf("ParseKeyValue(%q) returned empty key in result: %v",
					input, result)
			}
		}
	})
}

func TestReverseBasic(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"👋", "👋"},
	}

	for _, tt := range tests {
		got := Reverse(tt.input)
		if got != tt.want {
			t.Errorf("Reverse(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseconfigBasic(t *testing.T) {
	cfg, err := ParseConfig([]byte(`{"host":"localhost","port":8080}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Host != "localhost" || cfg.Port != 8080 {
		t.Errorf("got %+v; want {Host:localhost Port:8080}", cfg)
	}
}
