package main

import (
	"bufio"
	//"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// include_regexp, _ := regexp.Compile(`^\s*include\s+"(.*?)"`)
	rf := ReadFile{
		doneFiles:     make(map[string]bool),
		includeRegexp: regexp.MustCompile(`^\s*include\s+"(.*?)"`),
	}
	rf.readFile("test")
	//for _, line := range rf.readFile("test") {
	//	fmt.Println(line)
	//}
}

type ReadFile struct {
	doneFiles     map[string]bool
	includeRegexp *regexp.Regexp
}

func (rf ReadFile) readFile(file_name string) []string {
	var lines []string
	if file_name[0] == '/' {
		return rf.readOneFile(file_name, lines)
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
			lines_ := rf.readOneFile(file_name, lines)
			lines = append(lines, lines_...)
		}
	}
	return lines
}

func stripnb(hostname string) string {
	return strings.Trim(hostname, "0123456789")
}

func (rf ReadFile) readOneFile(file_name string, lines []string) []string {
	rf.doneFiles[file_name] = true
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// re := regexp.MustCompile(`^\s*include\s+"(.*?)"`)
	re := rf.includeRegexp
	for scanner.Scan() {
		line := scanner.Text()
		includes := re.FindSubmatch([]byte(line))
		if len(includes) >= 1 {
			// fmt.Printf("! include %s\n", includes[1])
			include_fname := string("includes/" + string(includes[1]))
			lines = append(lines, rf.readFile(include_fname)...)
		} else {
			lines = append(lines, line)
		}
	}
	return lines
}
