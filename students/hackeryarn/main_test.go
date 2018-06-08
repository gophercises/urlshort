package main

import (
	"reflect"
	"testing"
)

type flaggerMock struct {
	stringVarCalls  int
	varNames        []string
	varUsages       []string
	varStringValues []string
}

func (f *flaggerMock) StringVar(p *string, name, value, usage string) {
	f.stringVarCalls++
	f.varNames = append(f.varNames, name)
	f.varStringValues = append(f.varStringValues, value)
	f.varUsages = append(f.varUsages, usage)
}

func TestConfigFlags(t *testing.T) {
	flagger := &flaggerMock{}

	ConfigFlags(flagger)

	if flagger.stringVarCalls != 2 {
		t.Error("it should set string vars")
	}

	assertFlags(t, flagger)
}

func assertFlags(t *testing.T, flagger *flaggerMock) {
	t.Helper()

	expectedNames := []string{YAMLFlag, JSONFlag}
	expectedUsages := []string{YAMLFlagUsage, JSONFlagUsage}
	expectedStringValues := []string{YAMLFlagValue, JSONFlagValue}

	if !reflect.DeepEqual(expectedNames, flagger.varNames) {
		t.Errorf("it should setup flag names to be %v, got %v",
			expectedNames, flagger.varNames)
	}

	if !reflect.DeepEqual(expectedUsages, flagger.varUsages) {
		t.Errorf("it should setup flag usages to be %v, got %v",
			expectedUsages, flagger.varUsages)
	}

	if !reflect.DeepEqual(expectedStringValues, flagger.varStringValues) {
		t.Errorf("it should setup string values to be %v, got %v",
			expectedStringValues, flagger.varStringValues)
	}
}
