

Based on the provided context, I will create a documentation for the files inside the `reports/app/models/reports` directory.

**Directory Overview**

The `reports/app/models/reports` directory contains data models for working with reports. The models are defined using the Pydantic library.

**Files**

### general.py

This file contains the following data models:

* `Time_In_Audio`: Represents time in an audio recording.
	+ Fields:
		- `start`: start time in `time` format
		- `end`: end time in `time` format
* `Proposal`: Represents a key proposal.
	+ Fields:
		- `text`: proposal text in `str` format
		- `context`: proposal context in `str` format
		- `audio_time`: time in audio recording in `Time_In_Audio` format
* `Block`: Represents a main block.
	+ Fields:
		- `name_block`: block name in `str` format
		- `proposals`: list of proposals in `list[Proposal]` format

### unofficial_report.py

This file contains the `UnofficialReport` data model.

* `UnofficialReport`: Represents an unofficial report.
	+ Fields:
		- `title`: report title in `str` format
		- `content`: report content in `str` format

### official_report.py

This file contains the `OfficialReport` data model.

* `OfficialReport`: Represents an official report.
	+ Fields:
		- `title`: report title in `str` format
		- `content`: report content in `str` format
		- `blocks`: list of blocks in `list[Block]` format

### meeting_transcript.py

This file contains the `MeetingTranscript` data model.

* `MeetingTranscript`: Represents a meeting transcript.
	+ Fields:
		- `title`: transcript title in `str` format
		- `content`: transcript content in `str` format
