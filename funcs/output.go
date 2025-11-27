package funcs

import (
	"bufio"
	"fmt"
	"os"
)

func writeParamsToFile(path string, params []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, p := range params {
		if _, err := fmt.Fprintln(w, p); err != nil {
			return err
		}
	}
	return w.Flush()
}
