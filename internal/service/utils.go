package service

import (
	"fmt"
	"regexp"
)

// RegexpResponseValuePlaceholder is a regular expression that matches a specific format for response value placeholders.
// The format is: ${req.<param_type>:<param_name>}
// where <param_type> can be one of the following: headers, queryParams, pathParams, formParams, body
// and <param_name> can be any alphanumeric string, underscore, or hyphen.
var RegexpResponseValuePlaceholder = regexp.MustCompile(
	fmt.Sprintf("^\\$\\{(req)\\.(%s|%s|%s|%s|%s):([a-zA-Z0-9_-]+|\\*)\\}$",
		Headers,
		Query,
		Path,
		Form,
		Body,
	),
)

// RegexpRequestValuePlaceholder is a regular expression that matches a specific format for request value placeholders.
// The format is: ${regexp:<param_name>}
// where <param_name> can be any alphanumeric string, underscore, or hyphen.
// The placeholder is used to indicate that the value should be treated as a regular expression.
var RegexpRequestValuePlaceholder = regexp.MustCompile(
	fmt.Sprintf(
		"^\\$\\{(%s):([a-zA-Z0-9_-]+|\\*)\\}$",
		RegexpValuePlaceholder,
	),
)

const (
	// regexpValuePlaceholder is a constant string used to identify the type of value in a placeholder.
	RegexpValuePlaceholder = "regexp"
)

const (
	// anyValuePlaceholder is a constant string used to identify the type of value in a placeholder.
	// It indicates that the value can be of any type.
	AnyValuePlaceholder = "${...}"
	// fileParamName is a constant string used to identify the type of value in a placeholder.
	// It indicates that the value is a file parameter.
	FileParamName = "${file}"
)

type Parameter string

// The following constants represent the different types of placeholders that can be used.
const (
	Headers Parameter = "headers"
	Query   Parameter = "queryParams"
	Path    Parameter = "pathParams"
	Form    Parameter = "formParams"
	Body    Parameter = "body"
)
