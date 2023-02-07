export default function useUnstableAPI(a) {
  Deno.createHttpClient({});
  return a;
}
