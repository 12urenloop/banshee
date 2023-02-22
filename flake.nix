{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
    in rec {
      defaultPackage = pkgs.buildGoModule {
        name = "banshee";
        src = pkgs.stdenv.mkDerivation {
          name = "gosrc";
          srcs = [ ./go.mod ./go.sum ./cmd ./internal ./public ./vendor ];
          phases = "installPhase";
          installPhase = ''
            mkdir $out
            for src in $srcs; do
              for srcFile in $src; do
                cp -r $srcFile $out/$(stripHash $srcFile)
              done
            done
          '';
        };
        CGO_ENABLED = 0;
        vendorSha256 = null;
        ldFlages = [
          "-S" "-W"
        ];
      };
      devShell = pkgs.mkShell rec {
        buildInputs = with pkgs; [
          go
          nix
          git
          gotools
          go-tools
          gotestsum
          gofumpt
          golangci-lint
        ];
      };
      nixosModules.website = { config, lib, pkgs, ...}:
        with lib;
        let
          cfg = config.urenloop.services.bashee;
        in {
          options.urenloop.services.bashee = {
            enable = mkEnableOption "enables bashee service";
            port = mkOption {
              type = types.int;
              default = 8080;
              example = 8080;
              description = "The port number for the bashee to run on";
            };
          };

          config = mkIf cfg.enable {
            users.users.bashee = {
              createHome = true;
              isSystemUser = true;
              group = "banshee";
              description = "https://github.com/12urenloop/banshee";
            };
            users.groups.banshee.members = [ "bashee" ];
            systemd.services.website = {
              enable = true;
              serviceConfig = {
                EnvironmentFile = "/var/lib/banshee/.env";
                WorkingDirectory = "/var/lib/bashee";
                User = "bashee";
                Group = "banshee";
                ExecStart = "${defaultPackage}/bin/banshee";
              };
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];
            };
          };
        };
    });
}
