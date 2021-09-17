package main

import (
	"testing"
)

func TestRule3(t *testing.T) {
	ruleTest(t, difficult1, 3, 51, 20)
}

// looping rule3
func TestRule3L(t *testing.T) {
	PrepPmat(difficult1)

	totCnt := RuleLoop(rule3, RuleTable[3], Zero)
	if totCnt != 51 {
		t.Fatalf("Expected to find 51 but got %d counts.\n", totCnt)
	}
}

func TestRule_3a(t *testing.T) {

	input := "7..15..6991.37645.5.694.371..5.1.69.6.95.41.7.716.95...57.319.6.9386..1516..95..."

	ruleTest(t, input, 3, 31, 0)
}

func TestRule_3c(t *testing.T) {

	input := "..78265.16.1395.47..5147.6.3..2.1...172.8.356...6.3..4....687..82.71.6.57..5324.."

	ruleTest(t, input, 3, 36, 0)
}

func TestRule_3d(t *testing.T) {

	input := "4378265916813952472951478633..2.1978172.8.3569.8673124....687.282.71.6.57..532489"

	ruleTest(t, input, 3, 16, 0)
}

func TestRule_3e(t *testing.T) {

	input := "43782659168139524729514786336.251978172.893569586731245.396871282971.635716532489"

	ruleTest(t, input, 3, 4, 0)
}

func TestRule_3f(t *testing.T) {
	PrepPmat(difficult3)

	totCnt := RuleLoop(rule3, RuleTable[3], Zero)
	if totCnt != 29 {
		t.Fatalf("Expected to find 29 but got %d counts.\n", totCnt)
	}
}

func TestRule_3g(t *testing.T) {
	PrepPmat(difficult4)

	totCnt := RuleLoop(rule3, RuleTable[3], Zero)
	if totCnt != 51 {
		t.Fatalf("Expected to find 51 but got %d counts.\n", totCnt)
	}
}
