package mp

import (
	"reflect"
	"testing"
)

func TestBuildAngular(t *testing.T) {
	type args struct {
		args    []string
		profile string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"ng build with profile and additional args",
			args{
				args:    []string{"--base-href", "/kann/das/sein"},
				profile: "prod",
			},
			[]string{"ng", "build", "--prod", "--base-href", "/kann/das/sein"},
		},
		{
			"ng build without profile",
			args{
				args:    []string{},
				profile: "",
			},
			[]string{"ng", "build", "--prod"},
		},
		{
			"ng build with non default profile",
			args{
				args:    []string{},
				profile: "herbert",
			},
			[]string{"ng", "build", "--herbert"},
		},
		{
			"ng build with non default profile and additional args",
			args{
				args:    []string{"--arg"},
				profile: "herbert",
			},
			[]string{"ng", "build", "--herbert", "--arg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildAngular(tt.args.args, tt.args.profile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildAngular() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildNpm(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"npm build command without arguments",
			args{},
			[]string{"npm", "run", "build"},
		},
		{
			"npm build command with arguments",
			args{[]string{"-o", "output.file"}},
			[]string{"npm", "run", "build", "-o", "output.file"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildNpm(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildNpm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildGradle(t *testing.T) {
	type args struct {
		args            []string
		useNativeGradle bool
		ignoreTest      bool
		isWindows       bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"run native gradle build with tests",
			args{
				[]string{},
				true,
				false,
				false,
			},
			[]string{"gradle", "clean", "build"},
		},
		{
			"run native gradle build without test",
			args{
				[]string{},
				true,
				true,
				false,
			},
			[]string{"gradle", "clean", "build", "-x", "test"},
		},
		{
			"run gradle script for linux with tests",
			args{
				[]string{},
				false,
				false,
				false,
			},
			[]string{"sh", "gradlew", "clean", "build"},
		},
		{
			"run gradle script for linux without tests",
			args{
				[]string{},
				false,
				true,
				false,
			},
			[]string{"sh", "gradlew", "clean", "build", "-x", "test"},
		},
		{
			"run gradle script for linux with tests and arguments",
			args{
				[]string{"-Pprofile=dev"},
				false,
				false,
				false,
			},
			[]string{"sh", "gradlew", "clean", "build", "-Pprofile=dev"},
		},
		{
			"run native gradle with tests and arguments",
			args{
				[]string{"-Pprofile=dev"},
				true,
				false,
				false,
			},
			[]string{"gradle", "clean", "build", "-Pprofile=dev"},
		},
		{
			"run gradle script for windows with tests",
			args{
				[]string{},
				false,
				false,
				true,
			},
			[]string{"cmd.exe", "/C", "gradlew.bat", "clean", "build"},
		},
		{
			"run gradle script for windows without tests and arguments",
			args{
				[]string{"-Pprofile=dev"},
				false,
				true,
				true,
			},
			[]string{"cmd.exe", "/C", "gradlew.bat", "clean", "build", "-x", "test", "-Pprofile=dev"},
		},
		{
			"run gradle script for windows with tests and arguments",
			args{
				[]string{"-Pprofile=dev"},
				false,
				false,
				true,
			},
			[]string{"cmd.exe", "/C", "gradlew.bat", "clean", "build", "-Pprofile=dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildGradle(tt.args.args, tt.args.useNativeGradle, tt.args.ignoreTest, tt.args.isWindows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildGradle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildGo(t *testing.T) {
	type args struct {
		args        []string
		mugfileArgs []string
		target      string
	}
	tests := []struct {
		name    string
		args    args
		wantCmd []string
		wantEnv []string
	}{
		{
			"go build without specifying target",
			args{},
			[]string{"go", "build"},
			[]string{},
		},
		{
			"go build without specifying target but with arguments",
			args{
				[]string{"--arg"},
				[]string{},
				"",
			},
			[]string{"go", "build", "--arg"},
			[]string{},
		},
		{
			"go build without specifying target but with arguments and mugfile arguments",
			args{
				[]string{"--arg"},
				[]string{"--from-file", "13"},
				"",
			},
			[]string{"go", "build", "--from-file", "13", "--arg"},
			[]string{},
		},
		{
			"go build with target linux and with arguments and mugfile arguments",
			args{
				[]string{"--arg"},
				[]string{"--from-file", "13"},
				"linux",
			},
			[]string{"go", "build", "--from-file", "13", "--arg"},
			[]string{"GOOS=linux", "GOARCH=amd64"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCmd, gotEnv := BuildGo(tt.args.args, tt.args.mugfileArgs, tt.args.target); !reflect.DeepEqual(gotCmd, tt.wantCmd) || (len(gotEnv) > 0 && !reflect.DeepEqual(gotEnv, tt.wantEnv)) {
				t.Errorf("BuildGo(Command) = %v, want %v", gotCmd, tt.wantCmd)
				t.Errorf("BuildGo(Environment) = %v, want %v", gotEnv, tt.wantEnv)
			}
		})
	}
}

func TestBuildDocker(t *testing.T) {
	type args struct {
		image string
		tags  []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"docker build with default tag",
			args{
				image: "image",
				tags:  []string{},
			},
			[]string{"docker", "build", "--tag", "image", "."},
		},
		{
			"docker build with a single tag",
			args{
				image: "image",
				tags:  []string{"tag"},
			},
			[]string{"docker", "build", "--tag", "image:tag", "."},
		},
		{
			"docker build with a multiple tags",
			args{
				image: "image",
				tags:  []string{"tag", "gat"},
			},
			[]string{"docker", "build", "--tag", "image:tag", "--tag", "image:gat", "."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildDocker(tt.args.image, tt.args.tags, ".", []string{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildDocker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunAngular(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"ng serve without args",
			args{},
			[]string{"ng", "serve"},
		},
		{
			"ng serve with arguments",
			args{
				[]string{"--base-href", "/kern"},
			},
			[]string{"ng", "serve", "--base-href", "/kern"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunAngular(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunAngular() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunNpm(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"npm start without args",
			args{},
			[]string{"npm", "start"},
		},
		{
			"ng serve with arguments",
			args{
				[]string{"--base-href", "/kern"},
			},
			[]string{"npm", "start", "--", "--base-href", "/kern"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunNpm(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunNpm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunGradle(t *testing.T) {
	type args struct {
		args            []string
		useNativeGradle bool
		profile         string
		isWindows       bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"run native gradle bootRun without profile",
			args{
				[]string{},
				true,
				"",
				false,
			},
			[]string{"gradle", "bootRun"},
		},
		{
			"run native gradle bootRun with profile",
			args{
				[]string{},
				true,
				"heribert",
				false,
			},
			[]string{"gradle", "bootRun", "-Pprofile=heribert"},
		},
		{
			"run native gradle bootRun with profile and arguments",
			args{
				[]string{"-Dsome=JAVA"},
				true,
				"heribert",
				false,
			},
			[]string{"gradle", "bootRun", "-Pprofile=heribert", "-Dsome=JAVA"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunGradle(tt.args.args, tt.args.useNativeGradle, tt.args.profile, tt.args.isWindows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunGradle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunGo(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"go run without arguments",
			args{},
			[]string{"go", "run", "main.go"},
		},
		{
			"go run with arguments",
			args{
				[]string{"--some", "arg"},
			},
			[]string{"go", "run", "main.go", "--some", "arg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunGo(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunGo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrdCommit(t *testing.T) {
	type args struct {
		args         []string
		branch       string
		addAll       bool
		overrideType string
	}
	tests := []struct {
		name       string
		args       args
		wantAdd    []string
		wantCommit []string
	}{
		{
			"frd commit without adding files with feature commit",
			args{
				[]string{"kann das sein"},
				"f-333-abc",
				false,
				"",
			},
			[]string{},
			[]string{"git", "commit", "-m", "[FEATURE][FRD-333] kann das sein"},
		},
		{
			"frd commit with adding files with refactor commit",
			args{
				[]string{"kann das sein"},
				"r-333-abc",
				true,
				"",
			},
			[]string{"git", "add", "."},
			[]string{"git", "commit", "-m", "[REFACTORING][FRD-333] kann das sein"},
		},
		{
			"frd commit with adding files with internal commit that is overwritten with style",
			args{
				[]string{"kann das sein"},
				"i-333-abc",
				true,
				"s",
			},
			[]string{"git", "add", "."},
			[]string{"git", "commit", "-m", "[STYLE][FRD-333] kann das sein"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAdd, gotCommit := FrdCommit(tt.args.args, tt.args.branch, tt.args.addAll, tt.args.overrideType); !reflect.DeepEqual(gotCommit, tt.wantCommit) || (len(gotAdd) > 0 && !reflect.DeepEqual(gotAdd, tt.wantAdd)) {
				t.Errorf("FrdCommit() = %v, want %v", gotAdd, tt.wantAdd)
				t.Errorf("FrdCommit() = %v, want %v", gotCommit, tt.wantCommit)
			}
		})
	}
}

func TestLogDocker(t *testing.T) {
	type args struct {
		args        []string
		containerID string
		limit       string
		follow      bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"docker log without following and default log limit",
			args{
				[]string{},
				"conti",
				"10",
				false,
			},
			[]string{"docker", "logs", "conti", "--tail", "10"},
		},
		{
			"docker log without following and specified log limit",
			args{
				[]string{},
				"conti",
				"333",
				false,
			},
			[]string{"docker", "logs", "conti", "--tail", "333"},
		},
		{
			"docker log with following and default log limit",
			args{
				[]string{},
				"conti",
				"10",
				true,
			},
			[]string{"docker", "logs", "conti", "--tail", "10", "-f"},
		},
		{
			"docker log with following and specified log limit and arguments",
			args{
				[]string{"-since", "never"},
				"conti",
				"654",
				true,
			},
			[]string{"docker", "logs", "conti", "--tail", "654", "-f", "-since", "never"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LogDocker(tt.args.args, tt.args.containerID, tt.args.limit, tt.args.follow); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogDocker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogGit(t *testing.T) {
	type args struct {
		args      []string
		format    string
		limit     string
		fileNames bool
		graph     bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"git log with default limit",
			args{
				[]string{},
				"ABC",
				"10",
				false,
				true,
			},
			[]string{"git", "log", "--format=ABC", "-10", "--graph"},
		},
		{
			"git log with specified log limit",
			args{
				[]string{},
				"ABC",
				"33",
				false,
				true,
			},
			[]string{"git", "log", "--format=ABC", "-33", "--graph"},
		},
		{
			"git log with file names",
			args{
				[]string{},
				"ABC",
				"10",
				true,
				true,
			},
			[]string{"git", "log", "--format=ABC", "-10", "--graph", "--name-only"},
		},
		{
			"git log with default limit, file names and arguments",
			args{
				[]string{"--oneline"},
				"ABC",
				"10",
				true,
				true,
			},
			[]string{"git", "log", "--format=ABC", "-10", "--graph", "--name-only", "--oneline"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LogGit(tt.args.args, tt.args.format, tt.args.limit, tt.args.fileNames, tt.args.graph); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogGit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeployDocker(t *testing.T) {
	type args struct {
		image string
		tag   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"simple docker deploy command",
			args{
				"hans/dampf",
				"a.b",
			},
			[]string{"docker", "push", "hans/dampf:a.b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeployDocker(tt.args.image, tt.args.tag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeployDocker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditConfigFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"open config file with micro",
			args{
				"/home/coole/config",
			},
			[]string{"micro", "/home/coole/config"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EditConfigFile(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EditConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
