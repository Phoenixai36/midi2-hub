#include <juce_audio_processors/juce_audio_processors.h>

// ---------------------------------------------------------------------------
// Midi2HubProcessor  –  MIDI 2.0 hub plugin skeleton
// Receives MIDI 2.0 UMP packets, forwards them to the Go session server
// via WebSocket, and injects back incoming events from remote peers.
// ---------------------------------------------------------------------------

class Midi2HubProcessor : public juce::AudioProcessor
{
public:
  Midi2HubProcessor()
    : AudioProcessor(BusesProperties()) {}

  ~Midi2HubProcessor() override = default;

  // --- AudioProcessor overrides -------------------------------------------
  const juce::String getName() const override { return "MIDI 2 Hub"; }
  bool acceptsMidi()  const override { return true;  }
  bool producesMidi() const override { return true;  }
  bool isMidiEffect() const override { return true;  }
  double getTailLengthSeconds() const override { return 0.0; }

  int  getNumPrograms()    override { return 1; }
  int  getCurrentProgram() override { return 0; }
  void setCurrentProgram(int) override {}
  const juce::String getProgramName(int) override { return {}; }
  void changeProgramName(int, const juce::String&) override {}

  void prepareToPlay(double, int) override {}
  void releaseResources() override {}

  void processBlock(juce::AudioBuffer<float>& buffer,
                    juce::MidiBuffer& midiMessages) override
  {
    buffer.clear();

    // TODO Sprint-2: serialize incoming MIDI to JSON and send via WebSocket
    // TODO Sprint-2: inject remote MIDI events received from Go server
    (void)midiMessages;
  }

  bool hasEditor() const override { return false; }
  juce::AudioProcessorEditor* createEditor() override { return nullptr; }

  void getStateInformation(juce::MemoryBlock&) override {}
  void setStateInformation(const void*, int) override {}
};

// Required by JUCE plugin factory
juce::AudioProcessor* JUCE_CALLTYPE createPluginFilter()
{
  return new Midi2HubProcessor();
}
