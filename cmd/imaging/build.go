package imaging

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func Build(path string, imName string, imTag string) error {
	// goal:
	// 1. download nginx
	// 2. /usr/share/nginx/html <- delete this dir (new layer, appended on top of nginx)
	// 3. copy my blog there (new layer, appended on top of nginx)

	img, err := crane.Pull("nginx:latest")
	if err != nil {
		panic(err)
	}

	deleteMap := map[string][]byte{
		"usr/share/nginx/.wh.html": []byte{},
	}
	deleteLayer, err := crane.Layer(deleteMap)
	if err != nil {
		panic(err)
	}

	addLayer, err := layerFromDir(path)
	if err != nil {
		panic(err)
	}

	newImg, err := mutate.AppendLayers(img, deleteLayer, addLayer)
	if err != nil {
		panic(err)
	}

	tag, err := name.NewTag("adamlahbib/" + imName + ":" + imTag)
	if err != nil {
		panic(err)
	}

	if s, err := daemon.Write(tag, newImg); err != nil {
		panic(err)
	} else {
		fmt.Println(s)
	}

	log.Println(newImg)
	log.Println(tag.String())

	// push to remote registry
	if err := crane.Push(newImg, tag.String()); err != nil {
		panic(err)
	} else {
		log.Println("pushed image to adamlahbib's DockerHub")
	}

	// delete folder repos
	err = os.RemoveAll(path)
	if err != nil {
		panic(err)
	}

	return nil
}

func layerFromDir(root string) (v1.Layer, error) {
	var b bytes.Buffer
	tw := tar.NewWriter(&b)

	targetPath := "usr/share/nginx/html"

	err := filepath.Walk(root, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, err := filepath.Rel(root, fp)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		hdr := &tar.Header{
			Name: path.Join(targetPath, filepath.ToSlash(rel)),
			Mode: int64(info.Mode()),
		}

		if !info.IsDir() {
			hdr.Size = info.Size()
		}

		if info.Mode().IsDir() {
			hdr.Typeflag = tar.TypeDir
		} else if info.Mode().IsRegular() {
			hdr.Typeflag = tar.TypeReg
		} else {
			return fmt.Errorf("not implemented archiving file type %s (%s)", info.Mode(), rel)
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}
		if !info.IsDir() {
			f, err := os.Open(fp)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, f); err != nil {
				return fmt.Errorf("failed to read file into the tar: %w", err)
			}
			f.Close()
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan files: %w", err)
	}
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to finish tar: %w", err)
	}
	return tarball.LayerFromReader(&b)
}
