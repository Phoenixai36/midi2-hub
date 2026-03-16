cmake_minimum_required(VERSION 3.22)
project(Midi2HubPlugin VERSION 0.1.0)

# --------------------------------------------------------------------------
# JUCE (add as submodule: git submodule add https://github.com/juce-framework/JUCE)
# --------------------------------------------------------------------------
add_subdirectory(JUCE)

juce_add_plugin(Midi2HubPlugin
  COMPANY_NAME          "Phoenixai36"
  IS_SYNTH              FALSE
  NEEDS_MIDI_INPUT      TRUE
  NEEDS_MIDI_OUTPUT     TRUE
  IS_MIDI_EFFECT        TRUE
  EDITOR_WANTS_KEYBOARD_FOCUS FALSE
  COPY_PLUGIN_AFTER_BUILD     FALSE
  PLUGIN_MANUFACTURER_CODE    Pxa3
  PLUGIN_CODE             M2Hb
  FORMATS                 VST3 Standalone
  PRODUCT_NAME           "MIDI 2 Hub"
)

target_sources(Midi2HubPlugin
  PRIVATE
    src/PluginProcessor.cpp
)

target_compile_definitions(Midi2HubPlugin
  PUBLIC
    JUCE_WEB_BROWSER=0
    JUCE_USE_CURL=0
    JUCE_VST3_CAN_REPLACE_VST2=0
)

target_link_libraries(Midi2HubPlugin
  PRIVATE
    juce::juce_audio_utils
    juce::juce_audio_processors
  PUBLIC
    juce::juce_recommended_config_flags
    juce::juce_recommended_lto_flags
    juce::juce_recommended_warning_flags
)
