package main

import (
	"bufio"
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func HashFile(t, filename string) (Sha1, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Sha1{}, err
	}
	h := sha1.New()
	fmt.Fprintf(h, "%s %d\000%s", t, len(data), data)
	var val Sha1
	for idx, v := range h.Sum(nil) {
		val[idx] = v
	}
	return val, nil
}
func HashObject(c *Client, args []string) {
	var t string
	var write, stdin, stdinpaths bool
	flag.StringVar(&t, "t", "blob", "-t object type")
	flag.BoolVar(&write, "w", false, "-w")
	flag.BoolVar(&stdin, "stdin", false, "--stdin to read an object from stdin")
	flag.BoolVar(&stdinpaths, "stdin-paths", false, "--stdin-paths to read a list of files from stdin")

	fakeargs := []string{"git-hash-object"}
	os.Args = append(fakeargs, args...)
	flag.Parse()

	if stdin && stdinpaths {
		fmt.Fprintln(os.Stderr, "Can not use both --stdin and --stdin-paths\n")
		return
	}

	var data []byte
	var err error
	if stdin {
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		h := sha1.New()
		fmt.Fprintf(h, "%s %d\000%s", t, len(data), data)
		fmt.Printf("%x\n", h.Sum(nil))
		if write {
			wsha, err := c.WriteObject(t, data)
			if err != nil {
				panic(err)
			}
		}
		return
	} else if stdinpaths {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		buffReader := bufio.NewReader(os.Stdin)
		for val, err := buffReader.ReadString('\n'); err == nil; val, err = buffReader.ReadString('\n') {
			h := sha1.New()
			data, ferr := ioutil.ReadFile(val[:len(val)-1])
			if ferr != nil {
				fmt.Fprintf(os.Stderr, "%v\n", ferr)
				return
			}
			fmt.Fprintf(h, "%s %d\000%s", t, len(data), data)
			fmt.Printf("%x\n", h.Sum(nil))
			if write {
				wsha, err := c.WriteObject(t, data)
				if err != nil {
					panic(err)
				}
			}

		}
		return
	} else {
		files := flag.Args()
		for _, file := range files {
			h := sha1.New()
			data, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				return
			}
			fmt.Fprintf(h, "%s %d\000%s", t, len(data), data)
			fmt.Printf("%x\n", h.Sum(nil))
			if write {
				wsha, err := c.WriteObject(t, data)
				if err != nil {
					panic(err)
				}
			}

		}
	}
}
