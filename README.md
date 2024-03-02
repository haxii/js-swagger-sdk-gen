# js-swagger-sdk-gen
yet another javascript swagger package generate tool, using axios written in pure go.
CI ready, you can generate and publish a npm package without node env.

## How To Use

just one command to generate a pet-store-api@1.0.9 using the pet-store swagger and publish it to your registry: 

```bash
js-swagger-sdk-gen -f https://petstore.swagger.io/v2/swagger.json -n pet-store-api -v 1.0.9 -p
```

in this example

- `-f` specifies the path of swagger
- `-n` specifies the npm package name
- `-v` indicates the package version should be v1.0.9
- `-p` tells the command to publish to you default registry using the settings in `.npmrc`

you can try this generated pet store package at https://www.npmjs.com/package/pet-store-api

following example shows how to use the pet-store-api to get all the pending pets.
```javascript
const axios = require("axios");
const { default: createAxios } = require("pet-store-api");
const { setDomain, findPetsByStatus } = require("pet-store-api");

// init axios and set domain
createAxios(axios.default);
setDomain("https://petstore.swagger.io/v2");

// print all the pending pets
;(async function () {
  const pets = await findPetsByStatus({ status: "pending" });
  console.log(pets?.data);
})();
```
## Features

- swagger's OperationID as function name in the generated code, support both commonJS and ES module
- publish the generated package directly onto your registry  
- pure go implementation without any dependency, yes, you can make and publish a npm package without node env.
- small binary size ~3m friendly on CI/CD env.
- with `-ui` option to output a swagger-ui dest folder for documentation serving

for more options, try the `--help` command:

```
> js-swagger-sdk-gen --help

Usage:
  js-swagger-sdk-gen [OPTIONS]

Generate and publish a JavaScript SDK using axios with given swagger v2 specification.

Application Options:
  -f, --file=             swagger file, support both local and remote json/yaml files (default: ./swagger.json)
  -t, --target=           target dir to generate the SDK (default: ./)
  -p, --publish           publish to the registry directly, if set, the tarball will not write to the target dir
      --ui=               generate swagger ui to this dir for distribution

SDK Package Options:
  -n, --pkg-name=         sdk name, required, default from swagger's info.package_name
  -v, --pkg-version=      sdk version, required, default from swagger's version
      --pkg-description=  sdk description, default from swagger's info.description
      --pkg-author-name=  sdk author name, default from swagger's info.contact.name
      --pkg-author-email= sdk author email, default from swagger's info.contact.email
      --pkg-homepage=     sdk homepage, default from swagger's info.homepage
      --pkg-license=      sdk license, default from swagger's info.license

NPM Registry Options:
      --registry-url=     npm registry url to publish the SDK, default from .npmrc [$NPM_REGISTRY_URL]
      --registry-token=   npm registry token to publish the SDK, default from .npmrc [$NPM_REGISTRY_TOKEN]

Miscellaneous Options:
      --version           display application version
      --verbose           verbose the output

Help Options:
  -h, --help              Show this help message
```