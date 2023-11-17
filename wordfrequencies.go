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
	"strconv"
)

const (
	//ngramFileName = "dutch123grams.tsv"
	ngramFileName = "english1234grams.tsv"

	//maxTrigramsPercentage = 0.1
	//maxBigramsPercentage  = 1.0

	//maxBigramsFactor  = int(100.0 / maxBigramsPercentage)
	//maxTrigramsFactor = int(100.0 / maxTrigramsPercentage)
)

var (
	//onlyConsonants, _ = regexp.Compile(`[aeiouy]{3}`)
	//onlyVowels, _     = regexp.Compile(`[bcdfghjklmnpqrstvwxz]{3}`)
	onlyConsonants = `[aeiouy][aeiouy][aeiouy]`
	onlyVowels     = `www`
)

func randomInt(max int) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(r.Int64()) // Implicit conversion from Int64 to int
}

//func sortedValues(m map[string]int, keylen int, ordering func(i, j int) bool) []string {
//	s := make([]string, 0, len(m))
//	for k := range m {
//		if len(k) == keylen {
//			s = append(s, k)
//		}
//	}
//
//	sort.Slice(s, func(i, j int) bool {
//		return ordering(m[s[i]], m[s[j]])
//	})
//
//	return s
//}

func main() {
	const passwordLength = 12

	var excludedNgrams = [...]string{"www", "aaa"}

	//re := regexp.MustCompile(`^([a-z]{1,3})\t(\d*)\t([\d.]*) %`)
	//re := regexp.MustCompile(`^[0-9]*\. *([a-z]{1,4}) *\((\d*) *, *([\d.]*)%\)`)
	re := regexp.MustCompile(`^[0-9]*\.[[:blank:]]*([a-z]{1,4})[[:blank:]]*\((\d*)[[:blank:]]*,[[:blank:]]*([\d.]*)%\)`)

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

	nGramFreqIn := make(map[string]int)

	for scanner.Scan() {
		result := re.FindAllStringSubmatch(scanner.Text(), -1)
		if len(result) > 0 {
			nGram := result[0][1]
			count, _ := strconv.Atoi(result[0][2])
			if len(nGram) > 0 && len(nGram) < 5 {
				nGramFreqIn[nGram] = count
			}
		}
	}

	for _, k := range excludedNgrams {
		delete(nGramFreqIn, k)
	}

	allLetters := "abcdefghijklmnopqrstuvwxyz"
	for l1 := range allLetters {
		for l2 := range allLetters {
			for l3 := range allLetters {
				triGram := string(l1) + string(l2) + string(l3)
				if _, ok := nGramFreqIn[triGram]; ok {
					if match, _ := regexp.MatchString(onlyConsonants, triGram); match {
						fmt.Println("Only consonants: ", triGram)
					}
					if match, _ := regexp.MatchString(onlyVowels, triGram); match {
						fmt.Println("Only vowels: ", triGram)
					}
				}
			}
		}
	}

	//biGramValues := sortedValues(nGramFreqIn, 2, func(i, j int) bool { return i > j })
	//
	//triGramValues := sortedValues(nGramFreqIn, 3, func(i, j int) bool { return i > j })

	fmt.Printf("nGrams: %d\n", len(nGramFreqIn))

	for k, v := range nGramFreqIn {
		if len(k) == 1 {
			fmt.Println(k, ": ", v)
		}
	}
	for k, v := range nGramFreqIn {
		if len(k) == 2 {
			fmt.Println(k, ": ", v)
		}
	}
	for k, v := range nGramFreqIn {
		if len(k) == 3 {
			fmt.Println(k, ": ", v)
		}
	}
	for k, v := range nGramFreqIn {
		if len(k) == 4 {
			fmt.Println(k, ": ", v)
		}
	}

	//for i, v := range biGramValues {
	//	if nGramFreqIn[v]*maxBigramsFactor < nGramFreqIn[biGramValues[0]] {
	//		delete(nGramFreqIn, biGramValues[i])
	//	}
	//}
	//for i, v := range triGramValues {
	//	if nGramFreqIn[v]*maxTrigramsFactor < nGramFreqIn[triGramValues[0]] {
	//		delete(nGramFreqIn, triGramValues[i])
	//	}
	//}

	nGramFreqCnt := make(map[string]int)
	for k, v := range nGramFreqIn {
		nGramFreqCnt[k[0:len(k)-1]] += v
	}

	var genParams passwords.GeneratorParams
	genParams.Randomiser = randomInt
	genParams.NGramFreq, genParams.NGramCnt = nGramFreqIn, nGramFreqCnt

	passwordDictionary := make(map[string]int)
	password := ""
	maxNumOfPasswords := 100
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
	//nGrams := make(map[string]int)
	//
	//for _, s := range passwords.MonoGrams {
	//	passwords.Unravel(&nGrams, &s)
	//}
	//for _, s := range passwords.BiGrams {
	//	passwords.Unravel(&nGrams, &s)
	//}
	//for _, s := range passwords.TriGrams {
	//	passwords.Unravel(&nGrams, &s)
	//}
	//fmt.Printf("%#v", nGrams)
	//
	//nGrams = map[string]int{"aa": 1}
	//fmt.Println(nGrams)
}
