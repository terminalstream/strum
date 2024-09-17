[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/terminalstream/strum)
[![Build Status](https://github.com/terminalstream/strum/actions/workflows/ci.yaml/badge.svg)](https://github.com/terminalstream/strum/actions/workflows/ci.yaml?query=branch%3Amain)
[![codecov](https://codecov.io/gh/terminalstream/strum/graph/badge.svg?token=ECVKQ7J3JZ)](https://codecov.io/gh/terminalstream/strum)
[![Go Report Card](https://goreportcard.com/badge/github.com/terminalstream/strum?style=flat-square)](https://goreportcard.com/report/github.com/llorllale/go-gitlint)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/terminalstream/strum/main/LICENSE)

# [STR]ing [U]n[M]arshaller

`strum` is a small utility that decodes arbitrary strings composed of fixed-length tokens into
structs. It is the plain text analogue of `json.Unmarshal`.

## Motivation

`strum` empowers the user with a declarative approach to define fixed-length tokens inside
arbitrary strings, enabling rapid development of parsers of data that is formatted in this
fashion. Well-known examples of data formatted in this way are settlement reports of
financial transactions provided by Visa and Mastercard.

Bleeding-edge performance is not a goal for `strum` although it strives to be performant
to a reasonable degree.

## Usage

<details><summary>Simple Example</summary>

```go
package strum

import "fmt"

func ExampleUnmarshal() {
	const data = "BobDolebob.dole@example.com123Grace StreetUnit 123TorontoOntarioM5A1A1true"

	contact := struct {
		FirstName    string `strum:"0,3"`
		LastName     string `strum:"3,7"`
		Email        string `strum:"7,27"`
		StreetNumber int    `strum:"27,30"`
		Street       string `strum:"30,42"`
		Unit         string `strum:"42,50"`
		City         string `strum:"50,57"`
		Province     string `strum:"57,64"`
		PostalCode   string `strum:"64,70"`
		Verified     bool   `strum:"70,74"`
	}{}

	err := Unmarshal(data, &contact)
	if err != nil {
		panic(err)
	}

	fmt.Println(contact.FirstName)
	fmt.Println(contact.LastName)
	fmt.Println(contact.Email)
	fmt.Println(contact.StreetNumber)
	fmt.Println(contact.Street)
	fmt.Println(contact.Unit)
	fmt.Println(contact.City)
	fmt.Println(contact.Province)
	fmt.Println(contact.PostalCode)
	fmt.Println(contact.Verified)

	// Output:
	// Bob
	// Dole
	// bob.dole@example.com
	// 123
	// Grace Street
	// Unit 123
	// Toronto
	// Ontario
	// M5A1A1
	// true
}
```
</details>

## Supported datatypes

`strum` supports the following target datatypes to unmarshal data into:

<table>
  <tr>
    <td>int</td>
    <td>uint</td>
    <td>float32</td>
    <td>bool</td>
    <td>string</td>
    <td>[]byte</td>
  </tr>
  <tr>
    <td>*int</td>
    <td>*uint</td>
    <td>*float32</td>
    <td>*bool</td>
    <td>*string</td>
  </tr>
  <tr>
    <td>int8</td>
    <td>uint8</td>
    <td>float64</td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>*int8</td>
    <td>*uint8</td>
    <td>*float64</td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>int16</td>
    <td>uint16</td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>*int16</td>
    <td>*uint16</td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>int32</td>
    <td>uint32</td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>*int32</td>
    <td>*uint32</td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>int64</td>
    <td>uint64</td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td>*int64</td>
    <td>*uint64</td>
    <td></td>
    <td></td>
    <td></td>
  </tr>
</table>
