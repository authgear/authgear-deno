import { walk } from "https://deno.land/std@0.161.0/fs/walk.ts";

for await (const entry of walk("/")) {
  console.log(entry);
}
