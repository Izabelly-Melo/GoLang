package tickets_test

import (
	"testing"

	"github.com/bootcamp-go/desafio-go-bases/internal/tickets"
	"github.com/stretchr/testify/require"
)

func TestGetTotalTickets(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		result, _ := tickets.GetTotalTickets("Brazil")
		expected := 45
		require.Equal(t, expected, result)
	})

	t.Run("2", func(t *testing.T) {
		result, _ := tickets.GetTotalTickets("Philippines")
		expected := 46
		require.Equal(t, result, expected)
	})
}

func TestGetMornings(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		result, _ := tickets.GetMornings("noite")
		expected := 151
		require.Equal(t, result, expected)
	})

	t.Run("2", func(t *testing.T) {
		result, _ := tickets.GetMornings("madrugada")
		expected := 304
		require.Equal(t, result, expected)
	})

	t.Run("3", func(t *testing.T) {
		result, _ := tickets.GetMornings("manha")
		expected := 256
		require.Equal(t, result, expected)
	})

	t.Run("4", func(t *testing.T) {
		result, _ := tickets.GetMornings("tarde")
		expected := 289
		require.Equal(t, result, expected)
	})
}

func TestAverageDestination(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		result, _ := tickets.AverageDestination("Brazil")
		expected := 4
		require.Equal(t, result, expected)
	})

	t.Run("2", func(t *testing.T) {
		result, _ := tickets.AverageDestination("Philippines")
		expected := 4
		require.Equal(t, result, expected)
	})
}
