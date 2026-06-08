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
          version = "1.9.0";
          src = ./.;

          vendorHash = "sha256-fkbOmBaFsprjUaMqEhjNiH0+wxGjgL1+i5LBBL8Xzlg=";

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