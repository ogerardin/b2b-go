// Implements a name-to-type mechanism that allows a reflect.Type to be retrieved from a string representation
// at runtime, roughly similar to Java's Class.forName(). This in turn allows the mashalling/unmarshalling
// of structures that are not known at compile-time, e.g. for persistence in a Mongo database or as a JSON object.
//
// Important: only types that have been registered with typeregistry.Register() will be known. If you own the type and
// don't mind to tie your code to this package, you might want to add the call to Register() in an init function:
//
//      func init() {
//  	    typeregistry.Register(reflect.TypeOf((*MyType)(nil)).Elem())
//      }

package typeregistry
