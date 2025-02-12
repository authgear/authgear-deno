{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils = {
      url = "github:numtide/flake-utils";
    };
    deno_1_41_3.url = "github:NixOS/nixpkgs/080a4a27f206d07724b88da096e27ef63401a504";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      deno_1_41_3,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            # As of 2025-02-07, 1.23.6 is still unavailable in nixpkgs-unstable, so we need to use overlay to build 1.23.6 ourselves.
            (final: prev: {
              go = (
                prev.go.overrideAttrs {
                  version = "1.23.6";
                  src = prev.fetchurl {
                    url = "https://go.dev/dl/go1.23.6.src.tar.gz";
                    hash = "sha256-A5xbBOZSedrO7opvcecL0Fz1uAF4K293xuGeLtBREiI=";
                  };
                }
              );
            })
          ];
        };
        deno = deno_1_41_3.legacyPackages.${system}.deno;
      in
      {
        devShells.default = pkgs.mkShellNoCC {
          packages = [
            pkgs.go
            deno

            (pkgs.golangci-lint.overrideAttrs (
              prev:
              let
                version = "1.64.2";
              in
              {
                inherit version;
                src = pkgs.fetchFromGitHub {
                  owner = "golangci";
                  repo = "golangci-lint";
                  rev = "v${version}";
                  hash = "sha256-ODnNBwtfILD0Uy2AKDR/e76ZrdyaOGlCktVUcf9ujy8=";
                };
                vendorHash = "sha256-/iq7Ju7c2gS7gZn3n+y0kLtPn2Nn8HY/YdqSDYjtEkI=";
                # We do not actually override anything here,
                # but if we do not repeat this, ldflags refers to the original version.
                ldflags = [
                  "-s"
                  "-X main.version=${version}"
                  "-X main.commit=v${version}"
                  "-X main.date=19700101-00:00:00"
                ];
              }
            ))
          ];
        };
      }
    );
}
