# Translates JsonSchema to Katydid Relapse

No code here yet, only tests

## Known Issues

There are quite a few known issues:
  - the uniqueItems keyword is not supported (this does not fit into katydid's theoretical model)
  - the patternProperties keyword is not supported (currently katydid only supports OR, NOT and ANY operators for property names and not any regular expression)
  - relapse cannot distinguish between "type":"object" and "type":"array".
