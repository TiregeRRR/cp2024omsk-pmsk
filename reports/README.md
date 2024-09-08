

**API Documentation**
=====================

**Endpoints**
------------

### POST /reports/official

**Description**
---------------

Create an official report meeting.

**Request Body**
----------------

* `official_data`: Input data for the official report, expected to be of type `input_data.Official_Protocol_Data`

**Response**
------------

* A `FileResponse` object containing the generated report

### POST /reports/unofficial

**Description**
---------------

Create an unofficial report meeting.

**Request Body**
----------------

* `unofficial_data`: Input data for the unofficial report, expected to be of type `input_data.Unofficial_Protocol_Data`

**Response**
------------

* A `FileResponse` object containing the generated report

### POST /transcript

**Description**
---------------

Create a transcript.

**Request Body**
----------------

* `transcript_data`: Input data for the transcript

**Response**
------------

* A `FileResponse` object containing the generated transcript

### GET /

**Description**
---------------

Redirect to the documentation.

**Response**
------------

* A redirect to the documentation page (`/docs`)


**Notes**
--------

* All endpoints are subject to change without notice.
* This documentation is for informational purposes only.