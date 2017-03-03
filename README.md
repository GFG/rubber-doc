# RocketLabsRubberDoc

## How to get started
 - clone repository
 - make install
 - go run main.go help

## Configuration
The configuration is used to provide to the generator the location os the templates and how the will be generate and also the output's destination.

#### Global
| Property  | Description |
|:----------|:----------|
| combined | The option is a flag (true or false) and allows you to decide if the output will be combined in one file or will be located in different files defined on the Templates configuration.
| srcDir | Defines the location of the templates.
| dstDir | Defines the output directory.
| output | Destination of the combined output. In case the combined property is true, this property should be set.
| templates | Configuration for each template. See section Templates below.

#### Templates
| Property  | Description |
|:----------|:----------|
| src | Template's location.
| dst | Templates's output destination. This property will be omitted when the combined property is set to true.

This example shows an configuration for a single template -> output:
```yaml
combined: false
srcDir: "__TEMPLATES_SOURCE_DIRECTORY__"
dstDir: "__OUTPUT_DESTINATION_DIRECTORY__"
templates:
  -
    src: "simple.tmpl"
    dst: "simple.html"
```

This example shows an configuration for a multiple templates -> output.
```yaml
combined: true
srcDir: "__TEMPLATES_SOURCE_DIRECTORY__"
dstDir: "__OUTPUT_DESTINATION_DIRECTORY__"
output: "__OUTPUT_FILENAME__"
templates:
  -
    src: "base.tmpl"
  -
    src: "title.tmpl"
  -
    src: "version.tmpl"
  -
    src: "baseUri.tmpl"
  -
    src: "protocols.tmpl"
  -
    src: "mediaTypes.tmpl"
```

## Run examples
### From Blueprint
 - go run main.go generate --spec=./examples/spec/blueprint/simple.apib --config=./examples/html/config.yaml

### From RAML
 - go run main.go generate --spec=./examples/spec/raml/simple.raml --config=./examples/html/config.yaml