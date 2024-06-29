package hw10programoptimization

import (
	"bytes"
	"testing"
)

// бенчмарки проводились так:
// на master ветке был добавлен этот файл и выполнена команда:
// go test -bench=BenchmarkGetDomainStat -benchmem -count 10 > old
// далее переключаемся на ветку hw10_program_optimization и выполняем:
// go test -bench=BenchmarkGetDomainStat -benchmem -count 10 > new
// Далее выполняем:
// benchstat ./old ./new

//nolint:all
var data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

func BenchmarkGetDomainStatCom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(bytes.NewBufferString(data), "com")
	}
}

func BenchmarkGetDomainStatGov(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(bytes.NewBufferString(data), "gov")
	}
}
