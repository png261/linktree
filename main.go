// main.go
package main

import (
	"fmt"
	"github.com/png261/linktree/configs"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	clean()
	build()
}

func clean() error {
	err := os.RemoveAll("output")
	if err != nil {
		return err
	}

	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output/ directory: %w", err)
	}
	fmt.Println("Clean sucessfully")
	return nil
}

func build() {
	config, err := configs.LoadSiteConfig("config.yml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = generateHTML(config)
	if err != nil {
		log.Fatal("Error generating HTML: %v", err)
	}

	fmt.Println("Site generated sucessfully!")
	err = copyDir(fmt.Sprintf("themes/%s/assets", config.Theme), "output/assets")

	if err != nil {
		log.Fatalf("Error copying and minifying assets: %v", err)
	}

	fmt.Println("Site generated sucessfully!")
}

func generateHTML(config *configs.SiteConfig) error {
	themeFile := fmt.Sprintf("themes/%s/index.html", config.Theme)

	tmpl, err := template.ParseFiles(themeFile)
	if err != nil {
		return err
	}

	outputFile, err := os.Create("output/index.html")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	data := struct {
		Config *configs.SiteConfig
	}{
		Config: config,
	}

	return tmpl.Execute(outputFile, data)
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}
