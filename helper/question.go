package helper
import (
	"io"
	"os"
	"fmt"
	"bufio"
	"strings"
)

type Question struct {
	Msg string
	Writer io.Writer
	Reader io.Reader
}

func NewQuestion(msg string) *Question {
	return &Question{
		Msg:	msg,
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
}

func (q *Question) Ask() string {
	fmt.Fprint(q.Writer, q.Msg)

	reader := bufio.NewReader(q.Reader)
	str, _ := reader.ReadString('\n')
	str = strings.Trim(str, "\r\n")

	return str
}

