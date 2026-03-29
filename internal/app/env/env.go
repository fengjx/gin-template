package env

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/spf13/pflag"
)

type Environment string

const (
	Dev  Environment = "dev"
	Test Environment = "test"
	Prod Environment = "prod"
)

var (
	once         sync.Once
	current      Environment
	flagValue    string
	argsProvider = func() []string { return os.Args[1:] }
	lookupEnv    = os.LookupEnv
)

func BindFlags(flags *pflag.FlagSet) {
	if flags == nil || flags.Lookup("env") != nil {
		return
	}
	flags.StringVar(&flagValue, "env", "", "运行环境")
}

func Current() Environment {
	once.Do(func() {
		current = resolve()
	})
	return current
}

func IsDev() bool {
	return Current() == Dev
}

func IsTest() bool {
	return Current() == Test
}

func IsProd() bool {
	return Current() == Prod
}

func ResetForTest() {
	once = sync.Once{}
	current = ""
	flagValue = ""
	argsProvider = func() []string { return os.Args[1:] }
	lookupEnv = os.LookupEnv
}

func resolve() Environment {
	candidate := string(Dev)
	if value, ok := lookupEnv("APP_ENV"); ok && strings.TrimSpace(value) != "" {
		candidate = value
	}
	if strings.TrimSpace(flagValue) != "" {
		candidate = flagValue
	}

	flags := pflag.NewFlagSet("env", pflag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.String("env", candidate, "运行环境")
	if err := flags.Parse(argsProvider()); err == nil {
		if value, getErr := flags.GetString("env"); getErr == nil {
			candidate = value
		}
	}

	switch strings.ToLower(strings.TrimSpace(candidate)) {
	case string(Test):
		return Test
	case string(Prod):
		return Prod
	default:
		return Dev
	}
}
