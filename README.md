# yamltojson

Convert YAML files to JSON files.

A single file or a directory of files can be converted:

```
yamltojson /path/to/my/file.yaml

yamltojson /path/to/my/yaml/files/
```

The output is one or more `.json` files alongside the original input `.yaml`
files. The single file variant does not require a `.yaml` extension but
the directory version will skip files without a `.yaml` extension.
