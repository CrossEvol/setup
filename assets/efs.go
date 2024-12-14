package assets

import _ "embed"

//go:embed gitignore_pairs.json
var GitignorePairs []byte

//go:embed doc_pairs.json
var DocPairs []byte
