{
  description = "Dev environment for WhatGotDone";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-22.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils}:

    flake-utils.lib.eachDefaultSystem (system:
    let
      stdenv = pkgs.pkgsMusl.stdenv;
      overlays = [
        (self: super: rec {
          nodejs = super.nodejs-18_x;
        })
      ];
      pkgs = (import nixpkgs { inherit overlays system; }).pkgsMusl;
    in
    {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [ node2nix nodejs go_1_19 ];

        shellHook = ''
          echo "node `${pkgs.nodejs}/bin/node --version`"
        '';
      };
    });
}
