// Copyright 2020 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package libc2 // import "github.com/fumiama/sqlite3/internal/libc2"

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"modernc.org/libc"
	"modernc.org/libc/sys/types"
)

func todo(s string, args ...interface{}) string { //TODO-
	switch {
	case s == "":
		s = fmt.Sprintf(strings.Repeat("%v ", len(args)), args...)
	default:
		s = fmt.Sprintf(s, args...)
	}
	pc, fn, fl, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	var fns string
	if f != nil {
		fns = f.Name()
		if x := strings.LastIndex(fns, "."); x > 0 {
			fns = fns[x+1:]
		}
	}
	r := fmt.Sprintf("%s:%d:%s: TODOTODO %s", fn, fl, fns, s) //TODOOK
	fmt.Fprintf(os.Stdout, "%s\n", r)
	os.Stdout.Sync()
	return r
}

func trc(s string, args ...interface{}) string { //TODO-
	switch {
	case s == "":
		s = fmt.Sprintf(strings.Repeat("%v ", len(args)), args...)
	default:
		s = fmt.Sprintf(s, args...)
	}
	_, fn, fl, _ := runtime.Caller(1)
	r := fmt.Sprintf("\n%s:%d: TRC %s", fn, fl, s)
	fmt.Fprintf(os.Stdout, "%s\n", r)
	os.Stdout.Sync()
	return r
}

// int sched_yield(void);
func Xsched_yield(tls *libc.TLS) int32 {
	panic(todo(""))
}

// int pthread_create(pthread_t *thread, const pthread_attr_t *attr, void *(*start_routine) (void *), void *arg);
func Xpthread_create(tls *libc.TLS, thread, attr, start_routine, arg uintptr) int32 {
	panic(todo(""))
}

// int pthread_detach(pthread_t thread);
func Xpthread_detach(tls *libc.TLS, thread types.Pthread_t) int32 {
	panic(todo(""))
}

// int pthread_mutex_lock(pthread_mutex_t *mutex);
func Xpthread_mutex_lock(tls *libc.TLS, mutex uintptr) int32 {
	panic(todo(""))
}

// int pthread_cond_signal(pthread_cond_t *cond);
func Xpthread_cond_signal(tls *libc.TLS, cond uintptr) int32 {
	panic(todo(""))
}

// int pthread_mutex_unlock(pthread_mutex_t *mutex);
func Xpthread_mutex_unlock(tls *libc.TLS, mutex uintptr) int32 {
	panic(todo(""))
}

// int pthread_mutex_init(pthread_mutex_t *restrict mutex, const pthread_mutexattr_t *restrict attr);
func Xpthread_mutex_init(tls *libc.TLS, mutex, attr uintptr) int32 {
	panic(todo(""))
}

// int pthread_cond_init(pthread_cond_t *restrict cond, const pthread_condattr_t *restrict attr);
func Xpthread_cond_init(tls *libc.TLS, cond, attr uintptr) int32 {
	panic(todo(""))
}

// int pthread_cond_wait(pthread_cond_t *restrict cond, pthread_mutex_t *restrict mutex);
func Xpthread_cond_wait(tls *libc.TLS, cond, mutex uintptr) int32 {
	panic(todo(""))
}

// int pthread_cond_destroy(pthread_cond_t *cond);
func Xpthread_cond_destroy(tls *libc.TLS, cond uintptr) int32 {
	panic(todo(""))
}

// int pthread_mutex_destroy(pthread_mutex_t *mutex);
func Xpthread_mutex_destroy(tls *libc.TLS, mutex uintptr) int32 {
	panic(todo(""))
}

// int pthread_mutex_trylock(pthread_mutex_t *mutex);
func Xpthread_mutex_trylock(tls *libc.TLS, mutex uintptr) int32 {
	panic(todo(""))
}

// int pthread_cond_broadcast(pthread_cond_t *cond);
func Xpthread_cond_broadcast(tls *libc.TLS, cond uintptr) int32 {
	panic(todo(""))
}
