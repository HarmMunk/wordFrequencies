package main

import (
	"bufio"
	"crypto/rand"
	"example/passwords"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const (
	ngramFileName = "dutch123grams.tsv"
	//ngramFileName = "english1234grams.txt"
)

var (
	maxNgramsPercentages = []float32{0.1, 1.0}
)

func readNgrams(ngramFileName string, re *regexp.Regexp) (nGramFreqIn map[string]int, maxNgramLen int) {
	config()

	nGramsFile, err := os.Open(ngramFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func(nGramsFile *os.File) {
		err := nGramsFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(nGramsFile)

	scanner := bufio.NewScanner(nGramsFile)
	scanner.Split(bufio.ScanLines)

	nGramFreqIn = make(map[string]int)

	maxNgramLen = 0

	for scanner.Scan() {
		result := re.FindAllStringSubmatch(scanner.Text(), -1)
		if len(result) > 0 {
			nGram := result[0][1]
			count, _ := strconv.Atoi(result[0][2])
			maxNgramLen = max(maxNgramLen, len(nGram))
			nGramFreqIn[nGram] = count
		}
	}
	return nGramFreqIn, maxNgramLen
}

func randomInt(max int) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(r.Int64()) // Implicit conversion from Int64 to int
}

func sortedValues(m map[string]int, keylen int, ordering func(i, j int) bool) []string {
	s := make([]string, 0, len(m))
	for k := range m {
		if len(k) == keylen {
			s = append(s, k)
		}
	}

	sort.Slice(s, func(i, j int) bool {
		return ordering(m[s[i]], m[s[j]])
	})

	return s
}

func main() {
	const passwordLength = 12

	var excludedNgrams = []string{"www", "aaa"}

	re := regexp.MustCompile(`^([a-z]{1,3})\t(\d*)\t([\d.]*) %`)
	//re := regexp.MustCompile(`^[0-9]*\. *([a-z]{1,4}) *\((\d*) *, *([\d.]*)%\)`)
	//re := regexp.MustCompile(`^[0-9]*\.[[:blank:]]*([a-z]{1,4})[[:blank:]]*\((\d*)[[:blank:]]*,[[:blank:]]*([\d.]*)%\)`)

	nGramFreqIn, maxNgramLen := readNgrams(ngramFileName, re)

	for _, k := range excludedNgrams {
		delete(nGramFreqIn, k)
	}

	ngramValues := make(map[int][]string)

	for l := 0; l < maxNgramLen; l++ {
		ngramValues[l] = sortedValues(nGramFreqIn, l+1, func(i, j int) bool { return i > j })
	}

	for i := 0; i < min(len(ngramValues), len(maxNgramsPercentages)); i++ {
		for j, v := range ngramValues[i] {
			if nGramFreqIn[v]*int(100/maxNgramsPercentages[i]) < nGramFreqIn[ngramValues[i][0]] {
				delete(nGramFreqIn, ngramValues[i][j])
			}
		}
	}

	nGramFreqCnt := make(map[string]int)
	for k, v := range nGramFreqIn {
		nGramFreqCnt[k[0:len(k)-1]] += v
	}

	var genParams passwords.GeneratorParams
	genParams.Randomiser = randomInt
	genParams.NGramFreq, genParams.NGramCnt = nGramFreqIn, nGramFreqCnt
	genParams.MaxNGramLen = maxNgramLen

	passwordDictionary := make(map[string]int)
	password := ""
	maxNumOfPasswords := 1000000
	for i := 0; i < maxNumOfPasswords; i++ {
		doPrint := i%(maxNumOfPasswords/100) == 0
		if doPrint {
			fmt.Printf("%d", i)
		}
		password = genParams.Generate(passwordLength)
		if doPrint {
			fmt.Printf(": ")
		}
		if passwordDictionary[password] != 0 {
			//fmt.Println("Duplicate ", password, " after ", i, " passwords.")
		} else if doPrint {
			fmt.Println(password)
		}
		if password == "" {
			fmt.Println("Password is emtpy! How can that be?")
		}
		passwordDictionary[password] += 1
	}

	repeats := 0
	for k, v := range passwordDictionary {
		if v > 1 {
			fmt.Println(k, " repeats ", v, " times.")
			repeats++
		}
	}

	fmt.Println("In total, ", repeats, " passwords repeat.")
}
