#!/usr/bin/env bash
curl -X POST -i http://localhost:8080/vnfds --data  "@/usr/share/vnfdservice/test/yamlfiles/valid/parameterizedInput/validParameterizedVNFDInputWithOptionalPropConstraintsMissing.json"
curl -X POST -i http://localhost:8080/vnfds --data  "@/usr/share/vnfdservice/test/yamlfiles/valid/parameterizedInput/validParameterizedVNFDInputWithOptionalPropHAMissing.json"
curl -X POST -i http://localhost:8080/vnfds --data  "@/usr/share/vnfdservice/test/yamlfiles/valid/parameterizedInput/validParameterizedVNFDInputWithOptionalProps.json"
curl -X POST -i http://localhost:8080/vnfds --data  "@/usr/share/vnfdservice/test/yamlfiles/valid/parameterizedInput/validParameterizedVNFDInputWithOptionalPropScaleMissing.json"
curl -X POST -i http://localhost:8080/vnfds --data  "@/usr/share/vnfdservice/test/yamlfiles/valid/parameterizedInput/validParameterizedVNFDInputWithRequiredProps.json"