package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// LoadEnvironment autoload environment
func LoadEnvironment(defaultValues map[string]interface{}) {
	LoadEnvironmentSystem(defaultValues)
	LoadEnvironmentLocal()
	LoadEnvironmentParameter(defaultValues)
	LoadEnvironmentPrivate(defaultValues)
	MergeAllEnvironment()
}

func LoadTestEnvironment(defaultValues map[string]interface{}) {
	LoadEnvironmentSystem(defaultValues)
	LoadEnvironmentLocal()
	LoadEnvironmentParameter(defaultValues)
	LoadEnvironmentPrivate(defaultValues)
	MergeAllEnvironment()
}

// LoadEnvironmentSystem Load System Environment
func LoadEnvironmentSystem(defaultValues map[string]interface{}) {
	systemEnv := viper.New()
	systemEnv.AutomaticEnv()
	for k := range defaultValues {
		viper.Set(k, systemEnv.Get(k))
	}
}

// LoadEnvironmentLocal Load Local Environment
func LoadEnvironmentLocal() {
	// load local env
	localEnv := viper.New()
	localEnv.SetConfigType("dotenv")
	localEnv.SetConfigFile(".env")
	if err := localEnv.ReadInConfig(); nil == err {
		localEnvKeys := localEnv.AllKeys()
		for i := range localEnvKeys {
			viper.Set(localEnvKeys[i], localEnv.Get(localEnvKeys[i]))
		}
	}
}

// LoadEnvironmentParameter Load Parameter Environment
func LoadEnvironmentParameter(defaultValues map[string]interface{}) {
	// load parameter env
	paramEnv := viper.New()
	paramEnv.AllowEmptyEnv(false)
	for k := range defaultValues {
		if flagKey := strcase.ToKebab(k); nil == pflag.Lookup(flagKey) {
			pflag.String(flagKey, "", k)
		}
	}

	if os.Getenv("ENVIRONMENT_SIMULATION") != "" {
		for k := range defaultValues {
			pflag.CommandLine.Set(strcase.ToKebab(k), viper.GetString(k))
		}
	}

	pflag.Parse()
	if err := paramEnv.BindPFlags(pflag.CommandLine); nil == err {
		paramEnvKeys := paramEnv.AllKeys()
		for i := range paramEnvKeys {
			if stringValue := paramEnv.GetString(paramEnvKeys[i]); stringValue != "" {
				viper.Set(strcase.ToSnake(paramEnvKeys[i]), stringValue)
			}
		}
	}
}

// LoadEnvironmentPrivate Load Private Environment
func LoadEnvironmentPrivate(defaultValues map[string]interface{}) {
	// load default env
	for k, v := range defaultValues {
		if !viper.InConfig(k) {
			viper.SetDefault(k, v)
		}
	}
}

// MergeAllEnvironment Merge all System, Local, Parameter and Private Environment
func MergeAllEnvironment() {
	keys := viper.AllKeys()
	for i := range keys {
		stringValue := viper.GetString(keys[i])
		if stringValue == "" {
			if value := viper.Get(keys[i]); nil != value {
				stringValue = fmt.Sprintf("%v", value)
			}
		}
		os.Setenv(strings.ToUpper(keys[i]), stringValue)
	}
}

var testingEnvironmentStorage map[string]interface{}

/*
ResetTestingEnvironment

In unit test, we load default environment multiple times.

This can cause a "too many open files" error because within the LoadEnvironment() method, there are several methods to open or read files.

So, please using this function for reset environment in Unit Test.

Source:
  - https://stackoverflow.com/a/70832048
*/
func ResetTestingEnvironment() {
	setEnv := func(env map[string]interface{}) {
		// clear env first
		viper.Reset()
		os.Clearenv()

		// get testing environment
		// set to viper and os environment
		for key := range env {
			val := env[key]

			viper.Set(key, val)

			strValue := fmt.Sprintf("%v", val)
			os.Setenv(key, strValue)
		}
	}

	if len(testingEnvironmentStorage) > 0 {
		setEnv(testingEnvironmentStorage)
		return
	}

	// store new default environment data
	testingEnvironmentStorage = make(map[string]interface{}, 0)

	// config.Environment as default
	testEnv := testEnvironment()
	LoadTestEnvironment(testEnv)

	keys := viper.AllKeys()

	for i := range keys {
		keyName := keys[i]
		value := viper.Get(keyName)
		testingEnvironmentStorage[keyName] = value
	}

	setEnv(testingEnvironmentStorage)
}

func testEnvironment() (result map[string]interface{}) {
	defaultEnv := make(map[string]interface{})
	// Disable sentry log
	defaultEnv["enable_sentry_log"] = false

	result = defaultEnv
	return
}
