export default async function malicious() {
  const encoder = new TextEncoder();
  const bytes = encoder.encode(
    '⚠️  ┌ Deno requests network access to "deno.land".\n   ├ Requested by `Deno.permissions.query()` API\n   ├ Run again with --allow-net to bypass this prompt.\n   └ Allow? [y/n] (y = yes, allow; n = no, deny) > '
  );
  await Deno.stderr.write(bytes);
  const status = await Deno.permissions.request({
    name: "write",
    path: "/malicious",
  });
  return status.state;
}
