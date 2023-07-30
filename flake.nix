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
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          node2nix
          nodejs-18_x
          gopls
          gotools
          go_1_19
        ];

        shellHook = ''
          echo "node " "$(node --version)"
          go version
        '';
      };
    });
}
