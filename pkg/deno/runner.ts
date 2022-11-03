const filename = Deno.args[0];
const input = JSON.parse(await Deno.readTextFile(Deno.args[1]));
const m = await import(filename);
if (typeof m.default !== "function") {
  Deno.exit(1);
}
const output = await Promise.resolve(m.default(input));
await Deno.writeTextFile(Deno.args[2], JSON.stringify(output) + "\n");
