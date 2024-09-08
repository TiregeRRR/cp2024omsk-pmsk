import os

from dotenv import load_dotenv
from fastapi import FastAPI
from fastapi.responses import FileResponse, RedirectResponse

from . import report
from .models import input_data

load_dotenv()

app = FastAPI()


def check_reports_dir() -> None:
    """Check if the reports directory exists.

    If not, create it.
    """
    if not os.path.exists("reports"):
        os.mkdir("reports")


check_reports_dir()


@app.post("/reports/official")
async def official_report(
    official_data: input_data.Official_Protocol_Data,
) -> FileResponse:
    """Create an official report meeting.

    Args:
        official_data: Input data for the official report.

    Returns:
        A file response containing the generated report.
    """
    # Generate the report with the given input data
    file = report.generate_report(official_data)

    # Return the report as a file response
    return FileResponse(
        file,
        # Use the name provided in the input data, with the correct file extension
        filename=f"{official_data.name_report}.{official_data.document_type.value}",
    )


@app.post("/reports/unofficial")
async def unofficial_report(
    unofficial_data: input_data.Unofficial_Protocol_Data,
) -> FileResponse:
    """Create an unofficial report meeting.

    Args:
        unofficial_data: Input data for the unofficial report.

    Returns:
        A file response containing the generated report.
    """
    # Generate the report with the given input data
    file = report.generate_report(unofficial_data)

    # Return the report as a file response
    return FileResponse(
        file,
        # Use the name provided in the input data, with the correct file extension
        filename=f"{unofficial_data.name_report}.{unofficial_data.document_type.value}",
    )


@app.post("/transcript")
async def transcript(transcript_data: input_data.Meeting_Transcript_Data) -> FileResponse:
    """Create a transcript meeting.

    Args:
        transcript_data: Input data for the meeting transcript.

    Returns:
        A file response containing the generated transcript.
    """
    # Generate the report with the given input data
    file = report.generate_report(transcript_data)

    # Return the report as a file response
    return FileResponse(
        file,
        # Use the name provided in the input data, with the correct file extension
        filename=f"{transcript_data.name_report}.{transcript_data.document_type.value}",
    )


@app.get("/", response_class=RedirectResponse, include_in_schema=False)
async def index() -> RedirectResponse:
    """Redirect to the documentation.

    This endpoint is not included in the OpenAPI schema (i.e. it is not visible in the
    documentation) and is only used to redirect users who access the root URL of the service
    to the documentation.
    """
    return "/docs"
