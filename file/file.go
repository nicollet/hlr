package file

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

//func main() {
//	rf := NewReadFile()
//	text := rf.readFile("test")
//
//	fmt.Print(text)
//
//	l := lex()
//	for {
//		item, more := <-l.items
//		if !more {
//			break
//		}
//		fmt.Printf("item: %v\n", item)
//	}
//}

type ReadFile struct {
	doneFiles     map[string]bool
	includeRegexp *regexp.Regexp
}

func New() ReadFile {
	return ReadFile{
		doneFiles:     make(map[string]bool),
		includeRegexp: regexp.MustCompile(`^\s*include\s+"(.*?)"`),
	}
}

func (rf ReadFile) Read(file_name string) string {
	var text string
	if file_name[0] == '/' || file_name[0] == '.' {
		return rf.readOneFile(file_name, text)
	}

	for _, p := range []string{"hosts", "includes"} {
		name := p + "/" + file_name
		names := []string{name}
		name_strip := stripnb(name)
		if name_strip != name {
			names = append(names, name_strip)
		}
		for _, file_name := range names {
			if _, err := os.Stat(file_name); os.IsNotExist(err) {
				continue
			}
			_, hask := rf.doneFiles[file_name]
			if hask {
				continue
			}
			new_text := rf.readOneFile(file_name, text)
			text += new_text
		}
	}
	return text
}

func stripnb(hostname string) string {
	return strings.Trim(hostname, "0123456789")
}

func (rf ReadFile) readOneFile(file_name string, text string) string {
	rf.doneFiles[file_name] = true
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	re := rf.includeRegexp
	for scanner.Scan() {
		line := scanner.Text()
		includes := re.FindSubmatch([]byte(line))
		if len(includes) >= 1 {
			include_fname := string("includes/" + string(includes[1]))
			text += rf.Read(include_fname)
		} else {
			text += line + string('\n')
		}
	}
	return text
}

// vim: set sw=2 ts=2 list:
