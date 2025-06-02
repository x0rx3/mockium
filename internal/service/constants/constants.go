package constants

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
		"^\\$\\{(%s):(.*)\\}$",
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
	Query   Parameter = "query"
	Path    Parameter = "path"
	Form    Parameter = "form"
	Body    Parameter = "body"
)

const (
	// Text / HTML
	ContentTypeTextPlain      = "text/plain"
	ContentTypeTextHTML       = "text/html"
	ContentTypeTextCSS        = "text/css"
	ContentTypeTextCSV        = "text/csv"
	ContentTypeTextXML        = "text/xml"
	ContentTypeTextJavaScript = "text/javascript"

	// JSON / XML
	ContentTypeApplicationJSON        = "application/json"
	ContentTypeApplicationXML         = "application/xml"
	ContentTypeApplicationProblemJSON = "application/problem+json"
	ContentTypeApplicationProblemXML  = "application/problem+xml"

	// Forms
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeFormData       = "multipart/form-data"

	// Files
	ContentTypeOctetStream = "application/octet-stream"
	ContentTypePDF         = "application/pdf"
	ContentTypeZIP         = "application/zip"
	ContentTypeGZIP        = "application/gzip"
	ContentTypeTar         = "application/x-tar"

	// Images
	ContentTypeImageJPEG = "image/jpeg"
	ContentTypeImagePNG  = "image/png"
	ContentTypeImageGIF  = "image/gif"
	ContentTypeImageSVG  = "image/svg+xml"
	ContentTypeImageWebP = "image/webp"

	// Audio / Video
	ContentTypeAudioMPEG = "audio/mpeg"
	ContentTypeAudioOGG  = "audio/ogg"
	ContentTypeVideoMP4  = "video/mp4"
	ContentTypeVideoWEBM = "video/webm"

	// Stylesheets
	ContentTypeFontWOFF  = "font/woff"
	ContentTypeFontWOFF2 = "font/woff2"
	ContentTypeFontTTF   = "font/ttf"
	ContentTypeFontOTF   = "font/otf"
)
