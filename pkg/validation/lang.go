package validation

var errorMessages = map[string]string{
	"_":       "The `{field}` field did not pass validation",
	"_json":   "Invalid data was transmitted that can't be unmarshalled",
	"_noData": "It is necessary to transfer data for the request",
	// builtin
	"_validate": "The `{field}` field did not pass validation", // default validate message
	"_filter":   "The `{field}` field data is invalid",         // data filter error
	// int value
	"min": "The `{field}` field min value is %v",
	"max": "The `{field}` field max value is %v",
	// type check: float
	"isFloat": "The `{field}` field value must be an integer",
	// type check: int
	"isInt":   "The `{field}` field value must be an integer",
	"intType": "The `{field}` field value must be an integer",
	"isInt1":  "The `{field}` field value must be an integer and mix value is %d",      // has min check
	"isInt2":  "The `{field}` field value must be an integer and in the range %d - %d", // has min, max check
	"isInts":  "The `{field}` field value must be an int slice",
	"isUint":  "The `{field}` field value must be an unsigned integer(>= 0)",
	// type check: string
	"isString":  "The `{field}` field value must be a string",
	"isString1": "The `{field}` field value must be a string and min length is %d", // has min len check
	// length
	"minLength": "The `{field}` field min length is %d",
	"maxLength": "The `{field}` field max length is %d",
	// string length. calc rune
	"stringLength":  "The `{field}` field length must be in the range %d - %d",
	"stringLength1": "The `{field}` field min length is %d",
	"stringLength2": "The `{field}` field length must be in the range %d - %d",

	"isURL":     "The `{field}` field must be a valid URL address",
	"isFullURL": "The `{field}` field must be a valid full URL address",
	"regexp":    "The `{field}` field must match pattern %s",

	"isFile":  "The `{field}` field must be an uploaded file",
	"isImage": "The `{field}` field must be an uploaded image file",

	"enum":  "The `{field}` field value must be in the enum %v",
	"range": "The `{field}` field value must be in the range %d - %d",
	// int compare
	"lt": "The `{field}` field value should be less than %v",
	"gt": "The `{field}` field value should be greater than %v",
	// required
	"required":           "The `{field}` field is required to not be empty",
	"requiredIf":         "The `{field}` field is required when {args0} is {args1end}",
	"requiredUnless":     "The `{field}` field field is required unless {args0} is in {args1end}",
	"requiredWith":       "The `{field}` field field is required when {values} is present",
	"requiredWithAll":    "The `{field}` field field is required when {values} is present",
	"requiredWithout":    "The `{field}` field field is required when {values} is not present",
	"requiredWithoutAll": "The `{field}` field field is required when none of {values} are present",
	// field compare
	"eqField":  "The `{field}` field value must be equal the field %s",
	"neField":  "The `{field}` field value cannot be equal to the field %s",
	"ltField":  "The `{field}` field value should be less than the field %s",
	"lteField": "The `{field}` field value should be less than or equal to the field %s",
	"gtField":  "The `{field}` field value must be greater than the field %s",
	"gteField": "The `{field}` field value should be greater or equal to the field %s",
	// data type
	"bool":    "The `{field}` field value must be a bool",
	"float":   "The `{field}` field value must be a float",
	"slice":   "The `{field}` field value must be a slice",
	"map":     "The `{field}` field value must be a map",
	"array":   "The `{field}` field value must be an array",
	"strings": "The `{field}` field value must be a []string",
	"notIn":   "The `{field}` field value must not be in the given enum list %d",
	//
	"contains":    "The `{field}` field value does not contain %s",
	"notContains": "The `{field}` field value contains %s",
	"startsWith":  "The `{field}` field value does not start with %s",
	"endsWith":    "The `{field}` field value does not end with %s",
	"email":       "The `{field}` field value is an invalid email address",
	"regex":       "The `{field}` field value does not pass the regex check",
	"file":        "The `{field}` field value must be a file",
	"image":       "The `{field}` field value must be an image",
	// date
	"date":    "The `{field}` field value should be a date string",
	"gtDate":  "The `{field}` field value should be after %s",
	"ltDate":  "The `{field}` field value should be before %s",
	"gteDate": "The `{field}` field value should be after or equal to %s",
	"lteDate": "The `{field}` field value should be before or equal to %s",

	"hasWhitespace":  "The `{field}` field value should contains spaces",
	"ascii":          "The `{field}` field value should be an ASCII string",
	"alpha":          "The `{field}` field value contains only alpha char",
	"alphaNum":       "The `{field}` field value contains only alpha char and num",
	"alphaDash":      "The `{field}` field value contains only letters, num, dashes (-) and underscores (_)",
	"multiByte":      "The `{field}` field value should be a multiByte string",
	"base64":         "The `{field}` field value should be a base64 string",
	"dnsName":        "The `{field}` field value should be a DNS string",
	"dataURI":        "The `{field}` field value should be a DataURL string",
	"empty":          "The `{field}` field value should be empty",
	"hexColor":       "The `{field}` field value should be a color string in hexadecimal",
	"hexadecimal":    "The `{field}` field value should be a hexadecimal string",
	"json":           "The `{field}` field value should be a json string",
	"lat":            "The `{field}` field value should be a latitude coordinate",
	"lon":            "The `{field}` field value should be a longitude coordinate",
	"num":            "The `{field}` field value should be a num (>=0) string",
	"mac":            "The `{field}` field value should be a MAC address",
	"cnMobile":       "The `{field}` field value should be string of Chinese 11-digit mobile phone numbers",
	"printableASCII": "The `{field}` field value should be a printable ASCII string",
	"rgbColor":       "The `{field}` field value should be a RGB color string",
	"fullURL":        "The `{field}` field value should be a complete URL string",
	"full":           "The `{field}` field value should be a URL string",
	"ip":             "The `{field}` field value should be an IP (v4 or v6) string",
	"ipv4":           "The `{field}` field value should be an IPv4 string",
	"ipv6":           "The `{field}` field value should be an IPv6 string",
	"CIDR":           "The `{field}` field value should be a CIDR string",
	"CIDRv4":         "The `{field}` field value should be a CIDRv4 string",
	"CIDRv6":         "The `{field}` field value should be a CIDRv6 string",
	"uuid":           "The `{field}` field value should be a UUID string",
	"uuid3":          "The `{field}` field value should be a UUID3 string",
	"uuid4":          "The `{field}` field value should be a UUID4 string",
	"uuid5":          "The `{field}` field value should be a UUID5 string",
	"filePath":       "The `{field}` field value should be an existing file path",
	"unixPath":       "The `{field}` field value should be a unix path string",
	"winPath":        "The `{field}` field value should be a windows path string",
	"isbn10":         "The `{field}` field value should be a isbn10 string",
	"isbn13":         "The `{field}` field value should be a isbn13 string",
}
