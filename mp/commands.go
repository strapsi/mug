package mp

func Gradle(useScript bool) []string {
	cmd := []string{}
	if useScript {
		if IsWindows() {
			cmd  = append(cmd, "cmd.exe", "/C", "gradlew.bat")
		} else {
			cmd  = append(cmd, "sh", "gradlew")
		}
	} else {
		cmd = append(cmd, "gradle")
	}
	return cmd	
}
