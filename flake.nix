{
  description = "Dev environment for WhatGotDone";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-22.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self
    , nixpkgs
    , flake-utils
    }:

    flake-utils.lib.eachDefaultSystem (system:
    let
      overlays = [
        (self: super: rec {
          nodejs = super.nodejs-18_x;
        })
      ];
      pkgs = import nixpkgs { inherit overlays system; };
    in
    {
      devShells.default = pkgs.mkShell.override { stdenv = pkgs.pkgsStatic.stdenv; } {
        packages = with pkgs; [
          go_1_19
          gopls
          gotools
          node2nix
          nodejs-18_x
        ];

        shellHook = ''
          echo "node" "$(node --version)"
          go version
        '';
      };
    });
}
