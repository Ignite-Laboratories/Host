//go:build darwin

package host

/*
#include <ApplicationServices/ApplicationServices.h>

void GetMouseCoordinates(int *x, int *y) {
    CGPoint point = CGEventGetLocation(kCGEventNull);
    *x = (int)point.x;
    *y = (int)point.y;
}
*/
import "C"
import "fmt"

func (m mouse) GetCoordinates() (x int, y int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get mouse position: %v", r)
		}
	}()

	var cX, cY C.int
	C.GetMouseCoordinates(&cX, &cY)
	return int(cX), int(cY), nil
}
