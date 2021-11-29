# SPM-dep-cloner

Clones dependencies from .resolved file of Swift Package Manager.
Useful for setup of new project with dependencies in another repos.

### How to build

You need Go SDK installed. Build checked with go1.17.3. Run: 

```bash
go build -o spm-dep-cloner main.go
```

### Usage:

Clones repos to the working directory. TODO: flag which defines the output.

Get help with:

```bash
./spm-dep-cloner --help
```

```
Usage of ./spm-dep-cloner:
  -re string
    	specify regexp for urls
  -rev
    	reverses regexp
```

Example:

```bash
./spm-dep-cloner --re='github.com' --rev your-project.xcworkspace/xcshareddata/swiftpm/Package.resolved
```
