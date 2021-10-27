package fscopy_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nofeaturesonlybugs/fscopy"
	"github.com/stretchr/testify/assert"
)

// writeFiles write each file in files[name]data into dst.
func writeFiles(dst string, files map[string]string) error {
	var err error
	//
	if err = os.MkdirAll(dst, 0770); err != nil {
		return err
	}
	for name, contents := range files {
		to := filepath.Join(dst, name)
		if err = os.WriteFile(to, []byte(contents), 0770); err != nil {
			return err
		}
	}
	//
	return nil
}

func TestDir(t *testing.T) {
	t.Run("dir", func(t *testing.T) {
		chk := assert.New(t)
		dst, err := os.MkdirTemp("", "gotest-dir*")
		chk.NoError(err)
		defer os.Remove(dst)
		//
		err = writeFiles(dst, map[string]string{
			"a.txt": "This is file a",
			"b.txt": "This is file b",
			"other": "other file",
		})
		chk.NoError(err)
		err = writeFiles(filepath.Join(dst, "subdir"), map[string]string{
			"suba.txt": "This is file a",
			"subb.txt": "This is file b",
			"subother": "other file",
		})
		chk.NoError(err)
		//
		dstCopy := dst + ".copy"
		err = fscopy.Dir(dstCopy, dst)
		chk.NoError(err)
		defer os.RemoveAll(dstCopy)
	})
}

func TestFile(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		chk := assert.New(t)
		sfd, err := os.CreateTemp("", "gotest-copy*.txt")
		chk.NoError(err)
		defer os.Remove(sfd.Name())
		_, err = io.Copy(sfd, strings.NewReader("Hello, World!"))
		chk.NoError(err)
		err = sfd.Close()
		chk.NoError(err)
		//
		to := sfd.Name() + ".cp"
		defer os.Remove(to)
		err = fscopy.File(to, sfd.Name())
		chk.NoError(err)
	})
	t.Run("src-does-not-exist", func(t *testing.T) {
		chk := assert.New(t)
		src := filepath.Join(os.TempDir(), "thisfilecannotpossiblyexistonthissystemaldjfwoifjweoafijlwefjalwkejfbunchofjunk.txt")
		to := src + ".cp"
		defer os.Remove(to)
		err := fscopy.File(to, src)
		chk.Error(err)
	})
}

func TestInto(t *testing.T) {
	t.Run("dir", func(t *testing.T) {
		chk := assert.New(t)
		dst, err := os.MkdirTemp("", "gotest-dir*")
		chk.NoError(err)
		defer os.Remove(dst)
		//
		err = writeFiles(dst, map[string]string{
			"a.txt": "This is file a",
			"b.txt": "This is file b",
			"other": "other file",
		})
		chk.NoError(err)
		subdir := filepath.Join(dst, "subdir")
		err = writeFiles(subdir, map[string]string{
			"suba.txt": "This is file a",
			"subb.txt": "This is file b",
			"subother": "other file",
		})
		chk.NoError(err)
		//
		into, err := os.MkdirTemp("", "gotest-into*")
		chk.NoError(err)
		defer os.RemoveAll(into)
		sources := []string{
			filepath.Join(dst, "a.txt"),
			filepath.Join(dst, "b.txt"),
			filepath.Join(dst, "other"),
			dst,
			subdir,
		}
		err = fscopy.Into(into, sources...)
		chk.NoError(err)
		//
		info, err := os.Stat(filepath.Join(into, filepath.Base(dst)))
		chk.NoError(err)
		chk.True(info.IsDir())
		//
		info, err = os.Stat(filepath.Join(into, filepath.Base(dst), "subdir"))
		chk.NoError(err)
		chk.True(info.IsDir())
		//
		info, err = os.Stat(filepath.Join(into, "subdir"))
		chk.NoError(err)
		chk.True(info.IsDir())
		//
		files := []string{
			filepath.Join(into, "a.txt"),
			filepath.Join(into, "b.txt"),
			filepath.Join(into, "other"),
			filepath.Join(into, filepath.Base(dst), "a.txt"),
			filepath.Join(into, filepath.Base(dst), "b.txt"),
			filepath.Join(into, filepath.Base(dst), "other"),
			filepath.Join(into, filepath.Base(dst), "subdir", "suba.txt"),
			filepath.Join(into, filepath.Base(dst), "subdir", "subb.txt"),
			filepath.Join(into, filepath.Base(dst), "subdir", "subother"),
			filepath.Join(into, "subdir", "suba.txt"),
			filepath.Join(into, "subdir", "subb.txt"),
			filepath.Join(into, "subdir", "subother"),
		}
		for _, file := range files {
			info, err := os.Stat(file)
			chk.NoError(err)
			chk.False(info.IsDir())
		}
	})
}
