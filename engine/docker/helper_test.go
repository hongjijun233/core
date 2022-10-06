package docker

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	coreutils "github.com/projecteru2/core/utils"

	"github.com/stretchr/testify/assert"
)

func TestCreateTarStream(t *testing.T) {
	buff := bytes.NewBufferString("test")
	rc := io.NopCloser(buff)
	fname, err := coreutils.TempFile(rc)
	assert.NoError(t, err)
	_, err = CreateTarStream(fname)
	assert.NoError(t, err)
}

func TestWithDumpFiles(t *testing.T) {
	data := map[string][]byte{
		"/tmp/test-1": []byte("1"),
		"/tmp/test-2": []byte("2"),
	}
	fp := []string{}
	for target, content := range data {
		withTarfileDump(context.TODO(), target, content, 0, 0, int64(0), func(target, tarfile string) error {
			assert.True(t, strings.HasPrefix(target, "/tmp/test"))
			fp = append(fp, tarfile)
			_, err := os.Stat(tarfile)
			assert.Nil(t, err)
			return nil
		})
	}
	for _, path := range fp {
		_, err := os.Stat(path)
		assert.True(t, os.IsNotExist(err))
	}
}
