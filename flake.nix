{
  description = "Undux 🦆";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          undux = pkgs.buildGoModule {
            pname = "undux";
            version = "0.1.0";
            src = ./.;

            vendorHash = null;

            meta = with pkgs.lib; {
              description = "Undux - Unduck but faster";
              homepage = "https://github.com/warbacon/undux";
              license = licenses.mit;
            };
          };

          default = self.packages.${system}.undux;
        };

        devShells.default = import ./shell.nix { inherit pkgs; };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.undux}/bin/undux";
        };
      }
    );
}
