{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-q+fMrVPdFHWqGEPU4Jfn5fpNFhflWhsGDyfiqd9/4q4=";
  version = "1.0";
}
