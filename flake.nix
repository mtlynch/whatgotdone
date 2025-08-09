{
  description = "Dev environment for WhatGotDone";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.23.1 release
    go-nixpkgs.url = "github:NixOS/nixpkgs/5ed627539ac84809c78b2dd6d26a5cebeb5ae269";

    # 3.44.2 release
    sqlite-nixpkgs.url = "github:NixOS/nixpkgs/5ad9903c16126a7d949101687af0aa589b1d7d3d";

    # 20.6.1 release
    nodejs-nixpkgs.url = "github:NixOS/nixpkgs/78058d810644f5ed276804ce7ea9e82d92bee293";

    # 0.9.0 release
    shellcheck-nixpkgs.url = "github:NixOS/nixpkgs/8b5ab8341e33322e5b66fb46ce23d724050f6606";

    # 1.2.1 release
    sqlfluff-nixpkgs.url = "github:NixOS/nixpkgs/7cf5ccf1cdb2ba5f08f0ac29fc3d04b0b59a07e4";

    # 0.3.13 release
    litestream-nixpkgs.url = "github:NixOS/nixpkgs/a343533bccc62400e8a9560423486a3b6c11a23b";
  };

  outputs = {
    self,
    flake-utils,
    go-nixpkgs,
    sqlite-nixpkgs,
    nodejs-nixpkgs,
    shellcheck-nixpkgs,
    sqlfluff-nixpkgs,
    litestream-nixpkgs,
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      go = go-nixpkgs.legacyPackages.${system}.go_1_23;
      sqlite = sqlite-nixpkgs.legacyPackages.${system}.sqlite;
      nodejs = nodejs-nixpkgs.legacyPackages.${system}.nodejs_20;
      shellcheck = shellcheck-nixpkgs.legacyPackages.${system}.shellcheck;
      sqlfluff = sqlfluff-nixpkgs.legacyPackages.${system}.sqlfluff;
      litestream = litestream-nixpkgs.legacyPackages.${system}.litestream;
    in {
      devShells.default =
        go-nixpkgs.legacyPackages.${system}.mkShell.override
        {
          stdenv = go-nixpkgs.legacyPackages.${system}.pkgsStatic.stdenv;
        }
        {
          packages = [
            go-nixpkgs.legacyPackages.${system}.gotools
            go-nixpkgs.legacyPackages.${system}.gopls
            go-nixpkgs.legacyPackages.${system}.go-outline
            go-nixpkgs.legacyPackages.${system}.gopkgs
            go-nixpkgs.legacyPackages.${system}.gocode-gomod
            go-nixpkgs.legacyPackages.${system}.godef
            go-nixpkgs.legacyPackages.${system}.golint
            go
            sqlite
            nodejs
            shellcheck
            go-nixpkgs.legacyPackages.${system}.screen
            sqlfluff
            litestream
          ];

          shellHook = ''
            export GOROOT="${go}/share/go"

            echo "shellcheck" "$(shellcheck --version | grep '^version:')"
            sqlfluff --version
            echo "litestream" "$(litestream version)"
            echo "node" "$(node --version)"
            echo "npm" "$(npm --version)"
            echo "sqlite" "$(sqlite3 --version | cut -d ' ' -f 1-2)"
            go version
          '';
        };

      formatter = go-nixpkgs.legacyPackages.${system}.alejandra;
    });
}
