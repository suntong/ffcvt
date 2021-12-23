////////////////////////////////////////////////////////////////////////////
// Porgram: flag_array
// Purpose: a user-defined Go flag type that allows multiple occurrence
// Authors: Tong Sun (c) 2021, All rights reserved
// License: MIT
////////////////////////////////////////////////////////////////////////////
// Refs https://github.com/suntong/lang/blob/master/lang/Go/src/sys/CommandLineFlagArray.go

package main

// mFlags extend Go flags so that it can be specified multiple times
type mFlags []string

func (f *mFlags) String() string {
	return "n.a."
}

func (f *mFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

/*

Usage Example:

var multiple mFlags

func main() {
 flag.Var(&multiple, "list1", "Some description for this param.")
 flag.Parse()
}

*/
