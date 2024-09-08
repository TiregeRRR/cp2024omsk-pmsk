from datetime import datetime, time, timedelta

from pydantic import BaseModel

from . import general


class UnofficialProtocol(BaseModel):
    """Unofficial protocol.

    Attributes:
        date (datetime): The date of the protocol.
        time (time): The time of the protocol.
        duration (timedelta): The duration of the protocol.
        participants (list[str] | None): The participants of the protocol.
        agenda (str | None): The agenda of the protocol.
        blocks (list[general.Block]): The blocks of the protocol.
        audio_times (list[general.Time_In_Audio]): The audio times of the protocol.
    """

    date: datetime | None
    time: time | None
    duration: timedelta | None
    participants: list[str] | None
    agenda: list[str] | None
    blocks: list[general.Block] | None
    audio_times: list[general.Time_In_Audio] | None
