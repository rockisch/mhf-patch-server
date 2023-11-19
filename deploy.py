from zipfile import ZipFile
from pathlib import Path


deploy_dir = Path("deploy/")
if not deploy_dir.exists():
    deploy_dir.mkdir()

for exe, target in (
    ("mhf-patch-server", "Linux-amd64.zip"),
    ("mhf-patch-server.exe", "Windows-amd64.zip"),
):
    if Path(exe).exists():
        with ZipFile(deploy_dir / target, "w") as zip:
            zip.mkdir("game")
            zip.write("patch_config.json")
            zip.write(exe)
