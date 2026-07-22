package tests

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/require"
)

func TestE2ERoomRestartPurgesRoomAndDropsSocket(t *testing.T) {
	subject := newE2ESubject(t)
	binaryPath := buildE2EBackend(t)
	backend := startE2EBackend(t, binaryPath, subject)
	t.Cleanup(func() {
		backend.stop(t)
	})

	roomSubject := *subject
	roomSubject.baseURL = backend.baseURL
	room := roomSubject.createRoom(t, "e2e-restart-"+e2eHash(subject.userID))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := dialE2ERoomSocket(t, ctx, &roomSubject, room.ID)
	defer conn.Close(websocket.StatusNormalClosure, "test done")

	backend.stop(t)
	readCtx, readCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer readCancel()
	var event e2eRoomSocketEvent
	require.Error(t, wsjson.Read(readCtx, conn, &event), "server restart must drop the active websocket")

	restartedBackend := startE2EBackend(t, binaryPath, subject)
	t.Cleanup(func() {
		restartedBackend.stop(t)
	})
	restartedSubject := *subject
	restartedSubject.baseURL = restartedBackend.baseURL

	var response e2eErrorResponse
	restartedSubject.doJSON(t, http.MethodGet, "/api/rooms/"+room.ID+"/", nil, http.StatusNotFound, &response)
	require.Equal(t, "room.not_found", response.Code)
}

type e2eBackendProcess struct {
	baseURL string
	command *exec.Cmd
	stderr  bytes.Buffer
	stdout  bytes.Buffer
}

func buildE2EBackend(t *testing.T) string {
	t.Helper()

	root := e2eProjectRoot(t)
	directory := t.TempDir()
	filename := "miskatonic-e2e-backend"
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}
	binaryPath := filepath.Join(directory, filename)
	command := exec.Command("go", "build", "-o", binaryPath, "./cmd")
	command.Dir = root
	output, err := command.CombinedOutput()
	require.NoErrorf(t, err, "build isolated E2E backend: %s", output)
	return binaryPath
}

func startE2EBackend(t *testing.T, binaryPath string, subject *e2eSubject) *e2eBackendProcess {
	t.Helper()

	port := e2eFreePort(t)
	backend := &e2eBackendProcess{baseURL: "http://127.0.0.1:" + strconv.Itoa(port)}
	command := exec.Command(binaryPath)
	command.Dir = e2eProjectRoot(t)
	command.Env = e2eBackendEnvironment(port, t.TempDir())
	command.Stdout = &backend.stdout
	command.Stderr = &backend.stderr
	require.NoError(t, command.Start())
	backend.command = command

	t.Cleanup(func() {
		backend.stop(t)
	})
	backend.waitForReady(t, subject)
	return backend
}

func (b *e2eBackendProcess) stop(t *testing.T) {
	t.Helper()
	if b.command == nil || b.command.Process == nil || b.command.ProcessState != nil {
		return
	}

	require.NoError(t, b.command.Process.Kill())
	_ = b.command.Wait()
}

func (b *e2eBackendProcess) waitForReady(t *testing.T, subject *e2eSubject) {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if b.command.ProcessState != nil {
			t.Fatalf("isolated E2E backend stopped before ready:\n%s\n%s", b.stdout.String(), b.stderr.String())
		}

		req, err := http.NewRequest(http.MethodGet, b.baseURL+"/api/me", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", subject.authorization(t))
		response, err := subject.client.Do(req)
		if err == nil {
			response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(250 * time.Millisecond)
	}

	t.Fatalf("isolated E2E backend did not become ready:\n%s\n%s", b.stdout.String(), b.stderr.String())
}

func e2eBackendEnvironment(port int, portraitDirectory string) []string {
	environment := make([]string, 0, len(os.Environ())+4)
	overrides := map[string]string{
		"PORT":                 strconv.Itoa(port),
		"PORTRAIT_STORAGE_DIR": portraitDirectory,
		"PUBLIC_BACKEND_URL":   "http://127.0.0.1:" + strconv.Itoa(port),
	}
	for _, value := range os.Environ() {
		name, _, found := strings.Cut(value, "=")
		if _, overridden := overrides[name]; !found || !overridden {
			environment = append(environment, value)
		}
	}
	for name, value := range overrides {
		environment = append(environment, fmt.Sprintf("%s=%s", name, value))
	}
	return environment
}

func e2eFreePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func e2eProjectRoot(t *testing.T) string {
	t.Helper()

	directory, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
			return directory
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			t.Fatal("could not find project root for isolated E2E backend")
		}
		directory = parent
	}
}
