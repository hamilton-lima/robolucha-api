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

var badWordFragmentMap map[string]bool
var badWordMap map[string]bool

func SetupBadWordListFromFolder(folderName string) {
	badWordFragmentMap = createBadWordMap(filepath.Join(folderName, "badword/fragment"))
	badWordMap = createBadWordMap(filepath.Join(folderName, "badword/word"))
	for key := range badWordMap {
		delete(badWordFragmentMap, key)
	}
}

func createBadWordMap(folderName string) map[string]bool {
	var badWordList []string
	log.WithFields(log.Fields{
		"folderName": folderName,
	}).Info("createBadWordList")

	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.WithFields(log.Fields{
			"folderName": folderName,
			"error":      err,
		}).Error("Error loading badWordList files")
		return nil
	}

	for _, file := range files {
		fullPath := filepath.Join(folderName, file.Name())
		log.WithFields(log.Fields{
			"filename": fullPath,
		}).Info("Loading badWordList")

		badWordList = append(badWordList, createListFromFile(fullPath)...)
	}
	badWordMap := make(map[string]bool)
	for _, word := range badWordList {
		badWordMap[word] = true
	}
	return badWordMap
}

func createListFromFile(fileName string) []string {
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
	return isBad(changeSpecialChars(strings.ToLower(sentence)))
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
	vet := strings.Split(sentence, " ")
	for _, s := range vet {
		if badWordFragmentMap[s] || badWordMap[s] {
			return true
		}
	}

	for word := range badWordFragmentMap {
		if strings.Contains(strings.ReplaceAll(sentence, " ", ""), strings.ReplaceAll(word, " ", "")) {
			return true
		}
	}

	return false
}

func removeDuplicatedChar(sentence string) string {
	var result []byte
	for _, c := range sentence {
		if len(result) == 0 || result[len(result)-1] != byte(c) {
			result = append(result, byte(c))
		}
	}
	return string(result)
}

func removeNoChars(sentence string) string {
	re := regexp.MustCompile("[^a-zA-Z ]")
	return string(re.ReplaceAll([]byte(sentence), []byte("")))
}
