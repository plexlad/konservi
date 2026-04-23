{
  description = "konservi dev and testing environment";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        dbName = "konservi_test";
        dbUser = "postgres";
        dbPass = "postgres";
        dbPort = "5432";
        dbUrl = "postgres://${dbUser}:${dbPass}@localhost:${dbPort}/${dbName}?sslmode=disable&host=/tmp";

        # Starts a temporary postgres instance in /tmp
        start-db = pkgs.writeShellScriptBin "start-db" ''
          export PGDATA=/tmp/konservi-pgdata
          export PGPORT=${dbPort}

          if [ ! -d "$PGDATA" ]; then
            echo "Initializing postgres database..."
            ${pkgs.postgresql}/bin/initdb -D "$PGDATA" --auth=trust --no-locale --encoding=UTF8
            sed -i "s|#unix_socket_directories = '/run/postgresql'|unix_socket_directories = '/tmp'|" "$PGDATA/postgresql.conf"
          fi

          if ! ${pkgs.postgresql}/bin/pg_ctl -D "$PGDATA" status > /dev/null 2>&1; then
            echo "Starting postgres..."
            ${pkgs.postgresql}/bin/pg_ctl -D "$PGDATA" -l /tmp/konservi-pg.log start
            sleep 2
            ${pkgs.postgresql}/bin/createdb -h /tmp -p ${dbPort} ${dbName} 2>/dev/null || true
            echo "Postgres started on port ${dbPort}"
          else
            echo "Postgres already running"
          fi
        '';

        stop-db = pkgs.writeShellScriptBin "stop-db" ''
          export PGDATA=/tmp/konservi-pgdata
          ${pkgs.postgresql}/bin/pg_ctl -D "$PGDATA" stop
          echo "Postgres stopped"
        '';

        reset-db = pkgs.writeShellScriptBin "reset-db" ''
          export PGDATA=/tmp/konservi-pgdata
          ${pkgs.postgresql}/bin/pg_ctl -D "$PGDATA" stop 2>/dev/null || true
          rm -rf "$PGDATA"
          echo "Database reset. Run 'nix run .#db-start' to reinitialize."
        '';

        generate = pkgs.writeShellScriptBin "generate" ''
          cd "$(git rev-parse --show-toplevel)/server"
          echo "Running ent codegen..."
          go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema
          echo "Done."
        '';

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
          echo "→ help, dev, test, backend, frontend"
          echo "→ db-start, db-stop, db-reset"
          echo "→ generate"
          echo ""
          echo "run commands with 'nix run .#command'"
        '';
      in
      {
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

          # Database management
          db-start = flake-utils.lib.mkApp { drv = start-db; };
          db-stop  = flake-utils.lib.mkApp { drv = stop-db; };
          db-reset = flake-utils.lib.mkApp { drv = reset-db; };

          # Ent codegen
          generate = flake-utils.lib.mkApp { drv = generate; };

          # Run all tests (starts db automatically)
          test = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "test" ''
              export PGDATA=/tmp/konservi-pgdata
              export PGPORT=${dbPort}
              export KONSERVI_DB_URL="${dbUrl}"

              # Start db if not running
              ${start-db}/bin/start-db

              echo "Running backend tests..."
              cd server && go test -v ./ent/schema/...
              TEST_EXIT=$?

              echo ""
              echo "Running frontend tests..."
              cd ../client && pnpm test
              FRONTEND_EXIT=$?

              exit $(( TEST_EXIT || FRONTEND_EXIT ))
            '';
          };

          # Run only backend/schema tests
          test-backend = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "test-backend" ''
              export KONSERVI_DB_URL="${dbUrl}"
              ${start-db}/bin/start-db
              cd "$(git rev-parse --show-toplevel)/server"
              go test -v ./ent/schema/...
            '';
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

          dev = flake-utils.lib.mkApp {
            drv = pkgs.writeShellScriptBin "dev" ''
              echo "Starting development environment..."
              ${start-db}/bin/start-db
              cd server && air &
              cd ../client && pnpm dev
            '';
          };
        };
      }
    );
}
