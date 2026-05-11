{
  description = "Piton programming language";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "piton";
          version = "1.8.2";
          src = ./.;

          vendorHash = "sha256-lSQ2BIQ631yqEPfuIuA4xasHy6Uz8M2E0ac5vp3CUMw=";

          subPackages = [ "cmd/piton" ];
          
          env = {
            CGO_ENABLED = 0;
          };

          ldflags = [ "-s" "-w" ];
          
          deleteVendor = false; 
          proxyVendor = true; 
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.gotools
            pkgs.nixpkgs-fmt
            pkgs.mdbook
            pkgs.tinygo
            pkgs.binaryen # wasm-opt
          ];
        };
      });
}