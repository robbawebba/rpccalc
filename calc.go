package rpccalc

// CalcService provides simple arithmetic operations such as Add and Subtract
type CalcService struct{}

// Add returns the result of adding a and b
func (s *CalcService) Add(a, b int) int {
	return a + b
}

// Subtract returns the result of subtracting b from a
func (s *CalcService) Subtract(a, b int) int {
	return a - b
}
