from enum import Enum

from pydantic import BaseModel

from .reports import meeting_transcript, official_report, unofficial_report


class Type_Document(str, Enum):
    """Enumeration of document types.

    Attributes:
        docx (str): Represents a docx document.
        pdf (str): Represents a pdf document.
    """

    docx = "docx"
    pdf = "pdf"


class Meeting_Transcript_Data(BaseModel):
    """Data model for meeting transcript.

    Attributes:
        name_report (str): The name of the report.
        document_type (Type_Document): The type of document (docx or pdf).
        password (str | None): The password for encrypted documents (optional).
        data (meeting_transcript.MeetingTranscript): The meeting transcript data.
    """

    name_report: str
    document_type: Type_Document
    password: str | None = None
    data: meeting_transcript.MeetingTranscript


class Official_Protocol_Data(BaseModel):
    """Data model for official protocol.

    Attributes:
        name_report (str): The name of the report.
        document_type (Type_Document): The type of document (docx or pdf).
        password (str | None): The password for encrypted documents (optional).
        data (official_report.OfficialProtocol): The official protocol data.
    """

    name_report: str
    document_type: Type_Document
    password: str | None = None
    data: official_report.OfficialProtocol


class Unofficial_Protocol_Data(BaseModel):
    """Data model for unofficial protocol.

    Attributes:
        name_report (str): The name of the report.
        document_type (Type_Document): The type of document (docx or pdf).
        password (str | None): The password for encrypted documents (optional).
        data (unofficial_report.UnofficialProtocol): The unofficial protocol data.
    """

    name_report: str
    document_type: Type_Document
    password: str | None = None
    data: unofficial_report.UnofficialProtocol
