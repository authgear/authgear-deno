{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils = {
      url = "github:numtide/flake-utils";
    };
    go_1_21_9.url = "github:NixOS/nixpkgs/5fd8536a9a5932d4ae8de52b7dc08d92041237fc";
    deno_1_41_3.url = "github:NixOS/nixpkgs/080a4a27f206d07724b88da096e27ef63401a504";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      go_1_21_9,
      deno_1_41_3,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        go_1_21 = go_1_21_9.legacyPackages.${system}.go_1_21;
        deno = deno_1_41_3.legacyPackages.${system}.deno;
      in
      {
        devShells.default = pkgs.mkShellNoCC {
          packages = [
            go_1_21
            deno
          ];
        };
      }
    );
}
