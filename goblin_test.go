package goblin

import (
	"os"
	"testing"
	"time"
)

func TestAddNumbersSucceed(t *testing.T) {
	fakeTest := testing.T{}
	g := Goblin(&fakeTest)

	g.Describe("Numbers", func() {
		g.It("Should add numbers", func() {
			sum := 1 + 1
			g.Assert(sum).Equal(2)
		})
	})

	if fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestAddNumbersFails(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)

	g.Describe("Numbers", func() {
		g.It("Should add numbers", func() {
			sum := 1 + 1
			g.Assert(sum).Equal(4)
		})
	})

	if !fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestMultipleIts(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)

	count := 0
	g.Describe("Numbers", func() {
		g.It("Should add numbers", func() {
			count++
			sum := 1 + 1
			g.Assert(sum).Equal(2)
		})

		g.It("Should add numbers", func() {
			count++
			sum := 1 + 1
			g.Assert(sum).Equal(4)
		})
	})

	if count != 2 {
		t.Fatal("Failed")
	}
}

func TestMultipleDescribes(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)

	count := 0
	g.Describe("Numbers", func() {

		g.Describe("Addition", func() {
			g.It("Should add numbers", func() {
				count++
				sum := 1 + 1
				g.Assert(sum).Equal(2)
			})
		})

		g.Describe("Subtraction", func() {
			g.It("Should subtract numbers", func() {
				count++
				sub := 5 - 5
				g.Assert(sub).Equal(1)
			})
		})
	})

	if count != 2 {
		t.Fatal("Failed")
	}
}

func TestPending(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)

	g.Describe("Numbers", func() {

		g.It("Should add numbers")

		g.Describe("Subtraction", func() {
			g.It("Should subtract numbers")
		})

	})

	if fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestExcluded(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)

	count := 0
	g.Describe("Numbers", func() {

		g.Xit("Should add numbers", func() {
			count++
			sum := 1 + 1
			g.Assert(sum).Equal(2)
		})

		g.Describe("Subtraction", func() {
			g.Xit("Should subtract numbers", func() {
				count++
				sub := 5 - 5
				g.Assert(sub).Equal(1)
			})
		})

	})

	if count != 0 {
		t.Fatal("Failed")
	}

	if fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestJustBeforeEach(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)
	const (
		before = iota
		beforeEach
		nBeforeEach
		justBeforeEach
		nJustBeforeEach
		it
		nIt
	)

	var (
		res [9]int
		i   int
	)

	g.Describe("Outer", func() {
		g.Before(func() {
			res[i] = before
			i++
		})

		g.BeforeEach(func() {
			res[i] = beforeEach
			i++
		})

		g.JustBeforeEach(func() {
			res[i] = justBeforeEach
			i++
		})

		g.It("should run all before handles by now", func() {
			res[i] = it
			i++
		})

		g.Describe("Nested", func() {
			g.BeforeEach(func() {
				res[i] = nBeforeEach
				i++
			})

			g.JustBeforeEach(func() {
				res[i] = nJustBeforeEach
				i++
			})

			g.It("should run all before handles by now", func() {
				res[i] = nIt
				i++
			})
		})
	})

	expected := [...]int{
		before,
		beforeEach,
		justBeforeEach,
		it,
		beforeEach,
		nBeforeEach,
		justBeforeEach,
		nJustBeforeEach,
		nIt,
	}

	if res != expected {
		t.Fatalf("expected %v to equal %v", res, expected)
	}
}

func TestNotRunBeforesOrAfters(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)
	var count int

	g.Describe("Numbers", func() {
		g.Before(func() {
			count++
		})

		g.BeforeEach(func() {
			count++
		})

		g.JustBeforeEach(func() {
			count++
		})

		g.After(func() {
			count++
		})
		g.AfterEach(func() {
			count++
		})

		g.Describe("Letters", func() {
			g.Before(func() {
				count++
			})

			g.BeforeEach(func() {
				count++
			})

			g.JustBeforeEach(func() {
				count++
			})

			g.After(func() {
				count++
			})
			g.AfterEach(func() {
				count++
			})
		})
	})

	if count != 0 {
		t.Fatal("Failed")
	}
}

func TestFailOnError(t *testing.T) {
	fakeTest := testing.T{}

	g := Goblin(&fakeTest)

	g.Describe("Numbers", func() {
		g.It("Does something", func() {
			g.Fail("Something")
		})
	})

	g.Describe("Errors", func() {
		g.It("Should fail with structs", func() {
			var s struct{ error string }
			s.error = "Error"
			g.Fail(s)
		})
	})

	if !fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestRegex(t *testing.T) {
	fakeTest := testing.T{}
	os.Args = append(os.Args, "-goblin.run=matches")
	parseFlags()
	g := Goblin(&fakeTest)

	g.Describe("Test", func() {
		g.It("Doesn't match regex", func() {
			g.Fail("Regex shouldn't match")
		})

		g.It("It matches regex", func() {})
		g.It("It also matches", func() {})
	})

	if fakeTest.Failed() {
		t.Fatal("Failed")
	}

	// Reset the regex so other tests can run
	runRegex = nil
}

func TestFailImmediately(t *testing.T) {
	fakeTest := testing.T{}
	g := Goblin(&fakeTest)
	reached := false
	g.Describe("Errors", func() {
		g.It("Should fail immediately for sync test", func() {
			g.Assert(false).IsTrue()
			reached = true
			g.Assert("foo").Equal("bar")
		})
		g.It("Should fail immediately for async test", func(done Done) {
			go func() {
				g.Assert(false).IsTrue()
				reached = true
				g.Assert("foo").Equal("bar")
				done()
			}()
		})
	})

	if reached {
		t.Fatal("Failed")
	}
}

func TestAsync(t *testing.T) {
	fakeTest := testing.T{}
	g := Goblin(&fakeTest)

	g.Describe("Async test", func() {
		g.It("Should fail when Fail is called immediately", func(done Done) {
			g.Fail("Failed")
		})
		g.It("Should fail when Fail is called", func(done Done) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				g.Fail("foo is not bar")
			}()
		})

		g.It("Should fail if done receives a parameter", func(done Done) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				done("Error")
			}()
		})

		g.It("Should pass when done is called", func(done Done) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				done()
			}()
		})

		g.It("Should fail if done has been called multiple times", func(done Done) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				done()
				done()
			}()
		})
	})

	if !fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestTimeout(t *testing.T) {
	fakeTest := testing.T{}
	os.Args = append(os.Args, "-goblin.timeout=10ms", "-goblin.run=")
	parseFlags()
	g := Goblin(&fakeTest)

	g.Describe("Test", func() {
		g.It("Should fail if test exceeds the specified timeout with sync test", func() {
			time.Sleep(100 * time.Millisecond)
		})

		g.It("Should fail if test exceeds the specified timeout with async test", func(done Done) {
			time.Sleep(100 * time.Millisecond)
			done()
		})
	})

	if !fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestItTimeout(t *testing.T) {
	fakeTest := testing.T{}
	os.Args = append(os.Args, "-goblin.timeout=10ms")
	parseFlags()
	g := Goblin(&fakeTest)

	g.Describe("Test", func() {
		g.It("Should override default timeout", func() {
			g.Timeout(20 * time.Millisecond)
			time.Sleep(15 * time.Millisecond)
		})

		g.It("Should revert for different it", func() {
			g.Assert(g.timeout).Equal(10 * time.Millisecond)
		})

	})
	if fakeTest.Failed() {
		t.Fatal("Failed")
	}
}

func TestSkipIt(t *testing.T) {
	g := Goblin(t)

	g.Describe("Describe with it, xit, skip.it", func() {

		g.It("This will run", func() {
			g.Assert(4).Equal(4)
		})

		g.Xit("This will not run", func() {
			g.Assert(4).Equal(2)
		})

		g.Skip.It("This will also not run", func() {
			g.Assert(2).Equal(2)
		})
	})
}
