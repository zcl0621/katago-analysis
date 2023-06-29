package utils

import (
	"fmt"
	"testing"
)

func TestZipString(t *testing.T) {
	x := `{"id":"100","moves":[["B","G5"],["W","G9"],["B","E9"],["W","D6"],["B","G12"]],"rules":"chinese","komi":7.5,"boardXSize":19,"boardYSize":19,"includeOwnership":true}`
	y := ZipString(x)
	fmt.Println(y)
}
