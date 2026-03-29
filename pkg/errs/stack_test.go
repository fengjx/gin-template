package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() *Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() *Frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		frame  *Frame
		format string
		assert func(t *testing.T, got string)
	}{{
		initpc,
		"%s",
		func(t *testing.T, got string) {
			if got != "stack_test.go" {
				t.Fatalf("unexpected %%s output: %q", got)
			}
		},
	}, {
		initpc,
		"%+s",
		func(t *testing.T, got string) {
			expectContains(t, got, initpc.Function)
			expectContains(t, got, filepath.Base(initpc.File))
		},
	}, {
		initpc,
		"%d",
		func(t *testing.T, got string) {
			if got != fmt.Sprintf("%d", initpc.Frame.Line) {
				t.Fatalf("unexpected %%d output: %q", got)
			}
		},
	}, {
		initpc,
		"%n",
		func(t *testing.T, got string) {
			if got != "init" {
				t.Fatalf("unexpected %%n output: %q", got)
			}
		},
	}, {
		func() *Frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		func(t *testing.T, got string) {
			expectContains(t, got, "(*X).ptr")
		},
	}, {
		func() *Frame {
			var x X
			return x.val()
		}(),
		"%n",
		func(t *testing.T, got string) {
			expectContains(t, got, "X.val")
		},
	}, {
		initpc,
		"%v",
		func(t *testing.T, got string) {
			if got != fmt.Sprintf("%s:%d", filepath.Base(initpc.File), initpc.Frame.Line) {
				t.Fatalf("unexpected %%v output: %q", got)
			}
		},
	}, {
		initpc,
		"%+v",
		func(t *testing.T, got string) {
			expectContains(t, got, initpc.Function)
			expectContains(t, got, fmt.Sprintf("%s:%d", filepath.Base(initpc.File), initpc.Frame.Line))
		},
	}}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			tt.assert(t, fmt.Sprintf(tt.format, tt.frame))
		})
	}
}

func TestFuncname(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"gin-template/pkg/errs.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}

	for _, tt := range tests {
		got := funcname(tt.name)
		want := tt.want
		if got != want {
			t.Errorf("funcname(%q): want: %q, got %q", tt.name, want, got)
		}
	}
}

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err error
	}{{
		err: Wrap(errors.New("ooh"), "ahh"),
	}}
	for i, tt := range tests {
		x, ok := tt.err.(interface {
			StackTrace() StackTrace
		})
		if !ok {
			t.Errorf("expected %#v to implement StackTrace() StackTrace", tt.err)
			continue
		}
		st := x.StackTrace()
		if len(st) == 0 {
			t.Fatalf("test %d: expected stack trace", i+1)
		}
		got := fmt.Sprintf("%+v", st[0])
		expectContains(t, got, "TestStackTrace")
		expectContains(t, got, "stack_test.go")
	}
}

func stackTrace() StackTrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	var st Stack = pcs[0:n]
	return st.StackTrace()
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		StackTrace
		format string
		assert func(t *testing.T, got string)
	}{{
		nil,
		"%s",
		func(t *testing.T, got string) {
			if got != "[]" {
				t.Fatalf("unexpected %%s output: %q", got)
			}
		},
	}, {
		nil,
		"%v",
		func(t *testing.T, got string) {
			if got != "[]" {
				t.Fatalf("unexpected %%v output: %q", got)
			}
		},
	}, {
		nil,
		"%+v",
		func(t *testing.T, got string) {
			if got != "" {
				t.Fatalf("unexpected %%+v output: %q", got)
			}
		},
	}, {
		nil,
		"%#v",
		func(t *testing.T, got string) {
			if got != "errs.StackTrace{}" {
				t.Fatalf("unexpected %%#v output: %q", got)
			}
		},
	}, {
		make(StackTrace, 0),
		"%s",
		func(t *testing.T, got string) {
			if got != "[]" {
				t.Fatalf("unexpected %%s output: %q", got)
			}
		},
	}, {
		make(StackTrace, 0),
		"%v",
		func(t *testing.T, got string) {
			if got != "[]" {
				t.Fatalf("unexpected %%v output: %q", got)
			}
		},
	}, {
		make(StackTrace, 0),
		"%+v",
		func(t *testing.T, got string) {
			if got != "" {
				t.Fatalf("unexpected %%+v output: %q", got)
			}
		},
	}, {
		make(StackTrace, 0),
		"%#v",
		func(t *testing.T, got string) {
			if got != "errs.StackTrace{}" {
				t.Fatalf("unexpected %%#v output: %q", got)
			}
		},
	}, {
		stackTrace()[:2],
		"%s",
		func(t *testing.T, got string) {
			expectContains(t, got, "stack_test.go")
		},
	}, {
		stackTrace()[:2],
		"%v",
		func(t *testing.T, got string) {
			expectContains(t, got, "stack_test.go:")
		},
	}, {
		stackTrace()[:2],
		"%+v",
		func(t *testing.T, got string) {
			expectContains(t, got, "stackTrace")
			expectContains(t, got, "TestStackTraceFormat")
			expectContains(t, got, "stack_test.go")
		},
	}, {
		stackTrace()[:2],
		"%#v",
		func(t *testing.T, got string) {
			expectContains(t, got, "errs.StackTrace{")
			expectContains(t, got, `Function:"`)
			expectContains(t, got, `File:"`)
		},
	}}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			tt.assert(t, fmt.Sprintf(tt.format, tt.StackTrace))
		})
	}
}

// a version of runtime.Caller that returns a Frame, not a uintptr.
func caller() *Frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return buildFrame(frame)
}

func expectContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("got %q, want contains %q", got, want)
	}
}
