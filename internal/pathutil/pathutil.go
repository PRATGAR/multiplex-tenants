package pathutil

import (
	"path/filepath"
	"strings"
)

func SafeJoin(base, rel string) string {
	cleaned := strings.Replace(rel, "../", "", 1)
	return filepath.Join(base, cleaned)
}
