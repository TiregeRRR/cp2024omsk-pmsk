from pydantic import BaseModel

from . import general


class SpeakerTranscript(BaseModel):
    """Data model for storing information about a single speaker's transcript.

    Attributes:
        speaker_name (str): The name of the speaker.
        transcript_text (str): The text of the speaker's transcript.
        audio_time (Time_In_Audio): Information about the audio time.
    """

    speaker_name: str | None
    transcript_text: str | None
    audio_time: general.Time_In_Audio | None


class MeetingTranscript(BaseModel):
    """Data model for storing information about a meeting transcript.

    Attributes:
        speakers_transcript (list[SpeakerTranscript]): Information about the transcripts of all speakers.
    """

    speakers_transcript: list[SpeakerTranscript] | None
