package main

import "testing"

func TestReplaceProfane(t *testing.T) {
  s := "This is a kerfuffle opinion I need to share with the world"
  repl_s := ReplaceProfane(s)
  if repl_s != "This is a **** opinion I need to share with the world" {
    t.Fatalf("replaced string %s, differs from the expected one", repl_s)
  }
}

func TestReplaceProfaneDouble(t *testing.T) {
  s := "I really need a kerfuffle to go to bed sooner, Fornax !"
  repl_s := ReplaceProfane(s)
  if repl_s != "I really need a **** to go to bed sooner, **** !" {
    t.Fatalf("replaced string %s, differs from the expected one", repl_s)
  }
}
