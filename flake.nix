{
  description = "Dev environment for WhatGotDone";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-22.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };

      pkgs_for_sqlfluff = import (builtins.fetchTarball {
        url = "https://github.com/NixOS/nixpkgs/archive/6adf48f53d819a7b6e15672817fa1e78e5f4e84f.tar.gz";
        sha256 = "0p7m72ipxyya5nn2p8q6h8njk0qk0jhmf6sbfdiv4sh05mbndj4q";
      }) {};
    in
    {
      devShells.default = pkgs.mkShell.override { stdenv = pkgs.pkgsStatic.stdenv; } {
        packages = with pkgs; [
          go_1_19
          gopls
          gotools
          node2nix
          nodejs
          pkgs_for_sqlfluff.sqlfluff
        ];

        shellHook = ''
          echo "node" "$(node --version)"
          go version
        '';
      };
    });
}
