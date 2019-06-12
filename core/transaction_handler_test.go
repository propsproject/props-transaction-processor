package core

import (
	"testing"
	"github.com/propsproject/pending-props/core/state"
)

var tHandler = NewTransactionHandler()

func testFamilyName(t *testing.T)  {
	if name := tHandler.FamilyName(); name != FamilyName {
		t.Fatalf("expected transaction handler family name to be (%s) got (%s)", FamilyName, tHandler.FamilyName())
	}
}

func testFamilyVersions(t *testing.T)  {
	versions := tHandler.FamilyVersions()
	for key, value := range versions {
		if  value != FamilyVersions[key] {
			t.Fatalf("expected transaction handler versions to be (%s) got (%s)", FamilyVersions, tHandler.FamilyVersions())
		}
	}
}

func testNamespaces(t *testing.T)  {
	namespaces := tHandler.Namespaces()
	for key, value := range namespaces {
		if  value != state.NamespaceManager.Namespaces()[key] {
			t.Fatalf("expected transaction handler namespaces to be (%s) got (%s)", state.NamespaceManager.Namespaces, tHandler.Namespaces())
		}
	}
}

var testSuite = []interface{}{
	"testFamilyName", testFamilyName,
	"testFamilyVersions", testFamilyVersions,
	"testNamespaces", testNamespaces,
}

func TestTransactionHandler(t *testing.T)  {
	for i := 0; i < len(testSuite); i+=2 {
		testCase := testSuite[i+1].( func(t *testing.T))
		t.Run(testSuite[i].(string), testCase)
	}
}