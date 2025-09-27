package faker_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

type FirstStruct struct {
	OneName string
	Second  []SecondStruct
}
type SecondStruct struct {
	SecondName string
	Third      map[string]ThirdStruct
}
type ThirdStruct struct {
	ThirdName string
	Forth     *ForthStruct
}
type ForthStruct struct {
	ForthName string
	First     *FirstStruct
}

// Single fake function can be used for retrieving particular values.
func Example_singleFakeDatax() {
	var fs FirstStruct
	_ = faker.FakeData(&fs, options.WithRandomMapAndSliceMinSize(2), options.WithRandomMapAndSliceMaxSize(3),
		options.WithMaxFieldDepthOption(3))

	jsonData, _ := json.MarshalIndent(fs, "", "  ")
	fmt.Printf("\n%s\n", jsonData)
	/*
	   Result:
	   {
	     "OneName": "YXjukTxAUhrZUYjwUoomQleqE",
	     "Second": [
	       {
	         "SecondName": "TftEhuEhbxnyVQwNJPGuUmTYy",
	         "Third": {
	           "LnZwwGyrYdiouITSkGlIRCpLj": {
	             "ThirdName": "",
	             "Forth": null
	           },
	           "vpQoyQXORYvyTHlckVJiSHsPO": {
	             "ThirdName": "",
	             "Forth": null
	           }
	         }
	       },
	       {
	         "SecondName": "CVCBTSJmaIHXSnfQPhHbAIRru",
	         "Third": {
	           "BmDkdCFiOVylYPljyptkmArIb": {
	             "ThirdName": "",
	             "Forth": null
	           },
	           "XdvTacNkbmSSvcAounuORynrp": {
	             "ThirdName": "",
	             "Forth": null
	           }
	         }
	       }
	     ]
	   }
	*/

}
