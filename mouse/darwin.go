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
	"github.com/ignite-laboratories/core/std"
)

func init() {
	fmt.Println("[host] - Darwin - mouse")
}

func getCoordinates() std.XY[int] {
	var cX, cY C.int
	C.GetMouseCoordinates(&cX, &cY)
	return std.XY[int]{
		X: int(cX),
		Y: int(cY),
	}
}
