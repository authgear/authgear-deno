# authgear-deno

authgear-deno is a HTTP server that takes a JavaScript / TypeScript file and an JSON value.
The file is expected to have a default export of a function taking one argument, and return a value.
The function can be async or sync.

authgear-deno takes care of granting permission as the script runs.
Only network access to remote is granted.

## Setup

Install Deno according to [.tool-versions](./.tool-versions).

## Run

```
$ make start
```

## Examples

### Evaluate a pure function

```
$ curl --request POST \
  --url http://localhost:8090/ \
  --header 'Content-Type: application/json' \
  --data '{
	"script": "export default async function addOne(a) { return a + 1; }",
	"input": 42
}'
{"output":43}
```

### Evaluate a function with side-effects

```
$ curl --request POST \
  --url http://localhost:8090/ \
  --header 'Content-Type: application/json' \
  --data '{
	"script": "export default async function addOne(a) { console.log('\''hello'\''); return a + 1; }",
	"input": 42
}'
{"output":43,"stdout":"hello\n"}
```

### Evaluate a malicious function

```
$ curl --request POST \
  --url http://localhost:8090/ \
  --header 'Content-Type: application/json' \
  --data '{
	"script": "export default async function malicious() { Deno.remove('\''/'\'', { recursive: true}) }",
	"input": 42
}'
{"error":"exit status 1","stderr":"⚠️  ┌ Deno requests write access to \"/\".\r\n   ├ Requested by `Deno.remove()` API\r\n   ├ Run again with --allow-write to bypass this prompt.\r\n   └ Allow? [y/n] (y = yes, allow; n = no, deny) \u003e n\r\n\u001b[4A\u001b[0J❌ Denied write access to \"/\".\r\nerror: Uncaught (in promise) PermissionDenied: Requires write access to \"/\", run again with the --allow-write flag\r\n    at async Object.remove (deno:runtime/js/30_fs.js:167:5)\r\n"}
```
