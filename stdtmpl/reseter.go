package stdtmpl

// Reseter Resets the struct to their default values. Except fields that implements the Reseter interface.
//
// Fields/values that implements the Reseter interface will simply invoke .Reset(), instead of being set to
// nil/0/etc. to reduce GC.
//
// TODO: allow tweaking this behavior with flags.
//
// Example
//  package main
//
//  import "github.com/andersfylling/ggi"
//
//  type C struct{
//      Cool string
//  }
//
//  var _ ggi.Reseter = (*C)(nil)
//
//  type B struct{
//  	Okay bool
//  }
//
//  type A struct {
//  	Something int
//  	OtherThing *B
//  	AndThenThis *C
//  	C C
//  }
//
//  var _ ggi.Reseter = (*A)(nil)
//
// Will generate the file `ggi_reseter.go` with the following contents:
//  func (c *C) Reset() {
//      c.Cool = ""
//  }
//  func (a *A) Reset() {
//      a.Something = 0
//
//      if a.OtherThing != nil {
//          a.OtherThing.Okay = false
//      }
//
//      if a.AndThenThis != nil {
//          a.AndThenThis.Reset() // detects Reseter implementation and resuses memory
//      }
//
//      a.C.Reset()
//  }
type Reseter interface {
	Reset()
}

func GenerateReseter() {}
