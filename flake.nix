{
  description = "QRY - Ask. Get SQL.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            golangci-lint
            goreleaser
          ];

          shellHook = ''
            echo ""
            echo "  ██████╗ ██████╗ ██╗   ██╗"
            echo " ██╔═══██╗██╔══██╗╚██╗ ██╔╝"
            echo " ██║   ██║██████╔╝ ╚████╔╝ "
            echo " ██║▄▄ ██║██╔══██╗  ╚██╔╝  "
            echo " ╚██████╔╝██║  ██║   ██║   "
            echo "  ╚══▀▀═╝ ╚═╝  ╚═╝   ╚═╝   "
            echo ""
            echo "  Ask. Get SQL."
            echo ""
            echo "  go run .           Run QRY"
            echo "  go build           Build binary"
            echo "  go test ./...      Run tests"
            echo ""
          '';

          CGO_ENABLED = "0";
        };

        packages.default = pkgs.buildGoModule rec {
          pname = "qry";
          version = "0.3.1";
          src = ./.;
          vendorHash = null;

          ldflags = [
            "-s" "-w"
            "-X github.com/amansingh-afk/qry/internal/ui.version=${version}"
          ];

          meta = with pkgs.lib; {
            description = "Ask. Get SQL.";
            homepage = "https://qry.dev";
            license = licenses.mit;
          };
        };
      }
    );
}
