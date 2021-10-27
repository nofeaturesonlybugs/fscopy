package fscopy

import (
	"io"
	"os"
	"path/filepath"

	"github.com/nofeaturesonlybugs/errors"
)

// Dir recursively copies src/* into dst/.  If dst does not exist it is created with os.MkdirAll().
func Dir(dst, src string) error {
	seen := map[string]struct{}{}
	//
	var dir func(dst, src string) error
	dir = func(dst, src string) error {
		var entries []os.DirEntry
		var info os.FileInfo
		var err error
		//
		// This is a low effort attempt at preventing circular references.
		key := dst + "\x1e" + src
		if _, ok := seen[key]; ok {
			return nil
		}
		seen[key] = struct{}{}
		//
		if info, err = os.Stat(src); err != nil {
			return errors.Go(err)
		} else if !info.IsDir() {
			return errors.Errorf("src is not a directory: %v", src)
		}
		//
		if err = os.MkdirAll(dst, info.Mode()); err != nil {
			return errors.Go(err)
		}
		//
		if entries, err = os.ReadDir(src); err != nil {
			return errors.Go(err)
		}
		for _, entry := range entries {
			from := filepath.Join(src, entry.Name())
			to := filepath.Join(dst, entry.Name())
			if entry.IsDir() {
				if err = dir(to, from); err != nil {
					return errors.Go(err)
				}
				continue
			}
			if err = File(to, from); err != nil {
				return errors.Go(err)
			}
		}
		//
		return nil
	}
	//
	return dir(dst, src)
}

// File copies src file to dst.  If dst already exists it is overwritten.
func File(dst, src string) error {
	var dfd, sfd *os.File
	var info os.FileInfo
	var err error
	//
	if sfd, err = os.Open(src); err != nil {
		return errors.Go(err)
	}
	defer sfd.Close()
	//
	if dfd, err = os.Create(dst); err != nil {
		return errors.Go(err)
	}
	defer dfd.Close()
	//
	if _, err = io.Copy(dfd, sfd); err != nil {
		return errors.Go(err)
	}
	//
	if err = dfd.Close(); err != nil {
		return errors.Go(err)
	}
	if err = sfd.Close(); err != nil {
		return errors.Go(err)
	}
	//
	if info, err = os.Stat(src); err != nil {
		return errors.Go(err)
	}
	if err = os.Chmod(dst, info.Mode()); err != nil {
		return errors.Go(err)
	}
	//
	return nil
}

// Into copies each src into dst/.  If dst does not exist it is created with mode 0770.  If a src is
// a directory then its location in dst is dst/+basename(src).
func Into(dst string, src ...string) error {
	var info os.FileInfo
	var err error
	//
	if err = os.MkdirAll(dst, 0770); err != nil {
		return errors.Go(err)
	}
	for _, from := range src {
		to := filepath.Join(dst, filepath.Base(from))
		if info, err = os.Stat(from); err != nil {
			return errors.Go(err)
		}
		if info.IsDir() {
			if err = Dir(to, from); err != nil {
				return errors.Go(err)
			}
			continue
		}
		if err = File(to, from); err != nil {
			return errors.Go(err)
		}
	}
	//
	return nil
}
