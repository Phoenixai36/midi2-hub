"""
Emotional Engine
Maps abstract emotional macros to MIDI 2.0 CC (32-bit resolution).
Used for live sets with narrative arcs (R3B0RN project and others).

Endpoint: POST /map
Input:  { "tension": 0.87, "density": 0.65, "brightness": 0.32,
          "complexity": 0.91, "release": 0.18 }
Output: MIDI 2.0 CC messages as JSON
"""
from fastapi import FastAPI
from fastapi.responses import JSONResponse
from pydantic import BaseModel, Field
from typing import Optional

app = FastAPI(title="midi2-hub emotional-engine", version="0.1.0")

MAX_32BIT = 4294967295  # 2^32 - 1

# CC mapping: emotional macro -> MIDI 2.0 CC number
EMOTION_CC_MAP = {
    "tension":    74,   # Filter cutoff
    "density":    91,   # Reverb depth / layer density
    "brightness": 72,   # Release time / spectral brightness
    "complexity": 76,   # Vibrato rate / arp complexity
    "release":    73,   # Attack time / tension release
}

# Emotional stage presets (R3B0RN arc: stages of grief)
STAGE_PRESETS = {
    "denial":       {"tension": 0.2, "density": 0.3, "brightness": 0.7, "complexity": 0.2, "release": 0.8},
    "anger":        {"tension": 0.9, "density": 0.8, "brightness": 0.3, "complexity": 0.7, "release": 0.1},
    "bargaining":   {"tension": 0.6, "density": 0.5, "brightness": 0.5, "complexity": 0.6, "release": 0.4},
    "depression":   {"tension": 0.4, "density": 0.2, "brightness": 0.1, "complexity": 0.3, "release": 0.2},
    "acceptance":   {"tension": 0.2, "density": 0.4, "brightness": 0.8, "complexity": 0.3, "release": 0.9},
}


class EmotionState(BaseModel):
    tension:    float = Field(0.5, ge=0.0, le=1.0)
    density:    float = Field(0.5, ge=0.0, le=1.0)
    brightness: float = Field(0.5, ge=0.0, le=1.0)
    complexity: float = Field(0.5, ge=0.0, le=1.0)
    release:    float = Field(0.5, ge=0.0, le=1.0)
    channel:    int   = Field(1, ge=1, le=16)


def to_32bit(value: float) -> int:
    """Convert [0,1] float to 32-bit MIDI 2.0 CC value."""
    return int(max(0.0, min(1.0, value)) * MAX_32BIT)


@app.get("/health")
def health():
    return {"status": "ok", "engine": "emotional-engine"}


@app.get("/stages")
def list_stages():
    return {"stages": list(STAGE_PRESETS.keys())}


@app.get("/stage/{name}")
def get_stage(name: str):
    if name not in STAGE_PRESETS:
        return JSONResponse(status_code=404, content={"error": f"Stage '{name}' not found"})
    state = EmotionState(**STAGE_PRESETS[name])
    return _build_cc_messages(state)


@app.post("/map")
def map_emotion(state: EmotionState):
    """Map emotional macros to MIDI 2.0 CC messages (32-bit resolution)."""
    return _build_cc_messages(state)


def _build_cc_messages(state: EmotionState):
    values = {
        "tension":    state.tension,
        "density":    state.density,
        "brightness": state.brightness,
        "complexity": state.complexity,
        "release":    state.release,
    }
    messages = []
    for macro, value in values.items():
        cc_num = EMOTION_CC_MAP[macro]
        cc_32bit = to_32bit(value)
        messages.append({
            "type": "cc",
            "channel": state.channel,
            "cc": cc_num,
            "macro": macro,
            "value_float": value,
            "value_32bit": cc_32bit,
            # MIDI 2.0 high-res normalized value
            "value_normalized": cc_32bit / MAX_32BIT,
        })
    return {"channel": state.channel, "messages": messages, "state": values}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8003)
