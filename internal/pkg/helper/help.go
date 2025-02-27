package helper

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func Caller(i int) string {
	_, file, line, ok := runtime.Caller(i)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}

/**
 *
 * StackTrace returns a formatted stack trace as a string.
 * The stack trace starts from the caller specified by the skip parameter
 * and includes up to the number of frames specified by the depth parameter.
 * Parameters:
 * 	- skip: The number of stack frames to skip before recording the stack trace.
 * 	- depth: The maximum number of stack frames to record.
 *
 * Returns:
 * A string representing the formatted stack trace.
 */
func StackTrace(skip int, depth int) string {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	var stack strings.Builder
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		fmt.Fprintf(&stack, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
	}
	return stack.String()
}

func GetPort(port string) int {
	re := regexp.MustCompile("[0-9]+")
	port = re.FindString(port)
	portNo, err := strconv.Atoi(port)
	if err != nil {
		log.Panic(err)
		return 0
	}
	return portNo
}

// GetMachineID returns a unique identifier for the machine/container.
func GetMachineID() string {
	// Check if running inside a container
	if containerID := getContainerID(); containerID != "" {
		return containerID
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("cat", "/etc/machine-id")
	case "darwin": // macOS
		cmd = exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	case "windows":
		cmd = exec.Command("wmic", "csproduct", "get", "UUID")
	default:
		return uuid.NewString()
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return uuid.NewString()
	}

	machineID := strings.TrimSpace(out.String())

	// macOS: Extract the machine ID
	if runtime.GOOS == "darwin" {
		lines := strings.Split(machineID, "\n")
		for _, line := range lines {
			if strings.Contains(line, "IOPlatformUUID") {
				parts := strings.Fields(line)
				return parts[len(parts)-1]
			}
		}
	}

	// Windows: Extract UUID from output
	if runtime.GOOS == "windows" {
		lines := strings.Split(machineID, "\n")
		if len(lines) > 1 {
			return strings.TrimSpace(lines[1])
		}
	}

	return machineID
}

// getContainerID attempts to retrieve a container ID if running inside a container.
func getContainerID() string {
	// Check /etc/hostname (commonly used for container ID)
	if data, err := os.ReadFile("/etc/hostname"); err == nil {
		containerID := strings.TrimSpace(string(data))
		if len(containerID) > 0 {
			return containerID
		}
	}

	// Check /proc/self/cgroup (used in Docker & Kubernetes)
	if data, err := os.ReadFile("/proc/self/cgroup"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) == 3 && strings.Contains(parts[2], "/docker/") {
				return strings.TrimPrefix(parts[2], "/docker/")
			} else if len(parts) == 3 && strings.Contains(parts[2], "/kubepods/") {
				return strings.Split(parts[2], "/")[len(parts[2])-1]
			}
		}
	}

	return ""
}
