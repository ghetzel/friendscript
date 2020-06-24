package utils

type RunOptions struct {
	// If true, the scope of the running script will not inherit values from the calling scope.
	Isolated bool `json:"isolated"`

	// Specifies a key in the scope of the evaluate script that will be used as the result value of this command.
	ResultKey string `json:"result"`

	// Provides a set of initial variables to the script.
	Data map[string]interface{} `json:"data"`

	// Sets the base path from which relative file lookups will be performed
	BasePath string `json:"-"`
}
