package converter

import "github.com/yuin/goldmark/parser"

var (
	readerKey      = parser.NewContextKey()
	pathKey        = parser.NewContextKey()
	errorKey       = parser.NewContextKey()
	attachmentsKey = parser.NewContextKey()
)

func newContext(reader Reader, path string) parser.Context {
	ctx := parser.NewContext()
	ctx.Set(readerKey, reader)
	ctx.Set(pathKey, path)
	return ctx
}

func getReader(pc parser.Context) Reader {
	v := pc.Get(readerKey)
	if v == nil {
		return nil
	}
	return v.(Reader)
}

func getPath(pc parser.Context) string {
	v := pc.Get(pathKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

func getError(pc parser.Context) error {
	v := pc.Get(errorKey)
	if v == nil {
		return nil
	}
	return v.(error)
}

func getAttachments(pc parser.Context) []Attachment {
	v := pc.Get(attachmentsKey)
	if v == nil {
		return nil
	}
	return v.([]Attachment)
}

func addAttachment(pc parser.Context, attachment Attachment) {
	attachments := getAttachments(pc)
	pc.Set(attachmentsKey, append(attachments, attachment))
}
