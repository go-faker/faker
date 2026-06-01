package faker_test

import (
	"fmt"
	"regexp"
	"testing"
	"unicode/utf8"

	"github.com/go-faker/faker/v4"
)

type StructWithNestedTags struct {
	ListEmail          []string `faker:"[email]"`
	ListUrisWithLength []string `faker:"slice_len=3, [url]"`
	ListWordWithLang   []string `faker:"slice_len=4, [word lang=kor]"`
	// ListNumsOneOf      []float64 `faker:"[oneof: 3928.8383 83.837387]"`
}

func Example_withNestedTags() {
	a := &StructWithNestedTags{}
	_ = faker.FakeData(&a)
	fmt.Printf("%+v", a)
	// Result:
	/*
		{
			ListEmail:[
				INZkwAR@lNmQlrs.edu
				lZORYgB@oWLiiKp.net
				vLJLxiZ@eVBBKLC.com
				HkilwhV@hlukPgV.info
				JHywWlH@gJJfmyq.info
				NHSDUYy@OjYjDDq.edu
				QBaUGxv@iMYbxxO.edu
				gqGEPrh@WWtLPRd.com
				ElgLMOy@SmGWVsR.biz
				mhFwlZq@SqRfdNV.net
				DhVEmNX@qBfDqNA.ru
				cryscBF@NYgaWKl.org
				ZHlOYxS@XfGIUkv.ru
				xsDxPqs@blgGURo.net
				MnhlrGT@BVnHqvo.info
				NetrEHt@lgdJQII.info
				dyNpHcl@xAVVNSR.org
				RZgbeOw@DcgNwcj.biz
				YUgnrms@VxgOxZc.org
				GVpaHgb@vGJcZMi.info
				pbQddEH@IYDxjCq.com
				UlXqLjT@PlAhvNB.com
				AnYffFk@UimjjEs.ru
				jLSnZkh@JCdSWCo.top
				frfdUqj@kppfAYK.com
				RpapRMW@bPqqfuB.edu
				vJSpOIs@CvwQThL.org
				vZCPpBF@AIcXSVZ.edu
			]
			ListUrisWithLength:[
				https://zzshyph.top/LiMcQyr.html
				https://hdjxcmb.biz/BOSXMvB.jpg
				http://yllwaiu.org/
			]
			ListWordWithLang:[
				냶싹푣이솠휑넹뤭쾅쳎숩뎷쀷젆뉄쯃콅롾촸콝뉉똇쪈묒놬
				땭뚪휪뉯묽쫣둴뜢튔뻴쒿쬭쒃쐓멪펧뵼왌쵢뚰옐첋똠갬둎
				늒띻붬툮휣쒗쏦뙀셨쿫쎦뤿픓삠숊막뾰욤뾃쁂뼉륍펜툅샶
				샾족쭪낪챗엩바쇒뜇턝돁쏸봕맛뜼퐉뱌횋틖뫁뗽룋왙읠왂
			]
		}
	*/
}

func Test_withNestedTags(t *testing.T) {
	a := &StructWithNestedTags{}

	err := faker.FakeData(&a)

	if err != nil {
		t.Errorf("want %v, got %v", nil, err)
	}

	if len(a.ListEmail) < 1 {
		t.Errorf("Empty email list")
	}

	for _, e := range a.ListEmail {
		if m, _ := regexp.MatchString("[a-zA-Z0-9]+@[a-zA-Z0-9]+\\.[a-z]{2,3}", e); !m {
			t.Errorf("Invalid email generated %s", e)
			break
		}
	}

	if len(a.ListUrisWithLength) != 3 {
		t.Errorf("Expected list of length %d, got length %d", 3, len(a.ListUrisWithLength))
	}

	for _, e := range a.ListUrisWithLength {
		if m, _ := regexp.MatchString("(https?://)?.+\\.[a-zA-Z0-9]{2,4}", e); !m {
			t.Errorf("Invalid uri format generated %s", e)
			break
		}
	}

	if len(a.ListWordWithLang) != 4 {
		t.Errorf("Expected list of length %d, got length %d", 4, len(a.ListWordWithLang))
	}

	for _, e := range a.ListWordWithLang {
		if len(e) == utf8.RuneCountInString(e) {
			t.Errorf("String and rune of same length (%d)", len(e))
			break
		}
	}
}
