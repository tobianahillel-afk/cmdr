#!/usr/bin/env python3
from __future__ import annotations
import base64, json, lzma, os, pathlib, shutil, subprocess, sys

ROOT = pathlib.Path(__file__).resolve().parents[1]
os.chdir(ROOT)
PART_DIR = ROOT / ".cmdr-generator"

def run(*args: str) -> None:
    subprocess.run(args, check=True)

def write_files(data: dict[str, str], paths: list[str]) -> None:
    for rel in paths:
        target = ROOT / rel
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(data[rel], encoding="utf-8", newline="\n")

def top_domain(path: str) -> str:
    rel = path.removeprefix("cmdr-product-spec/")
    return rel.split("/", 1)[0] if "/" in rel else "__root__"

parts = sorted(PART_DIR.glob("payload-*.txt"))
if not parts:
    raise SystemExit("No payload parts found")
encoded = "".join(p.read_text(encoding="ascii").strip() for p in parts)
data = json.loads(lzma.decompress(base64.b64decode(encoded)).decode("utf-8"))
if not isinstance(data, dict) or len(data) != 781:
    raise SystemExit(f"Unexpected manifest size: {len(data) if isinstance(data, dict) else 'invalid'}")

endpoint_key = "cmdr-product-spec/11-endpoint-agent/README.md"
endpoint_requirements = """
## Exigences EDR natives obligatoires

L’Endpoint Agent est un EDR natif complet. Il couvre explicitement la télémétrie, la détection locale, l’inspection, la collecte, le terminal et l’exécution de commandes, la quarantaine et l’isolation, le rollback et sa vérification, le mode hors ligne, la protection anti-altération, la signature des commandes, l’audit local et le stockage sécurisé.
"""
if "## Exigences EDR natives obligatoires" not in data[endpoint_key]:
    data[endpoint_key] = data[endpoint_key].rstrip() + "\n\n" + endpoint_requirements.strip() + "\n"

root_readme = (ROOT / "README.md").read_text(encoding="utf-8")
if root_readme != "# cmdr\n":
    raise SystemExit("Root README changed before generation")

spec_root = ROOT / "cmdr-product-spec"
if spec_root.exists():
    shutil.rmtree(spec_root)

groups = [
    ("docs: align repository with canonical CMDR architecture",
     {"__root__", "00-governance", "01-product-vision", "02-brand"}),
    ("docs: complete design system experience and domain model",
     {"03-design-system", "04-experience-architecture", "05-domain-model"}),
    ("docs: add Command Investigate and Govern specifications",
     {"06-command", "07-investigate", "08-govern"}),
    ("docs: add CMDR Studio Platform Settings and Endpoint Agent",
     {"09-cmdr-studio", "10-platform-settings", "11-endpoint-agent"}),
    ("docs: complete shared journeys security and content",
     {"12-shared-capabilities", "13-user-journeys",
      "14-security-permissions-and-trust", "15-content-and-language"}),
    ("docs: add quality contracts roadmap templates and archive",
     {"16-quality-and-validation", "17-implementation-contracts",
      "18-roadmap-and-releases", "templates", "assets", "99-archive"}),
]

all_paths = sorted(data)
for message, domains in groups:
    selected = [p for p in all_paths if top_domain(p) in domains]
    write_files(data, selected)
    run("git", "add", "-A", "cmdr-product-spec")
    status = subprocess.run(
        ["git", "diff", "--cached", "--quiet"], check=False
    ).returncode
    if status != 0:
        run("git", "commit", "-m", message)

# Structural validation before cleanup.
expected = set(all_paths)
actual = {
    p.relative_to(ROOT).as_posix()
    for p in spec_root.rglob("*")
    if p.is_file()
}
if actual != expected:
    missing = sorted(expected - actual)
    extra = sorted(actual - expected)
    raise SystemExit(f"Manifest mismatch missing={missing[:10]} extra={extra[:10]}")
empty = [p for p in actual if not (ROOT / p).read_text(encoding="utf-8").strip()]
if empty:
    raise SystemExit(f"Empty files: {empty[:10]}")
if (ROOT / "README.md").read_text(encoding="utf-8") != "# cmdr\n":
    raise SystemExit("Root README changed during generation")

# Remove the one-shot machinery from the final branch.
for path in [ROOT / ".cmdr-generator", ROOT / "tools" / "generate_cmdr_docs.py",
             ROOT / ".github" / "workflows" / "generate-cmdr-docs.yml"]:
    if path.is_dir():
        shutil.rmtree(path)
    elif path.exists():
        path.unlink()
for parent in [ROOT / "tools", ROOT / ".github" / "workflows", ROOT / ".github"]:
    try:
        parent.rmdir()
    except OSError:
        pass
run("git", "add", "-A")
if subprocess.run(["git", "diff", "--cached", "--quiet"], check=False).returncode != 0:
    run("git", "commit", "-m", "docs: validate canonical CMDR manifest")
run("git", "push", "origin", "HEAD:docs/cmdr-product-spec-foundation")
print(f"Generated {len(actual)} files across {len(groups)} domain commits")
