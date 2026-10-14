package git

import (
	"testing"
)

func FuzzParseStatus(f *testing.F) {
	// Seed corpus
	f.Add("1 A. N... 000000 100644 100644 0000000 1234567 staged.txt\x00")
	f.Add("2 R. N... 100644 100644 100644 1234567 1234567 R100 new.txt\x00old.txt\x00")
	f.Add("u UU N... 100644 100644 100644 100644 1111111 2222222 3333333 c.txt\x00")
	f.Add("? untracked.go\x00")
	f.Add("! ignored.o\x00")
	f.Add("")

	f.Fuzz(func(t *testing.T, data string) {
		// Parsers should never panic on random input
		_, _ = ParseStatusPorcelainV2(data)
	})
}

func FuzzParseLog(f *testing.F) {
	// Seed corpus
	f.Add("1234\x0012\x00\x00t1\x00Alice\x00a@b.com\x002026-01-01T00:00:00Z\x00Alice\x00a@b.com\x002026-01-01T00:00:00Z\x00Subject\x00Body\x00HEAD -> main\x01")
	f.Add("aaa\x00a\x00bbb ccc\x00t\x00Name\x00n@n.com\x00bad-date\x00Name\x00n@n.com\x00bad-date\x00Sub\x00Body\x00tag: v1\x01")
	f.Add("")

	f.Fuzz(func(t *testing.T, data string) {
		_, _ = ParseLog(data)
	})
}

func FuzzParseDiff(f *testing.F) {
	// Seed corpus
	f.Add(`diff --git a/f b/f
index 1..2 100644
--- a/f
+++ b/f
@@ -1,1 +1,2 @@
 context
+added
`)
	f.Add("diff --git a/old b/new\nrename from old\nrename to new\n")
	f.Add("")

	f.Fuzz(func(t *testing.T, data string) {
		_, _ = ParseDiff(data)
		_ = ParseDiffStat(data)
	})
}
