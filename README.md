# retup

A tool that allows creating a distribution directory more flexible way for the mono repos

## Usage

Just run `retup` in the target directory contains `retup.yaml` file. 

More info : 

```shell script
Usage of retup:
      --config string   config file (default "./retup.yaml")
      --output string   output directory (default "dist")
      --version         Prints version info
```

## Configuration File

```yaml
artifacts:
  - name: public    // name of the artifact
    context: .      // context path of the artifact(by default it's ".").
    dependencies:   // dependencies of the artifact
      paths:        // paths should be set to the file dependencies for this artifact.
        - 'foo/build'
        - 'foo/bar.json'
      ignore:       // ignore specifies the paths that should be ignored. 
                    // If a file exists in both `paths` and in `ignore`, it will be ignored
                    // will only work in conjunction with `paths`
        - 'foo/build/**/*.html'
        - 'foo/build/ignore-this-dir'
```

