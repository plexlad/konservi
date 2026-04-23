{
  description = "konservi dev and testing environment";
  inputs = { nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        help-script = ''
            echo "Something is said in the gaps between all information."
            echo " - Taryn Simon"
            echo ""
            echo "konservi tool versions"
            echo "→ Go:   $(go version)"
            echo "→ Node: $(node --version)"
            echo "→ PNPM: $(pnpm --version)"
            echo ""
            echo "commands:"
            echo "→ dev, test, backend, frontend"
            echo "run commands with 'nix run .#command'"
          '';
        db-compile = "go run -mod=mod entgo.io/ent/cmd/ent generate --target ./ent ./schema";
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            nodejs_25
            pnpm
            air
            golangci-lint
            delve
            ent-go
          ];

          shellHook = help-script;
        };

        packages.default = pkgs.buildGoModule {
          name = "konservi-server";
          src = ./server;
          vendorHash = null;
          doCheck = true;
        };

        apps = {
          help = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "help" help-script;
          };
          backend = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "backend" ''
              pwd
              cd server
              air
            '';
          };
          frontend = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "frontend" ''
              cd client
              pnpm dev
            '';
          };
          test = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "test" ''
              cd server
              rm -rf ent/
              ${db-compile}
              go test -v ./... && cd ../client && pnpm test
            '';
          };
          dev = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "dev" ''
              echo "Starting development environment..."
              cd server && air &
              cd ../client && pnpm dev
            '';
          };
        };
      }
    );
}
