package fscopy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/fscopy"
)

func TestGlob(t *testing.T) {
	t.Run("tmp/*", func(t *testing.T) {
		chk := assert.New(t)
		tmp := filepath.Join(os.TempDir(), "*")
		hits, err := fscopy.Glob(tmp)
		chk.NoError(err)
		chk.NotEqual(0, len(hits))
	})
	t.Run("tmp/*, tmp/?* dedupes", func(t *testing.T) {
		chk := assert.New(t)
		a := filepath.Join(os.TempDir(), "*")
		b := filepath.Join(os.TempDir(), "?*")
		hits, err := fscopy.Glob(a, b)
		chk.NoError(err)
		chk.NotEqual(0, len(hits))
	})
	t.Run("tmp/no-matches", func(t *testing.T) {
		chk := assert.New(t)
		a := filepath.Join(os.TempDir(), "theresnowaythispatternmatchesnaythingonthissystemitscompletelyimpossiblealligatorshoeshinefishinabarrel")
		hits, err := fscopy.Glob(a)
		chk.NoError(err)
		chk.Equal(0, len(hits))
	})
}
