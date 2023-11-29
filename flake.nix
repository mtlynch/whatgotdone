{
  description = "Dev environment for WhatGotDone";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.19.4 release
    go_dep.url = "github:NixOS/nixpkgs/2d38b664b4400335086a713a0036aafaa002c003";

    # 20.6.1 release
    nodejs_dep.url = "github:NixOS/nixpkgs/78058d810644f5ed276804ce7ea9e82d92bee293";

    # 0.9.0 release
    shellcheck_dep.url = "github:NixOS/nixpkgs/8b5ab8341e33322e5b66fb46ce23d724050f6606";

    # 1.2.1 release
    sqlfluff_dep.url = "github:NixOS/nixpkgs/7cf5ccf1cdb2ba5f08f0ac29fc3d04b0b59a07e4";

    # 0.3.7 release
    litestream_dep.url = "github:NixOS/nixpkgs/02177737c5d977444df41e0f5d6124c48c64bba3";
  };

  outputs = { self, flake-utils, go_dep, nodejs_dep, shellcheck_dep, sqlfluff_dep, litestream_dep }@inputs :
    flake-utils.lib.eachDefaultSystem (system:
    let
      go_dep = inputs.go_dep.legacyPackages.${system};
      nodejs_dep = inputs.nodejs_dep.legacyPackages.${system};
      shellcheck_dep = inputs.shellcheck_dep.legacyPackages.${system};
      sqlfluff_dep = inputs.sqlfluff_dep.legacyPackages.${system};
      litestream_dep = inputs.litestream_dep.legacyPackages.${system};
    in
    {
      devShells.default = go_dep.mkShell.override { stdenv = go_dep.pkgsStatic.stdenv; } {
        packages = [
          go_dep.gopls
          go_dep.gotools
          go_dep.go
          nodejs_dep.nodejs_20
          shellcheck_dep.shellcheck
          sqlfluff_dep.sqlfluff
          litestream_dep.litestream
        ];

        shellHook = ''
          echo "shellcheck" "$(shellcheck --version | grep '^version:')"
          sqlfluff --version
          echo "litestream" "$(litestream version)"
          echo "node" "$(node --version)"
          echo "npm" "$(npm --version)"
          go version
        '';
      };
    });
}
