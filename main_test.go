package main

import "testing"

func TestHello(t *testing.T) {
	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("David", "Spanish")
		want := "Hola, David"

		assertCorrectMessage(t, got, want)
	})
	t.Run("in Korean", func(t *testing.T) {
		got := Hello("은영", "Korean")
		want := "안녕, 은영"

		assertCorrectMessage(t, got, want)
	})

	t.Run("in English", func(t *testing.T) {
		got := Hello("Chris", "English")
		want := "Hello, Chris"

		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
