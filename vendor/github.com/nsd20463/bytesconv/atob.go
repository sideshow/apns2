// Copyright 2016 Nicolas Dade. Derived from the go 1.7.4
// standard library, which itself was:
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytesconv

// ParseBool returns the boolean value represented by the string in []byte.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns an error.
func ParseBool(str []byte) (bool, error) {
	if len(str) < 1 {
		return false, syntaxError("ParseBool", str)
	}
	switch str[0] {
	case '1':
		if len(str) == 1 {
			return true, nil
		}
	case '0':
		if len(str) == 1 {
			return false, nil
		}
	case 't', 'T':
		switch len(str) {
		case 1:
			return true, nil
		case 4:
			if (str[1] == 'r' || str[1] == 'R') &&
				(str[2] == 'u' || str[2] == 'U') &&
				(str[3] == 'e' || str[3] == 'E') {
				return true, nil
			}
		}

	case 'f', 'F':
		switch len(str) {
		case 1:
			return false, nil
		case 5:
			if (str[1] == 'a' || str[1] == 'A') &&
				(str[2] == 'l' || str[2] == 'L') &&
				(str[3] == 's' || str[3] == 'S') &&
				(str[4] == 'e' || str[4] == 'E') {
				return false, nil
			}
		}
	}
	return false, syntaxError("ParseBool", str)
}
