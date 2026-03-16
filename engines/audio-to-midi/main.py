"""
audio-to-midi engine
Converts audio files to MIDI 2.0 clips using basic-pitch (Spotify).
Exposes a FastAPI endpoint: POST /convert
"""
import io
import json
from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.responses import JSONResponse
import numpy as np

app = FastAPI(title="midi2-hub audio-to-midi engine", version="0.1.0")


@app.get("/health")
def health():
    return {"status": "ok", "engine": "audio-to-midi"}


@app.post("/convert")
async def convert(file: UploadFile = File(...), slot: str = "melody"):
    """
    Upload an audio file, get back a MIDI 2.0 clip JSON.
    Slot: drums | bass | melody | lead
    """
    try:
        from basic_pitch.inference import predict
        from basic_pitch import ICASSP_2022_MODEL_PATH

        audio_bytes = await file.read()
        audio_buf = io.BytesIO(audio_bytes)

        model_output, midi_data, note_events = predict(
            audio_buf,
            ICASSP_2022_MODEL_PATH,
        )

        # Convert note events to MIDI 2.0 clip JSON format
        notes = []
        for note in note_events:
            notes.append({
                "pitch": int(note.pitch),
                "velocity": int(note.amplitude * 127),
                "start_beat": float(note.start_time),
                "end_beat": float(note.end_time),
                # MIDI 2.0 per-note attributes (32-bit resolution)
                "attributes": {
                    "pitch_bend": 0,
                    "pressure": int(note.amplitude * 4294967295),  # 32-bit
                }
            })

        clip = {
            "slot": slot,
            "resolution": 32,
            "format": "midi2-json",
            "notes": notes,
            "metadata": {
                "source": file.filename,
                "note_count": len(notes),
            }
        }

        return JSONResponse(content=clip)

    except ImportError:
        # basic_pitch not installed, return stub
        raise HTTPException(
            status_code=501,
            detail="basic-pitch not installed. Run: pip install basic-pitch"
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
