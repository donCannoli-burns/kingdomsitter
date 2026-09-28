package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DefaultAddr = "127.0.0.1:10423"

type healthResponse struct {
	OK                 bool   `json:"ok"`
	Name               string `json:"name"`
	ParserBackend      string `json:"parser_backend"`
	ExecutionAuthority bool   `json:"execution_authority"`
	StateSeen          bool   `json:"state_seen"`
	LastStateUpdate    string `json:"last_state_update,omitempty"`
}

func Start(addr string) (string, error) {
	if addr == "" {
		addr = DefaultAddr
	}
	if h, err := health(addr); err == nil && h.OK {
		return fmt.Sprintf("kingdomsitter sidecar already running on http://%s", addr), nil
	}

	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable: %w", err)
	}
	pidFile, logFile, err := runtimePaths()
	if err != nil {
		return "", err
	}
	logHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return "", fmt.Errorf("open sidecar log: %w", err)
	}

	cmd := exec.Command(exe, "serve", addr)
	cmd.Stdin = nil
	cmd.Stdout = logHandle
	cmd.Stderr = logHandle
	cmd.SysProcAttr = detachedSysProcAttr()

	if err := cmd.Start(); err != nil {
		_ = logHandle.Close()
		return "", fmt.Errorf("start detached sidecar: %w", err)
	}
	pid := cmd.Process.Pid
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)+"\n"), 0o600); err != nil {
		_ = cmd.Process.Kill()
		_ = logHandle.Close()
		return "", fmt.Errorf("write pid file: %w", err)
	}
	_ = cmd.Process.Release()
	_ = logHandle.Close()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if h, err := health(addr); err == nil && h.OK {
			return fmt.Sprintf("kingdomsitter sidecar started in background on http://%s (pid %d; log %s)", addr, pid, logFile), nil
		}
		time.Sleep(75 * time.Millisecond)
	}

	return "", fmt.Errorf("sidecar process started as pid %d but health check did not become ready; see %s", pid, logFile)
}

func Stop() (string, error) {
	pidFile, _, err := runtimePaths()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(pidFile)
	if os.IsNotExist(err) {
		return "kingdomsitter sidecar is not managed by a local pid file", nil
	}
	if err != nil {
		return "", fmt.Errorf("read pid file: %w", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return "", fmt.Errorf("invalid pid file %s: %w", pidFile, err)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return "", fmt.Errorf("find pid %d: %w", pid, err)
	}
	if err := stopProcess(process); err != nil {
		return "", fmt.Errorf("stop pid %d: %w", pid, err)
	}
	_ = os.Remove(pidFile)
	return fmt.Sprintf("kingdomsitter sidecar stopped (pid %d)", pid), nil
}

func Status(addr string) (string, error) {
	if addr == "" {
		addr = DefaultAddr
	}
	_, logFile, err := runtimePaths()
	if err != nil {
		return "", err
	}
	h, err := health(addr)
	if err != nil || !h.OK {
		return fmt.Sprintf("kingdomsitter sidecar OFF on http://%s (log %s)", addr, logFile), nil
	}
	state := "no KoLmafia state received yet"
	if h.StateSeen {
		state = "KoLmafia state seen"
		if h.LastStateUpdate != "" {
			state += " at " + h.LastStateUpdate
		}
	}
	return fmt.Sprintf("kingdomsitter sidecar ON on http://%s (%s; execution_authority=%t)", addr, state, h.ExecutionAuthority), nil
}

func InstallCondaHook() (string, error) {
	prefix := os.Getenv("CONDA_PREFIX")
	if prefix == "" {
		return "", fmt.Errorf("CONDA_PREFIX is not set; activate the target Conda environment first")
	}
	dir := filepath.Join(prefix, "etc", "conda", "activate.d")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create Conda activate.d: %w", err)
	}
	path := filepath.Join(dir, "kingdomsitter.sh")
	content := "#!/usr/bin/env bash\n" +
		"# Managed by: kingdomsitter conda-install\n" +
		"if command -v kingdomsitter >/dev/null 2>&1; then\n" +
		"  kingdomsitter start >/dev/null 2>&1 || true\n" +
		"fi\n"
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		return "", fmt.Errorf("write Conda activation hook: %w", err)
	}
	return path, nil
}

func RemoveCondaHook() (string, error) {
	prefix := os.Getenv("CONDA_PREFIX")
	if prefix == "" {
		return "", fmt.Errorf("CONDA_PREFIX is not set; activate the target Conda environment first")
	}
	path := filepath.Join(prefix, "etc", "conda", "activate.d", "kingdomsitter.sh")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("remove Conda activation hook: %w", err)
	}
	return path, nil
}

func health(addr string) (healthResponse, error) {
	client := &http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get("http://" + addr + "/health")
	if err != nil {
		return healthResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return healthResponse{}, fmt.Errorf("health returned %s", resp.Status)
	}
	var h healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&h); err != nil {
		return healthResponse{}, err
	}
	if !h.OK || h.Name != "kingdomsitter" {
		return healthResponse{}, fmt.Errorf("unexpected health response")
	}
	return h, nil
}

func runtimePaths() (pidFile string, logFile string, err error) {
	root, err := os.UserCacheDir()
	if err != nil {
		return "", "", fmt.Errorf("resolve user cache directory: %w", err)
	}
	dir := filepath.Join(root, "kingdomsitter")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", "", fmt.Errorf("create runtime directory: %w", err)
	}
	return filepath.Join(dir, "kingdomsitter.pid"), filepath.Join(dir, "kingdomsitter.log"), nil
}
