package util

import "testing"

func TestEscapeMarkdown(t *testing.T) {
	cases := []struct {
		test     string
		expected string
	}{
		{
			"\\,`,*,{,},[,],(,),#,+,-,.,!,_,>",
			"\\\\,\\`,\\*,{,},\\[,],(,),#,+,-,.,!,\\_,>",
		},
	}
	for _, c := range cases {
		actual := EscapeMarkdown(c.test)
		if actual != c.expected {
			t.Fatalf("%s != %s", actual, c.expected)
		}
	}
}
