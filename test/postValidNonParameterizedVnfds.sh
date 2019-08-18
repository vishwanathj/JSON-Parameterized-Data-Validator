#!/usr/bin/env bash
curl -X POST -i http://localhost:8080/vnfds --data  "@/usr/share/vnfdservice/test/yamlfiles/valid/nonParameterizedInput/validNonParameterizedVNFDInputWithOptionalPropConstraintsMissing.json"
