// Copyright 2021 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package sqlite // import "github.com/fumiama/sqlite3"

func setMaxOpenFiles(n int) error { return nil }
