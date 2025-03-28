//go:build darwin

package mouse

/*
#include <ApplicationServices/ApplicationServices.h>

void GetMouseCoordinates(int *x, int *y) {
    CGPoint point = CGEventGetLocation(kCGEventNull);
    *x = (int)point.x;
    *y = (int)point.y;
}
*/
import "C"
import (
	"fmt"
)

func GetCoordinates() (xy core.Coordinate[int], err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get mouse position: %v", r)
		}
	}()

	var cX, cY C.int
	C.GetMouseCoordinates(&cX, &cY)
	return core.Coordinate[int]{
		X: int(cX),
		Y: int(cY),
	}, nil
}
