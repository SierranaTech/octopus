#!/usr/bin/env python3

import re
import subprocess
from pathlib import Path


ROOT = Path(".")
ARTIFACTS = ROOT / ".artifacts"
BADGES = ROOT / ".github" / "badges"
README = ROOT / "README.md"

BADGES.mkdir(parents=True, exist_ok=True)


def read_text(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def parse_total_coverage() -> float:
    summary = read_text(ARTIFACTS / "coverage-summary.txt")
    match = re.search(r"total:\s+\(statements\)\s+([\d.]+)%", summary)
    if not match:
        raise RuntimeError("Could not parse total coverage from coverage-summary.txt")
    return float(match.group(1))


def parse_quality_grade() -> tuple[str, str | None]:
    text = read_text(ARTIFACTS / "goreportcard.txt")

    grade_match = re.search(r"Grade\s*[:. ]+\s*([A-F][+-]?)", text, re.IGNORECASE)
    if not grade_match:
        # fallback: first standalone grade-like token after 'grade'
        grade_match = re.search(r"grade.*?\b([A-F][+-]?)\b", text, re.IGNORECASE | re.DOTALL)

    if not grade_match:
        raise RuntimeError("Could not parse Go Report Card grade")

    score_match = re.search(r"\b([0-9]+(?:\.[0-9]+)?)%\b", text)
    return grade_match.group(1).upper(), score_match.group(1) if score_match else None


def color_for_coverage(value: float) -> str:
    if value >= 80:
        return "#4c1"
    if value >= 60:
        return "#97CA00"
    if value >= 40:
        return "#dfb317"
    return "#e05d44"


def color_for_grade(grade: str) -> str:
    if grade in {"A+", "A"}:
        return "#4c1"
    if grade in {"A-", "B+", "B"}:
        return "#97CA00"
    if grade in {"B-", "C+", "C"}:
        return "#dfb317"
    return "#e05d44"


def badge_svg(label: str, value: str, color: str) -> str:
    label_width = max(46, 7 * len(label) + 10)
    value_width = max(46, 7 * len(value) + 10)
    total_width = label_width + value_width

    return f"""<svg xmlns="http://www.w3.org/2000/svg" width="{total_width}" height="20" role="img" aria-label="{label}: {value}">
  <title>{label}: {value}</title>
  <linearGradient id="s" x2="0" y2="100%">
    <stop offset="0" stop-color="#fff" stop-opacity=".7"/>
    <stop offset=".1" stop-color="#aaa" stop-opacity=".1"/>
    <stop offset=".9" stop-color="#000" stop-opacity=".3"/>
    <stop offset="1" stop-color="#000" stop-opacity=".5"/>
  </linearGradient>
  <clipPath id="r">
    <rect width="{total_width}" height="20" rx="3" fill="#fff"/>
  </clipPath>
  <g clip-path="url(#r)">
    <rect width="{label_width}" height="20" fill="#555"/>
    <rect x="{label_width}" width="{value_width}" height="20" fill="{color}"/>
    <rect width="{total_width}" height="20" fill="url(#s)"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11">
    <text x="{label_width / 2}" y="15" fill="#010101" fill-opacity=".3">{label}</text>
    <text x="{label_width / 2}" y="14">{label}</text>
    <text x="{label_width + (value_width / 2)}" y="15" fill="#010101" fill-opacity=".3">{value}</text>
    <text x="{label_width + (value_width / 2)}" y="14">{value}</text>
  </g>
</svg>
"""


def update_readme(coverage: float, grade: str, score: str | None) -> None:
    generated = [
        "<!-- generated:quality:start -->",
        "",
        "[![Coverage](.github/badges/coverage.svg)](.artifacts/coverage-summary.txt)",
        "[![Quality](.github/badges/quality.svg)](#code-quality)",
        "",
        "## Code quality",
        "",
        f"- Coverage: **{coverage:.1f}%**",
        f"- Go Report Card grade: **{grade}**",
    ]

    if score is not None:
        generated.append(f"- Go Report Card score: **{score}%**")

    generated.extend([
        "",
        "_These badges are generated in CI for the private repo and committed back into the repository._",
        "",
        "<!-- generated:quality:end -->",
    ])

    block = "\n".join(generated)

    current = README.read_text(encoding="utf-8") if README.exists() else ""
    pattern = re.compile(
        r"<!-- generated:quality:start -->.*?<!-- generated:quality:end -->",
        re.DOTALL,
    )

    if pattern.search(current):
        updated = pattern.sub(block, current)
    else:
        updated = current.rstrip() + ("\n\n" if current.strip() else "") + block + "\n"

    README.write_text(updated, encoding="utf-8")


def main() -> None:
    coverage = parse_total_coverage()
    grade, score = parse_quality_grade()

    (BADGES / "coverage.svg").write_text(
        badge_svg("coverage", f"{coverage:.1f}%", color_for_coverage(coverage)),
        encoding="utf-8",
    )
    (BADGES / "quality.svg").write_text(
        badge_svg("quality", grade, color_for_grade(grade)),
        encoding="utf-8",
    )

    update_readme(coverage, grade, score)


if __name__ == "__main__":
    main()
