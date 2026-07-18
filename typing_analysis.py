# /// script
# requires-python = ">=3.13"
# dependencies = [matplotlib]
# ///

"""
Typing log analysis: fatigue proxies + character-level error/hand analysis,
with matplotlib visualizations.

Expects four keylogger JSON exports (one per condition) with this shape:
{
  "words": [
    {"text": "...", "progress": "...", "events": [
        {"key": "a", "date": "ISO8601", "expecteKey": "a"}, ...
    ]},
    ...
  ],
  "results": {"missing": ["word1", "word2", ...]},
  "startedAt": "ISO8601",
  "endedAt": "ISO8601"
}

Usage:
    python3 typing_analysis.py
Edit FILES below to point at your four JSON exports.
"""

import json
import statistics
from collections import Counter, defaultdict
from datetime import datetime

import matplotlib.pyplot as plt

FILES = {
    "flat-mechanical": "/tmp/flat-mechanical.json",
    "medium-angle-mechanical": "/tmp/medium-angle-mechanical.json",
    "high-angle-mechanical": "/tmp/high-angle-mechanical.json",
    "membrane": "/tmp/membrane.json",
}

LEFT = set("`1234qwertasdfgzxcvb")
RIGHT = set("5678290yuiophjklnm,./;'[]\\-=")


def parse_ts(s):
    s = s.rstrip("Z")
    if "." in s:
        head, frac = s.split(".")
        frac = (frac + "000000")[:6]
        s = f"{head}.{frac}"
    return datetime.fromisoformat(s)


def hand_of(ch):
    if ch is None:
        return None
    c = ch.lower()
    if c == " ":
        return "space"
    if c in LEFT:
        return "left"
    if c in RIGHT:
        return "right"
    return "other"


# ---------- Fatigue-proxy analysis (per file) ----------

def analyze_fatigue(path):
    with open(path) as f:
        data = json.load(f)

    words = data["words"]
    missing = set(data.get("results", {}).get("missing", []))

    all_events = []
    for wi, w in enumerate(words):
        for ev in w["events"]:
            all_events.append((wi, parse_ts(ev["date"])))
    all_events.sort(key=lambda x: x[1])

    start = parse_ts(data["startedAt"])
    end = parse_ts(data["endedAt"])
    total_seconds = (end - start).total_seconds()

    n_words = len(words)
    half = n_words // 2

    def half_stats(lo, hi):
        w_slice = words[lo:hi]
        first_t = min(parse_ts(ev["date"]) for w in w_slice for ev in w["events"])
        last_t = max(parse_ts(ev["date"]) for w in w_slice for ev in w["events"])
        span_min = max((last_t - first_t).total_seconds() / 60, 0.001)
        wpm = len(w_slice) / span_min
        errs = sum(1 for w in w_slice if w["text"] in missing)
        return wpm, errs / len(w_slice)

    first_half = half_stats(0, half)
    second_half = half_stats(half, n_words)

    intervals = []
    for i in range(1, len(all_events)):
        dt = (all_events[i][1] - all_events[i - 1][1]).total_seconds() * 1000
        if dt >= 0:
            intervals.append(dt)
    mid = len(intervals) // 2

    def stdev_safe(lst):
        return statistics.pstdev(lst) if len(lst) > 1 else 0

    total_errs = sum(1 for w in words if w["text"] in missing)

    return {
        "overall_wpm": n_words / (total_seconds / 60),
        "overall_error_rate": total_errs / n_words,
        "first_half_wpm": first_half[0],
        "second_half_wpm": second_half[0],
        "first_half_err_rate": first_half[1],
        "second_half_err_rate": second_half[1],
        "first_half_stdev_ms": stdev_safe(intervals[:mid]),
        "second_half_stdev_ms": stdev_safe(intervals[mid:]),
    }


# ---------- Error / hand analysis (combined across files) ----------

def analyze_errors(files):
    pair_counter = Counter()
    hand_error_counter = Counter()
    per_file_hand_errors = defaultdict(Counter)
    total_errors = 0

    for label, path in files.items():
        with open(path) as f:
            data = json.load(f)
        for w in data["words"]:
            for ev in w["events"]:
                typed, expected = ev.get("key"), ev.get("expecteKey")
                if typed is None or expected is None or typed == expected:
                    continue
                total_errors += 1
                pair_counter[(expected, typed)] += 1
                eh = hand_of(expected)
                hand_error_counter[eh] += 1
                per_file_hand_errors[label][eh] += 1

    return pair_counter, hand_error_counter, per_file_hand_errors, total_errors


def main():
    fatigue = {label: analyze_fatigue(path) for label, path in FILES.items()}
    pair_counter, hand_error_counter, per_file_hand_errors, total_errors = analyze_errors(FILES)

    labels = list(FILES.keys())

    fig, axes = plt.subplots(2, 2, figsize=(14, 10))
    fig.suptitle("Typing Test Analysis Across Keyboard Conditions", fontsize=15, fontweight="bold")

    # --- Panel 1: first vs second half WPM per condition ---
    ax = axes[0, 0]
    x = range(len(labels))
    width = 0.35
    first_wpm = [fatigue[l]["first_half_wpm"] for l in labels]
    second_wpm = [fatigue[l]["second_half_wpm"] for l in labels]
    ax.bar([i - width / 2 for i in x], first_wpm, width, label="First half", color="#3a7bd5")
    ax.bar([i + width / 2 for i in x], second_wpm, width, label="Second half", color="#e07a5f")
    ax.set_xticks(list(x))
    ax.set_xticklabels(labels, rotation=20, ha="right")
    ax.set_ylabel("WPM")
    ax.set_title("WPM: first half vs second half of trial")
    ax.legend()

    # --- Panel 2: first vs second half error rate per condition ---
    ax = axes[0, 1]
    first_err = [fatigue[l]["first_half_err_rate"] * 100 for l in labels]
    second_err = [fatigue[l]["second_half_err_rate"] * 100 for l in labels]
    ax.bar([i - width / 2 for i in x], first_err, width, label="First half", color="#3a7bd5")
    ax.bar([i + width / 2 for i in x], second_err, width, label="Second half", color="#e07a5f")
    ax.set_xticks(list(x))
    ax.set_xticklabels(labels, rotation=20, ha="right")
    ax.set_ylabel("Error rate (%)")
    ax.set_title("Error rate: first half vs second half of trial")
    ax.legend()

    # --- Panel 3: errors by hand, per condition (stacked) ---
    ax = axes[1, 0]
    hands = ["left", "right", "space"]
    colors = {"left": "#3a7bd5", "right": "#e07a5f", "space": "#8d99ae"}
    bottoms = [0] * len(labels)
    for hand in hands:
        vals = [per_file_hand_errors[l].get(hand, 0) for l in labels]
        ax.bar(labels, vals, bottom=bottoms, label=hand, color=colors[hand])
        bottoms = [b + v for b, v in zip(bottoms, vals)]
    ax.set_xticklabels(labels, rotation=20, ha="right")
    ax.set_ylabel("Error count")
    ax.set_title("Errors by hand (target key), per condition")
    ax.legend()

    # --- Panel 4: top 12 substitution pairs, combined ---
    ax = axes[1, 1]
    top_pairs = pair_counter.most_common(12)
    pair_labels = [f"{e!r}\u2192{t!r}" for (e, t), _ in top_pairs]
    pair_counts = [c for _, c in top_pairs]
    ax.barh(pair_labels[::-1], pair_counts[::-1], color="#3a7bd5")
    ax.set_xlabel("Count")
    ax.set_title("Top substitution pairs (combined, all 4 conditions)")

    plt.tight_layout(rect=[0, 0, 1, 0.96])
    plt.savefig("/mnt/user-data/outputs/typing_analysis.png", dpi=150)
    print("Saved plot to /mnt/user-data/outputs/typing_analysis.png")

    # Console summary
    print("\n=== Fatigue proxies ===")
    for l in labels:
        f = fatigue[l]
        print(f"{l}: overall_wpm={f['overall_wpm']:.1f} "
              f"first_half_wpm={f['first_half_wpm']:.1f} second_half_wpm={f['second_half_wpm']:.1f} "
              f"first_err={f['first_half_err_rate']*100:.1f}% second_err={f['second_half_err_rate']*100:.1f}%")

    print(f"\n=== Errors: total={total_errors} ===")
    for hand, cnt in hand_error_counter.most_common():
        print(f"{hand}: {cnt} ({cnt/total_errors*100:.1f}%)")


if __name__ == "__main__":
    main()
