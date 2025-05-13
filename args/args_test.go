package args

import ("testing"
"os")

func TestParseArgs(t *testing.T) {
	cases := []struct {
		name     string
		args []string
		want *Args
		wantErr bool
	} {
		{"all args", []string{"-i", "input.txt", "-o", "output.txt", "-d"}, &Args{"input.txt", "output.txt", true}, false},
		{"no args", []string{}, nil, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			os.Args = c.args;
			got, err := ParseArgs()
			// Check if the error matches the expected error
			if c.wantErr && (err != nil) {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			// If no error is expected, check if the returned args match the expected args
			if err == nil && *got != *c.want {
				t.Errorf("ParseArgs() = %v, want %v", got, c.want)
			}
		})
	}
}
