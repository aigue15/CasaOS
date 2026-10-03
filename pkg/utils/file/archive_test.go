package file

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mholt/archives"
)

func TestArchiveRoundTrip(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("bravo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "a.txt"), []byte("alpha"), 0o644); err != nil {
		t.Fatal(err)
	}

	cases := map[string]string{
		"":       ".zip",
		"zip":    ".zip",
		"tar":    ".tar",
		"targz":  ".tar.gz",
		"tarbz2": ".tar.bz2",
		"tarxz":  ".tar.xz",
		"tarlz4": ".tar.lz4",
		"tarsz":  ".tar.sz",
	}

	for algorithm, wantExt := range cases {
		t.Run(algorithm, func(t *testing.T) {
			ext, ar, err := GetCompressionAlgorithm(algorithm)
			if err != nil {
				t.Fatal(err)
			}
			if ext != wantExt {
				t.Fatalf("extension = %q, want %q", ext, wantExt)
			}

			commonDir := CommonPrefix(filepath.Separator, root)
			if err := AddFile(ar, root, commonDir); err != nil {
				t.Fatal(err)
			}

			var buf bytes.Buffer
			if err := ar.WriteTo(context.Background(), &buf); err != nil {
				t.Fatal(err)
			}

			got := map[string]string{}
			format, stream, err := archives.Identify(context.Background(), "download"+ext, bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatal(err)
			}
			err = format.(archives.Extractor).Extract(context.Background(), stream, func(_ context.Context, f archives.FileInfo) error {
				if f.IsDir() {
					got[f.NameInArchive] = "<dir>"
					return nil
				}
				rc, err := f.Open()
				if err != nil {
					return err
				}
				defer rc.Close()
				data, err := io.ReadAll(rc)
				got[f.NameInArchive] = string(data)
				return err
			})
			if err != nil {
				t.Fatal(err)
			}

			want := map[string]string{"b.txt": "bravo", "sub": "<dir>", "sub/a.txt": "alpha"}
			if algorithm == "" || algorithm == "zip" {
				want = map[string]string{"b.txt": "bravo", "sub/": "<dir>", "sub/a.txt": "alpha"}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("archive contents = %v, want %v", got, want)
			}
		})
	}

	if _, _, err := GetCompressionAlgorithm("rar"); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
