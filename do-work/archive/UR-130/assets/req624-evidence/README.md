# Replaying the changelog history proof

The proof verifies the historical trim at commit `75ecfde4505feea97aeb196de98e44f8eab1bff3`. Run it from the root of a separate, clean checkout of that commit:

```sh
python3 do-work/archive/UR-130/assets/req624-evidence/verify-req624-history.py
```

It compares that checkout and its `HEAD` archive with a fixed pre-trim baseline. Later releases are outside this proof's input and correctly fail its exact release-set assertions. The script and recorded evidence remain unchanged.
