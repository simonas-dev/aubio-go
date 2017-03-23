package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>

char_t* to_char_t_ptr(char* c) {
 return (char_t*)c;
}

*/
import "C"

func toCharTPtr(s string) *C.char_t {
	return C.to_char_t_ptr(C.CString(s))
}
