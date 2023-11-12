package validation

import "github.com/gookit/validate"

func newValidators() {
	validate.AddValidator("floatType", isFloat)
	validate.AddValidator("intType", isInt)
}

func isFloat(v any) bool {
	_, okFloat32 := v.(float32)
	_, okFloat64 := v.(float64)

	return okFloat32 || okFloat64
}

func isInt(v any) bool {
	_, okFloat64 := v.(float64)
	_, okInt := v.(int)

	return okFloat64 || okInt
}
