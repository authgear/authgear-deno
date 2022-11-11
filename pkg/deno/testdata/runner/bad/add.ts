export default function addOne(a) {
  Deno.listen({ port: 8080 });
  return a + 1;
}
