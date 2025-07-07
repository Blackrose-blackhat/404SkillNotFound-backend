package parser

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func ExtractResume(file io.Reader) (string, error) {
	tmpFile, err := os.CreateTemp("", "resume-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file:%w", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		return "", fmt.Errorf("Failed to Copy file %w", err)

	}
	tmpFile.Close()

	cmd := exec.Command("pdftotext", tmpFile.Name(), "-")
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}
	return string(output), nil
}
