/*
 Package mp : command helper functions
*/

package mp

// Gradle os specific gradle command
func Gradle(useScript bool) []string {
	var cmd []string
	if useScript {
		if IsWindows() {
			cmd = append(cmd, "cmd.exe", "/C", "gradlew.bat")
		} else {
			cmd = append(cmd, "sh", "gradlew")
		}
	} else {
		cmd = append(cmd, "gradle")
	}
	return cmd
}
