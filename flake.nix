{
  description = "Dev environment for WhatGotDone";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.19.6 release
    go_dep.url = "github:NixOS/nixpkgs/6adf48f53d819a7b6e15672817fa1e78e5f4e84f";

  };

  outputs = { self, flake-utils, go_dep,  }@inputs :
    flake-utils.lib.eachDefaultSystem (system:
    let
      go_dep = inputs.go_dep.legacyPackages.${system};
    in
    {
      devShells.default = go_dep.mkShell.override { stdenv = go_dep.pkgsStatic.stdenv; } {
        packages = [
          go_dep.go_1_19
          go_dep.gopls
          go_dep.gotools
          go_dep.nodejs-18_x
          go_dep.sqlfluff
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
