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

const ngramFileName = "C:\\Users\\munkh\\Downloads\\sonar_ngrams1.0\\sonar_ngrams1.0\\sonar123grams.tsv"
const (
	maxTrigramsPercentage = 100
	maxBigramsPercentage  = 1000
)

func randomInt(max int) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(r.Int64()) // Implicit conversion from Int64 to int
}

func sortedValues(m map[string]int) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}

	sort.SliceStable(ks, func(i, j int) bool {
		return m[ks[i]] > m[ks[j]]
	})

	return ks
}

func main() {
	const passwordLength = 12

	re := regexp.MustCompile(`^([a-z]{1,3})\t(\d*)\t([\d.]*) %`)
	/*	fmt.Println(result[0][1], result[0][2])
		fmt.Println(re.FindAllStringSubmatch("12	1	0.0 %", -1))*/

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

	monoGrams, biGramsIn, triGramsIn, nGramFreqIn := make(map[string]int), make(map[string]int), make(map[string]int), make(map[string]int)
	monoGramCount, biGramCountIn, triGramCountIn := make(map[string]int), make(map[string]int), make(map[string]int)

	for scanner.Scan() {
		result := re.FindAllStringSubmatch(scanner.Text(), -1)
		if len(result) > 0 {
			nGram := result[0][1]
			count, _ := strconv.Atoi(result[0][2])
			switch len(nGram) {
			case 1:
				monoGrams[nGram] = count
				monoGramCount[""] += count
			case 2:
				biGramsIn[nGram] = count
				biGramCountIn[nGram[0:1]] += count
			case 3:
				triGramsIn[nGram] = count
				triGramCountIn[nGram[0:2]] += count
			default:
				panic(fmt.Sprintf("Huh? The count for nGram %s is missing!\n", result[1][0]))
			}
			if len(nGram) > 0 && len(nGram) < 4 {
				nGramFreqIn[nGram] = count
			}
		}
	}

	biGramValues := sortedValues(biGramsIn)
	for k, v := range biGramValues {
		fmt.Println(v, k)
	}
	triGramValues := sortedValues(triGramsIn)

	maxNumOfBiGrams := 0
	for i := range biGramValues {
		if biGramsIn[biGramValues[i]]*maxTrigramsPercentage < biGramsIn[biGramValues[0]] {
			maxNumOfBiGrams = i
			break
		}
	}
	maxNumOfTriGrams := 0
	for i := range triGramValues {
		if triGramsIn[triGramValues[i]]*maxTrigramsPercentage < triGramsIn[triGramValues[0]] {
			maxNumOfTriGrams = i
			break
		}
	}
	fmt.Printf("max bigrams: %d (%d), max trigrams: %d (%d)\n", maxNumOfBiGrams, len(biGramsIn), maxNumOfTriGrams, len(triGramsIn))

	var genParams passwords.GeneratorParams
	genParams.Randomiser = randomInt
	genParams.NGramFreq, genParams.NGramCnt = monoGrams, monoGramCount
	//fmt.Println(genParams.NGramCnt, monoGramCount)

	for i, k := range biGramValues {
		if i > maxNumOfBiGrams {
			break
		}
		genParams.NGramFreq[k] = biGramsIn[k]
		genParams.NGramCnt[k[0:1]] += biGramsIn[k]
	}

	for i, k := range triGramValues {
		if i > maxNumOfTriGrams {
			break
		}
		genParams.NGramFreq[k] = triGramsIn[k]
		genParams.NGramFreq[k[0:2]] += triGramsIn[k]
		//	//fmt.Printf("%s=%d, ", k, triGramsIn[k])
	}

	for k, v := range genParams.NGramCnt {
		switch len(k) {
		case 0:
			if v != monoGramCount[k[0:len(k)]] {
				//fmt.Printf("%s: %d != %d!\n", k, v, monoGramCount[k[0:len(k)]])
			}
		case 1:
			if v != biGramCountIn[k[0:len(k)]] {
				//fmt.Printf("%s: %d != %d!\n", k, v, biGramCountIn[k[0:len(k)]])
			} else {
				//fmt.Printf("%s: %d == %d!\n", k, v, biGramCountIn[k[0:len(k)]])
			}
		}
	}
	//for i := 0; i < 20; i++ {
	//	password := genParams.Generate(12)
	//	fmt.Printf("password: %s\n", password)
	//}

	nGrams := make(map[string]int)

	for _, s := range passwords.MonoGrams {
		passwords.Unravel(&nGrams, &s)
	}
	for _, s := range passwords.BiGrams {
		passwords.Unravel(&nGrams, &s)
	}
	for _, s := range passwords.TriGrams {
		passwords.Unravel(&nGrams, &s)
	}
	fmt.Printf("%#v", nGrams)

	nGrams = map[string]int{"aa": 1}
	fmt.Println(nGrams)
}
