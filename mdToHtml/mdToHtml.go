package mdtohtml

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func Parse(md string) string {
	unsafeHtmlContent := blackfriday.Run([]byte(md))
	saveHtmlContent := bluemonday.UGCPolicy().SanitizeBytes(unsafeHtmlContent)

	return string(saveHtmlContent)
}
