export default function unstableAPI(a) {
  Deno.createHttpClient({});
  return a;
}
