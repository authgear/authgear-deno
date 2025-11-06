const filename = Deno.args[0];
const input = JSON.parse(await Deno.readTextFile(Deno.args[1]));
const m = await import(filename);
if (typeof m.default !== "function") {
  console.error(
    "The hook must export a default function. Check that you have `export default async function(...) { ... }` in your script.",
  );
  Deno.exit(1);
}
const output = await Promise.resolve(m.default(input));
let content = JSON.stringify(output);
if (content === undefined) {
  content = "null";
}
await Deno.writeTextFile(Deno.args[2], content + "\n");
