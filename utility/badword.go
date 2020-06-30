package utility

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var badWordFragmentList []string
var badWordList []string

func SetupBadWordListFromFolder(folderName string) {
	log.Info("SetupGameDefinitionFromFolder")
	badWordFragmentList = createBadWordList(filepath.Join(folderName, "badword/fragment"))
	log.Info(badWordFragmentList)
	badWordList = createBadWordList(filepath.Join(folderName, "badword/word"))
	log.Info(badWordList)
}

func createBadWordList(folderName string) []string{
	var result []string
	log.WithFields(log.Fields{
		"folderName": folderName,
	}).Info("createBadWordList")


	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.WithFields(log.Fields{
			"folderName": folderName,
			"error":      err,
		}).Error("Error loading badWordList files")
		return result
	}

	for _, file := range files {
		fullPath := filepath.Join(folderName, file.Name())
		log.WithFields(log.Fields{
			"filename": fullPath,
		}).Info("Loading badWordList")

		result = append(result, createListFromFile(fullPath)...)
	}
	return result
}

func createListFromFile(fileName string) []string{
	var result []string

	file, _ := os.Open(fileName)
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, error := reader.ReadString('\n')
		if error != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			result = append(result, strings.ReplaceAll(line, "\r\n", ""))
		}
	}
	return result
}

func ContainsBadWord(sentence string) bool {
	return isBad(removeDuplicatedChar(changeSpecialChars(strings.ToLower(sentence))))
}

func changeSpecialChars(sentence string) string {
	m := make(map[string]string)
	m["3"] = "e"
	m["1"] = "l"
	m["@"] = "a"
	m["$"] = "s"
	m["&"] = "e"
	m["!"] = "i"
	m["5"] = "s"
	m["0"] = "o"
	m["9"] = "g"

	for k, v := range m {
		sentence = strings.ReplaceAll(sentence, k, v)
	}

	sentence = removeNoChars(strings.ToLower(sentence))
	return sentence
}

func isBad(sentence string) bool {

	for _,word := range badWordList {
		vet :=strings.Split(sentence," ")
		for _, s := range vet {
			if s == word {
				return true
			}
		}
	}

	for _,word := range badWordFragmentList {
		if strings.Contains(strings.ReplaceAll(sentence, " ", ""), strings.ReplaceAll(word, " ", "")) {
			return true
		}
	}

	return false
}

func removeDuplicatedChar(sentence string) string {
	var result []byte
	for _, c := range sentence {
		if len(result)==0 || result[len(result)-1] != byte(c) {
			result = append(result, byte(c))
		}
	}
	return string(result)
}

func removeNoChars(sentence string) string {
	re := regexp.MustCompile("[^a-zA-Z ]")
	return string(re.ReplaceAll([]byte(sentence),[]byte("")))
}
