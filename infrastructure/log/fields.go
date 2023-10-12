package log

import (
	"runtime"
	"strconv"
	"strings"
)

type Fields map[string]interface{}

func (f Fields) With(key string, value interface{}) Fields {
	f[key] = value
	return f
}

func (f Fields) WithFields(fields Fields) Fields {
	for key, value := range fields {
		f[key] = value
	}
	return f
}

func NewFields() Fields {
	fields := map[string]interface{}{}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	path := strings.Builder{}
	path.WriteString(file)
	path.WriteString(":")
	path.WriteString(strconv.Itoa(line))

	if f := runtime.FuncForPC(pc); f == nil {
		fields["func_name"] = "unknown"
		fields["file_name"] = formatRelativePath(path.String())
	} else {
		fields["func_name"] = formatCaller(f.Name())
		fields["file_name"] = formatRelativePath(path.String())
	}

	return fields
}

func NewFieldsWithMessage(message string) Fields {
	fields := NewFields()
	fields["message"] = message
	return fields
}

func NewFieldsWithError(err error) Fields {
	fields := NewFields()
	fields["error"] = err
	return fields
}

func formatRelativePath(path string) string {
	return strings.TrimPrefix(path, workingDir)
}

func formatCaller(caller string) string {
	return strings.TrimPrefix(caller, pkgName)
}
