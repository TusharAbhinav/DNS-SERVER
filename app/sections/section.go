package section

import (
	"github.com/codecrafters-io/dns-server-starter-go/app/sections/header"
	"github.com/codecrafters-io/dns-server-starter-go/app/sections/question"
)

type Section struct {
	Header   header.Header
	Question question.Question
}
