package gomkv

import (
	"fmt"
	"strings"
)

// View ...
func View(r EBMLReader) error {
	els, err := r.GetElements()
	if err != nil {
		return err
	}

	eleTreeView := make([]string, 0)

	ele2str := func(e *Element) error {
		dr := e.DataRange()
		name := id2name[e.ID]
		prefix := strings.Repeat("--|", int(e.Level+1)) // depth starts from 0
		eleTreeView = append(eleTreeView, fmt.Sprintf("%s ID: 0x%x(%s) of level %d, size: %d @%d, data:[%d, %d)\n",
			prefix, e.ID, name, e.Level, e.Size, e.At, dr.Start, dr.End))
		return nil
	}

	if err = els.Iter(ele2str); err != nil {
		return err
	}

	fmt.Println(strings.Join(eleTreeView, ""))
	return nil
}
