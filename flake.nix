{
  description = "Dev environment for WhatGotDone";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-22.11";
    flake-utils.url = "github:numtide/flake-utils";

    # 1.19.6 release
    go_dep.url = "github:NixOS/nixpkgs/6adf48f53d819a7b6e15672817fa1e78e5f4e84f";

    # 18.14.1 release
    nodejs_dep.url = "github:NixOS/nixpkgs/6adf48f53d819a7b6e15672817fa1e78e5f4e84f";

    # 1.2.1 release
    sqlfluff_dep.url = "github:NixOS/nixpkgs/7cf5ccf1cdb2ba5f08f0ac29fc3d04b0b59a07e4";
  };

  outputs = { self, nixpkgs, flake-utils, go_dep, nodejs_dep, sqlfluff_dep }@inputs :
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = inputs.nixpkgs.legacyPackages.${system};
      go_dep = inputs.go_dep.legacyPackages.${system};
      sqlfluff_dep = inputs.sqlfluff_dep.legacyPackages.${system};
      nodejs_dep = inputs.nodejs_dep.legacyPackages.${system};
    in
    {
      devShells.default = go_dep.mkShell.override { stdenv = go_dep.pkgsStatic.stdenv; } {
        packages = with pkgs; [
          gopls
          gotools
          go_dep.go_1_19
          nodejs_dep.nodejs-18_x
          sqlfluff_dep.sqlfluff
        ];

        shellHook = ''
          echo "node" "$(node --version)"
          echo "npm" "$(npm --version)"
          go version
          sqlfluff --version
        '';
      };
    });
}
