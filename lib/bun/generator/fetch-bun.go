package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const MinBunCLIBytes = 1024
const MetaFile = "bun-fetch.json"

type BunMeta struct {
	Tag    string `json:"tag"`
	GOOS   string `json:"goos"`
	GOARCH string `json:"goarch"`
	Zip    string `json:"zip"`
}

func effectiveGOOS() string {
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = os.Getenv("INPUT_GOOS")
	}
	if goos == "" {
		goos = runtime.GOOS
	}
	return goos
}

func effectiveGOARCH() string {
	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = os.Getenv("INPUT_GOARCH")
	}
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	return goarch
}

func getBunPath(wd string) (string, error) {
	goos := effectiveGOOS()
	goarch := effectiveGOARCH()

	zipName, ok := bunZipName(goos, goarch)
	if !ok {
		log.Fatalf("fetchbun: no Bun zip for GOOS=%s GOARCH=%s", goos, goarch)
	}

	embedDir := filepath.Join(wd, "embed")
	if err := os.MkdirAll(embedDir, 0o755); err != nil {
		return "", err
	}

	out := filepath.Join(embedDir, "buncli")
	if goos == "windows" {
		out = filepath.Join(embedDir, "buncli.exe")
	}

	metaPath := filepath.Join(embedDir, MetaFile)

	log.Printf("fetchbun: target %s/%s (%s)", goos, goarch, zipName)

	apiClient := &http.Client{Timeout: 45 * time.Second}
	tag, apiErr := fetchLatestBunTag(apiClient)
	if apiErr != nil {
		log.Printf("fetchbun: notice: could not reach GitHub API (%v)", apiErr)
		if tryUseCachedBun(metaPath, out, goos, goarch, zipName) {
			return out, nil
		}
		if tryInstallFromSystemPATH(embedDir, out, metaPath, goos, goarch, zipName) {
			return out, nil
		}
		return "", fmt.Errorf("fetchbun: no cached Bun, GitHub unreachable, and no usable `bun` on PATH (same host OS/ARCH only)")
	}

	if cacheMatchesLatestGitHubRelease(metaPath, out, tag, goos, goarch, zipName) {
		log.Printf("fetchbun: notice: using local cache (already at latest %s for %s/%s)", tag, goos, goarch)
		return out, nil
	}

	if binaryUsable(out) {
		log.Printf("fetchbun: notice: fetching %s from GitHub (updating existing cache)", tag)
	} else {
		log.Printf("fetchbun: notice: no usable cache — downloading Bun %s from GitHub", tag)
	}

	body, dlErr := downloadReleaseZip(tag, zipName)
	if dlErr != nil {
		log.Printf("fetchbun: notice: download failed (%v)", dlErr)
		if tryUseCachedBun(metaPath, out, goos, goarch, zipName) {
			log.Printf("fetchbun: notice: kept local cached Bun instead of release zip")
			return out, nil
		}
		if tryInstallFromSystemPATH(embedDir, out, metaPath, goos, goarch, zipName) {
			return out, nil
		}
		return "", fmt.Errorf("fetchbun: download failed, no usable cache, and no PATH fallback")
	}

	payload, err := unzipBunBinary(body, goos, zipName)
	if err != nil {
		log.Printf("fetchbun: notice: unzip failed (%v)", err)
		if tryUseCachedBun(metaPath, out, goos, goarch, zipName) {
			log.Printf("fetchbun: notice: kept local cached Bun instead")
			return out, nil
		}
		if tryInstallFromSystemPATH(embedDir, out, metaPath, goos, goarch, zipName) {
			return out, nil
		}
		return "", err
	}

	meta := BunMeta{Tag: tag, GOOS: goos, GOARCH: goarch, Zip: zipName}
	if err := writeBunCLI(embedDir, out, metaPath, payload, meta); err != nil {
		return "", err
	}

	return out, nil
}

func bunZipName(goos, goarch string) (string, bool) {
	type key struct{ os, arch string }
	m := map[key]string{
		{"linux", "amd64"}:   "bun-linux-x64.zip",
		{"linux", "arm64"}:   "bun-linux-aarch64.zip",
		{"darwin", "amd64"}:  "bun-darwin-x64.zip",
		{"darwin", "arm64"}:  "bun-darwin-aarch64.zip",
		{"windows", "amd64"}: "bun-windows-x64.zip",
	}
	if z, ok := m[key{goos, goarch}]; ok {
		return z, true
	}
	return "", false
}

func cacheMatchesLatestGitHubRelease(metaPath, binPath, tag, goos, goarch, zipName string) bool {
	if !binaryUsable(binPath) {
		return false
	}
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return false
	}
	var m BunMeta
	if json.Unmarshal(data, &m) != nil {
		return false
	}
	return m.Tag == tag && m.GOOS == goos && m.GOARCH == goarch && m.Zip == zipName
}

func binaryUsable(binPath string) bool {
	st, err := os.Stat(binPath)
	return err == nil && st.Size() >= MinBunCLIBytes
}

func tryUseCachedBun(metaPath, binPath, goos, goarch, zipName string) bool {
	if !binaryUsable(binPath) {
		return false
	}
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("fetchbun: notice: using %s without %s (version unknown)", filepath.Base(binPath), MetaFile)
			return true
		}
		return false
	}
	var m BunMeta
	if json.Unmarshal(data, &m) != nil {
		log.Printf("fetchbun: notice: using %s; %s is invalid JSON (version unknown)", filepath.Base(binPath), MetaFile)
		return true
	}
	if m.GOOS != goos || m.GOARCH != goarch {
		return false
	}
	if m.Zip != zipName {
		return false
	}
	log.Printf("fetchbun: notice: using local cache (Bun %s); GitHub not used", m.Tag)
	return true
}

func tryInstallFromSystemPATH(embedDir, out, metaPath, goos, goarch, zipName string) bool {
	if goos != runtime.GOOS || goarch != runtime.GOARCH {
		log.Printf("fetchbun: notice: not using PATH `bun` (cross-compile to %s/%s from %s/%s)", goos, goarch, runtime.GOOS, runtime.GOARCH)
		return false
	}
	path, err := exec.LookPath("bun")
	if err != nil {
		log.Printf("fetchbun: notice: no `bun` on PATH (%v)", err)
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("fetchbun: notice: could not read %s (%v)", path, err)
		return false
	}
	if len(data) < MinBunCLIBytes {
		log.Printf("fetchbun: notice: %s too small (%d bytes)", path, len(data))
		return false
	}

	tag := strings.TrimSpace(bunVersionOutput(path))
	if tag == "" {
		tag = "system"
	}
	meta := BunMeta{Tag: tag, GOOS: goos, GOARCH: goarch, Zip: zipName}
	if err := writeBunCLI(embedDir, out, metaPath, data, meta); err != nil {
		log.Printf("fetchbun: notice: could not write embed copy (%v)", err)
		return false
	}
	log.Printf("fetchbun: notice: copied system Bun from %s into embed (%d bytes)", path, len(data))
	return true
}

func bunVersionOutput(bunExe string) string {
	cmd := exec.Command(bunExe, "--version")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func writeBunCLI(embedDir, out, metaPath string, payload []byte, meta BunMeta) error {
	_ = os.Remove(filepath.Join(embedDir, "buncli"))
	_ = os.Remove(filepath.Join(embedDir, "buncli.exe"))
	if err := os.WriteFile(out, payload, 0o644); err != nil {
		return err
	}
	if runtime.GOOS != "windows" {
		_ = os.Chmod(out, 0o755)
	}

	enc, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	enc = append(enc, '\n')
	return os.WriteFile(metaPath, enc, 0o644)
}

func downloadReleaseZip(tag, zipName string) ([]byte, error) {
	url := fmt.Sprintf("https://github.com/oven-sh/bun/releases/download/%s/%s", tag, zipName)
	client := &http.Client{Timeout: 15 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		snip, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("%s: %s", resp.Status, strings.TrimSpace(string(snip)))
	}
	return io.ReadAll(resp.Body)
}

func fetchLatestBunTag(client *http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/oven-sh/bun/releases/latest", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "deployit-fetchbun")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		snip, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("%s: %s", resp.Status, strings.TrimSpace(string(snip)))
	}
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.TagName == "" {
		return "", fmt.Errorf("empty tag_name in API response")
	}
	return payload.TagName, nil
}

func unzipBunBinary(zipBytes []byte, goos, zipName string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("zip: %w", err)
	}

	want := "bun"
	if goos == "windows" {
		want = "bun.exe"
	}

	for _, f := range zr.File {
		if filepath.Base(f.Name) != want {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open zip member: %w", err)
		}
		b, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read member: %w", err)
		}
		if len(b) < MinBunCLIBytes {
			return nil, fmt.Errorf("%s in %s too small (%d bytes)", want, zipName, len(b))
		}
		return b, nil
	}
	return nil, fmt.Errorf("%s not found in %s", want, zipName)
}
