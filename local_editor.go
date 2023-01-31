//go:build !linux && !windows && !darwin
// +build !linux,!windows,!darwin

package launch_editor

func getLocalEditor() string {
	return ""
}
