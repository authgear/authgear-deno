export default async function identity(a) {
  await Deno.readTextFile("/tmp/foobar");
  return a;
}
