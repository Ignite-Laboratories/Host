package math

import (
	"gonum.org/v1/gonum/mat"
)

// Ortho returns a 4x4 orthographic projection matrix.
//
// left, right: the visible range of X (horizontal axis in world coordinates)
// bottom, top: the visible range of Y (vertical axis in world coordinates)
// near, far: the visible range of Z (depth, used for layering 2D elements or 3D rendering)
//
// The matrix can be directly passed to OpenGL by first flattening into a 1D array
func Ortho(left, right, bottom, top, near, far float64) []float32 {
	// Create a 4x4 dense matrix with all elements set to zero
	ortho := mat.NewDense(4, 4, nil)

	// Compute orthographic projection elements
	ortho.Set(0, 0, 2/(right-left)) // row 0, col 0
	ortho.Set(1, 1, 2/(top-bottom)) // row 1, col 1
	ortho.Set(2, 2, -2/(far-near))  // row 2, col 2
	ortho.Set(3, 3, 1)              // row 3, col 3

	ortho.Set(0, 3, -(right+left)/(right-left)) // row 0, col 3
	ortho.Set(1, 3, -(top+bottom)/(top-bottom)) // row 1, col 3
	ortho.Set(2, 3, -(far+near)/(far-near))     // row 2, col 3

	flattened := make([]float32, 16)
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			flattened[col*4+row] = float32(ortho.At(row, col))
		}
	}

	return flattened
}
