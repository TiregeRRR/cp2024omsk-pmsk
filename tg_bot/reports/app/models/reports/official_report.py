from datetime import datetime, time

from pydantic import BaseModel

from .general import Block


class Errand(BaseModel):
    """An errand.

    Attributes:
        assignee (str): The person assigned the errand.
        context (str): The context of the errand.
        deadline (datetime): The deadline for completing the errand.
    """

    assignee: str | None # кому поручено
    context: str | None # контекст поручения
    deadline: datetime | None # срок выполнения


class ErrandProtocol(BaseModel):
    """An errand protocol.

    Attributes:
        list_errands (list[Errand]): A list of errands.
    """

    list_errands: list[Errand] | None


class OfficialProtocol(BaseModel):
    """The official agenda.

    Attributes:
        date (datetime): The date of the event.
        time (time): The time of the event.
        attendees (list[str]): The list of attendees.
        blocks (list[Block]): The main blocks.
        errand_protocol (ErrandProtocol): The errand protocol.
    """

    date: datetime | None  # дата проведения
    time: time | None # время проведения
    attendees: list[str] | None  # перечень присутствующих
    blocks: list[Block] | None # Основные блоки
    errand_protocol: ErrandProtocol | None  # Протокол поручений
