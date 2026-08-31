package nextselection

import "testing"

func TestCanonicalTargetTokenUsesNumericIdentity(t *testing.T) {
	for _, test := range []struct {
		input  string
		prefix string
		want   string
	}{
		{"req-42", "REQ-", "REQ-042"},
		{"REQ-0042", "REQ-", "REQ-042"},
		{"Ur-11", "UR-", "UR-011"},
	} {
		got, err := canonicalToken(test.input, test.prefix)
		if err != nil || got != test.want {
			t.Fatalf("canonicalToken(%q, %q) = %q, %v; want %q", test.input, test.prefix, got, err, test.want)
		}
	}
	if _, err := canonicalToken("REG-42", "REQ-"); err == nil {
		t.Fatal("canonicalToken accepted an unrecognized prefix")
	}
}

func TestNextOptionAxesStayIndependent(t *testing.T) {
	options, err := parseNextOptions([]string{"--wave", "2", "--fan-out", "3", "--skip-impact-negligible"})
	if err != nil {
		t.Fatal(err)
	}
	if options.WaveDepth == nil || *options.WaveDepth != 2 || options.FanOutLimit == nil || *options.FanOutLimit != 3 || !options.SkipImpactNegligible {
		t.Fatalf("parsed options collapsed independent axes: %#v", options)
	}
	if _, err := parseNextOptions([]string{"REQ-001", "--wave", "1"}); err == nil {
		t.Fatal("--wave combined with a target")
	}
	bare, err := parseNextOptions([]string{"--fan-out"})
	if err != nil || bare.FanOutLimit == nil || *bare.FanOutLimit != 2 {
		t.Fatalf("bare --fan-out = %#v, %v; want default 2", bare, err)
	}
}
