from datetime import time

from pydantic import BaseModel


class Time_In_Audio(BaseModel):
    """Represents a time interval in an audio recording.

    Attributes:
        start (time): The start time of the interval.
        end (time): The end time of the interval.
    """

    start: time | None
    end: time | None


class Proposal(BaseModel):
    """Represents a main block.

    Attributes:
        name_block (str): The name of the block.
        proposals (list[Proposal] | None): A list of proposals associated with the block, or None if no proposals are associated.
    """

    text: str | None
    context: str | None
    audio_time: Time_In_Audio | None


class Block(BaseModel):
    """Основной блок."""

    name_block: str | None
    proposals: list[Proposal] | None
