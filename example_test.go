// Copyright 2024 Terminal Stream Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strum_test

import (
	"fmt"

	"github.com/terminalstream/strum"
)

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

	err := strum.Unmarshal(data, &contact)
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

func ExampleUnmarshal_formatter() {
	const data = "abc"

	test := struct {
		Val int `strform:"lettersToNumbers" strum:"0"`
	}{}

	err := strum.Unmarshal(
		data,
		&test,
		strum.WithFormatter("lettersToNumbers", func(string) (string, error) {
			return "123", nil
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(test.Val)

	// Output: 123
}
