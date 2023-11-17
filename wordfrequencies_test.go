package main

import (
	"fmt"
	"testing"
)

func TestCounts(t *testing.T) {
	nGramFreqIn, nGramFreqCnt := make(map[string]int), make(map[string]int)
	count1 := 0
	for l1 := 'a'; l1 <= 'z'; l1++ {
		count1 += nGramFreqIn[string(l1)]
		count2 := 0
		for l2 := 'a'; l2 <= 'z'; l2++ {
			biGram := string(l1) + string(l2)
			count2 += nGramFreqIn[biGram]
			count3 := 0
			for l3 := 'a'; l3 <= 'z'; l3++ {
				triGram := string(l1) + string(l2) + string(l3)
				count3 += nGramFreqIn[triGram]
			}
			if count3 != nGramFreqCnt[string(l1)+string(l2)] {
				fmt.Printf("Error: %d != nGramFreqCnt[%s]=%d\n", count3, string(l1)+string(l2), nGramFreqCnt[string(l1)+string(l2)])
			} else {
				fmt.Printf("'%s' ok: %d == %d\n", string(l1)+string(l2), count3, nGramFreqCnt[string(l1)+string(l2)])
			}
		}
		if count2 != nGramFreqCnt[string(l1)] {
			fmt.Printf("Error: %d != nGramFreqCnt[%s]=%d\n", count2, string(l1), nGramFreqCnt[string(l1)])
		} else {
			fmt.Printf("'%s' ok.\n", string(l1))
		}
	}
	if count1 != nGramFreqCnt[""] {
		fmt.Printf("Error: %d != nGramFreqCnt[%s]=%d\n", count1, "", nGramFreqCnt[""])
	}
}
