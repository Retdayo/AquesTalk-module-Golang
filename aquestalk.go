//go:build linux && cgo
// +build linux,cgo

package aquestalk

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>

typedef void* (*AquesTalk_Synthe_t)(const char*, int, int*);
typedef void  (*AquesTalk_FreeWave_t)(void*);
typedef void  (*AquesTalk_SetDevKey_t)(const char*);

static void* must_dlopen(const char* path) {
    return dlopen(path, RTLD_LAZY);
}

static AquesTalk_Synthe_t load_synthe(void* h) {
    return (AquesTalk_Synthe_t)dlsym(h, "AquesTalk_Synthe");
}

static AquesTalk_FreeWave_t load_free(void* h) {
    return (AquesTalk_FreeWave_t)dlsym(h, "AquesTalk_FreeWave");
}

static AquesTalk_SetDevKey_t load_setkey(void* h) {
    return (AquesTalk_SetDevKey_t)dlsym(h, "AquesTalk_SetDevKey");
}

static void call_setkey(AquesTalk_SetDevKey_t f, const char* key) {
    f(key);
}

static void* call_synthe(
    AquesTalk_Synthe_t f,
    const char* text,
    int speed,
    int* size
) {
    return f(text, speed, size);
}

static void call_free(AquesTalk_FreeWave_t f, void* p) {
    f(p);
}
*/
import "C"

import (
    "errors"
    "fmt"
    "unsafe"

    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/transform"
)

type AquesTalk struct {
    handle unsafe.Pointer
    synthe C.AquesTalk_Synthe_t
    free   C.AquesTalk_FreeWave_t
    setKey C.AquesTalk_SetDevKey_t
}

func LoadAquesTalk(libPath, devKey string) (*AquesTalk, error) {
    cpath := C.CString(libPath)
    defer C.free(unsafe.Pointer(cpath))

    h := C.must_dlopen(cpath)
    if h == nil {
        return nil, fmt.Errorf("dlopen failed: %s", libPath)
    }

    synthe := C.load_synthe(h)
    free := C.load_free(h)
    setKey := C.load_setkey(h)

    if synthe == nil || free == nil || setKey == nil {
        return nil, errors.New("failed to resolve AquesTalk symbols")
    }

    ckey := C.CString(devKey)
    defer C.free(unsafe.Pointer(ckey))
    C.call_setkey(setKey, ckey)

    return &AquesTalk{
        handle: h,
        synthe: synthe,
        free:   free,
        setKey: setKey,
    }, nil
}

func toCP932Bytes(s string) ([]byte, error) {
    enc := japanese.ShiftJIS.NewEncoder()
    out, _, err := transform.String(enc, s)
    if err != nil {
        return nil, err
    }
    return []byte(out), nil
}

func (a *AquesTalk) Synthe(text string, speechSpeed int) ([]byte, error) {
    sjis, err := toCP932Bytes(text)
    if err != nil {
        return nil, err
    }

    var size C.int
    ptr := C.call_synthe(
        a.synthe,
        (*C.char)(unsafe.Pointer(&sjis[0])),
        C.int(speechSpeed),
        &size,
    )

    if ptr == nil || size <= 0 {
        return nil, errors.New("AquesTalk_Synthe failed")
    }

    wav := C.GoBytes(ptr, size)
    C.call_free(a.free, ptr)
    return wav, nil
}
