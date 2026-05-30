package gamepack

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/pipeline"
)

// PackOptions configures moonbasic pack.
type PackOptions struct {
	AssetsDir       string // relative to source file dir; default "assets"
	OutZip          string // output .zip path; default <basename>-pack.zip
	ExcludeRuntime  bool   // when true, omit moonbasic executable from the zip
}

// Pack compiles sourceMB and bundles the .mbc plus assets into a zip archive.
func Pack(sourceMB string, opts PackOptions) (string, error) {
	sourceMB, err := filepath.Abs(sourceMB)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(filepath.Ext(sourceMB), ".mb") {
		return "", fmt.Errorf("pack: expected .mb source file")
	}
	prog, err := pipeline.CompileFile(sourceMB)
	if err != nil {
		return "", err
	}
	base := strings.TrimSuffix(filepath.Base(sourceMB), filepath.Ext(sourceMB))
	outDir := filepath.Dir(sourceMB)
	mbcPath := filepath.Join(outDir, base+".mbc")
	data, err := pipeline.EncodeMOON(prog)
	if err != nil {
		return "", fmt.Errorf("pack: encode: %w", err)
	}
	if err := os.WriteFile(mbcPath, data, 0o644); err != nil {
		return "", fmt.Errorf("pack: write mbc: %w", err)
	}
	assetsRel := opts.AssetsDir
	if assetsRel == "" {
		assetsRel = "assets"
	}
	assetsDir := filepath.Join(outDir, filepath.FromSlash(assetsRel))
	outZip := opts.OutZip
	if outZip == "" {
		outZip = filepath.Join(outDir, base+"-pack.zip")
	}
	if err := os.MkdirAll(filepath.Dir(outZip), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(outZip)
	if err != nil {
		return "", err
	}
	defer f.Close()
	zw := zip.NewWriter(f)

	if err := addFileToZip(zw, mbcPath, base+".mbc"); err != nil {
		zw.Close()
		return "", err
	}
	if st, err := os.Stat(assetsDir); err == nil && st.IsDir() {
		if err := addDirToZip(zw, assetsDir, assetsRel); err != nil {
			zw.Close()
			return "", err
		}
	}
	readme := fmt.Sprintf("# %s pack\n\nExtract and run: moonbasic --run %s\nOr use bundled moonbasic: ./moonbasic --run %s\nSource: %s\n", base, base+".mbc", base+".mbc", sourceMB)
	if err := addBytesToZip(zw, "README-pack.txt", []byte(readme)); err != nil {
		zw.Close()
		return "", err
	}
	if !opts.ExcludeRuntime {
		if err := addRuntimeBinary(zw); err != nil {
			zw.Close()
			return "", err
		}
	}
	if err := zw.Close(); err != nil {
		return "", err
	}
	return outZip, nil
}

func addFileToZip(zw *zip.Writer, srcPath, nameInZip string) error {
	body, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return addBytesToZip(zw, filepath.ToSlash(nameInZip), body)
}

func addBytesToZip(zw *zip.Writer, name string, data []byte) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func addDirToZip(zw *zip.Writer, srcDir, prefixInZip string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(filepath.Join(prefixInZip, rel))
		return addFileToZip(zw, path, name)
	})
}

func addRuntimeBinary(zw *zip.Writer) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("pack: locate runtime: %w", err)
	}
	name := "moonbasic"
	if strings.EqualFold(filepath.Ext(exe), ".exe") {
		name = "moonbasic.exe"
	}
	return addFileToZip(zw, exe, name)
}
