package env

import (
	godotenv "github.com/joho/godotenv"
	envconfig "github.com/kelseyhightower/envconfig"
)

/*
Load will read your env file(s) and load them into environment variables for this process.
Call this function as close as possible to the start of your program (ideally in main).
If you call Load without any args it will default to loading .env in the current path.
You can otherwise tell it which files to load (there can be more than one) like:
godotenv.Load("fileone", "filetwo")
This way, fileone will override value in filetwo.

It's important to note that it WILL NOT OVERRIDE an env variable that already exists.
Consider the .env file to set development variables or sensible defaults.
*/
func Load(filenames ...string) error {
	return godotenv.Load(filenames...)
}

// Parse parses environment variables.
func Parse(prefix string, out interface{}) error {
	return envconfig.Process(prefix, out)
}

// LoadAndParse loads and parses environment variables.
func LoadAndParse(prefix string, out interface{}, filenames ...string) error {
	if err := Load(filenames...); err != nil {
		return err
	}

	if err := Parse(prefix, out); err != nil {
		return err
	}

	return nil
}
