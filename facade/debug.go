package facade

import (
	"fmt"

	"github.com/goark/gocli/rwi"
)

func debugPrint(ui *rwi.RWI, err error) error {
	if debugFlag && err != nil {
		fmt.Fprintf(ui.Writer(), "%+v\n", err)
		return nil
	}
	return err
}
