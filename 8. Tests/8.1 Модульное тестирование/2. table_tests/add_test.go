package main

import (
	"testing"
)

func TestAdd_TableDriven(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{2, 3, 5},
		{-1, 1, 0},
		{0, 0, 0},
		{100, -50, 50},
	}

	for _, test := range tests {
		result := Add(test.a, test.b)
		if result != test.expected {
			t.Errorf("Add(%d, %d) = %d; ожидалось %d",
				test.a, test.b, result, test.expected)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"положительные_числа", 2, 3, 5},
		{"отрицательные_числа", -2, -3, -5},
		{"смешанные_числа", -2, 3, 1},
		{"ноль", 0, 5, 5},
		{"большие_числа", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; ожидалось %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

